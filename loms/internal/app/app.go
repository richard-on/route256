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
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	lomsservice "gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
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

	model := domain.New(repo, tx)

	grpchealth.RegisterHealthServer(s, health.NewServer())
	reflection.Register(s)
	loms.RegisterLOMSServer(s, lomsservice.New(model))

	go func() {
		err = s.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal(err, "error while running grpc server")
		}
	}()
	log.Infof("grpc server listening at %v", listener.Addr())

	ticker := time.NewTicker(cfg.CancelInterval)
	// Start a separate goroutine to check and cancel unpaid orders.
	go func() {
		for {
			select {
			// Run cancelling on each tick.
			case <-ticker.C:
				errSlice := model.CancelUnpaidOrders(ctx, cfg.Service.PaymentTimeout)
				if len(errSlice) > 0 {
					for _, cancelErr := range errSlice {
						log.Error(cancelErr, "cancelling unpaid orders")
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
	cancel()

	log.Info("shutting down: loms service")
	s.GracefulStop()
}
