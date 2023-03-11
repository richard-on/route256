// Package app initializes server, establishes connection with database and other dependencies.
// It then runs the service, accepting incoming connections and listening for shutdown signals.
package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/config"
	checkoutservice "gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/api/checkout"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/productservice"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/transactor"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Run creates and runs the service using provided config.
func Run(cfg *config.Config) {
	log := logger.New(
		os.Stdout,
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info("config and logger init success")

	lomsConn, err := grpc.Dial(cfg.LOMS.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err, "error while dialing loms grpc server")
	}
	lomsClient := loms.NewClient(lomsConn)

	productConn, err := grpc.Dial(cfg.ProductService.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err, "error while dialing product service grpc server")
	}
	productClient := productservice.NewClient(productConn, cfg.ProductService.Token)

	listener, err := net.Listen("tcp", net.JoinHostPort("", cfg.GRPC.Port))
	if err != nil {
		log.Fatal(err, "error while creating listener")
	}

	grpcLogger := logger.New(
		os.Stdout,
		cfg.Log.Level,
		fmt.Sprintf("%v-grpc", cfg.Service.Name),
	)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			logger.UnaryServerInterceptor(grpcLogger),
		)),
	)

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB))
	if err != nil {
		log.Fatal(err, "parsing database config")
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgConfig)
	if err != nil {
		log.Fatal(err, "connecting to database")
	}
	defer pool.Close()

	tx := transactor.New(pool)
	repo := repository.New(tx, tx)

	model := domain.New(repo, tx, lomsClient, productClient, lomsClient)

	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, checkoutservice.New(model))

	go func() {
		if err = s.Serve(listener); !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal(err, "error while running grpc server")
		}
	}()
	log.Infof("grpc server listening at %v", listener.Addr())

	<-ctx.Done()
	cancel()

	log.Info("shutting down: checkout service")
	s.GracefulStop()
}
