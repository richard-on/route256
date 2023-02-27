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
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	lomsservice "gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	model := domain.New()

	reflection.Register(s)
	loms.RegisterLOMSServer(s, lomsservice.New(model))

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

	log.Info("shutting down: loms service")
	s.GracefulStop()
}
