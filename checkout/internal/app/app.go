package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/config"
	checkoutservice "gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/api/checkout"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/client/grpc/productservice"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

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

	model := domain.New(lomsClient, productClient, lomsClient)

	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, checkoutservice.New(model))

	go func() {
		if err = s.Serve(listener); !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal(err, "error while running grpc server")
		}
	}()
	log.Infof("grpc server listening at %v", listener.Addr())

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	<-ctx.Done()
	cancel()

	log.Info("shutting down: checkout service")
	s.GracefulStop()
}
