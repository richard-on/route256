package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	log := logger.New(
		os.Stdout, //zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info().Msg("config and logger init success")

	listener, err := net.Listen("tcp", net.JoinHostPort("", cfg.GRPC.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("error while creating listener")
	}

	model := domain.New()

	s := grpc.NewServer()

	reflection.Register(s)
	loms.RegisterLOMSServer(s, lomsservice.New(model))

	go func() {
		if err = s.Serve(listener); !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal().Err(err).Msg("error while running grpc server")
		}
	}()
	log.Info().Msgf("grpc server listening at %v", listener.Addr())

	<-ctx.Done()
	cancel()

	log.Info().Msg("shutting down: loms service")
	s.GracefulStop()
}
