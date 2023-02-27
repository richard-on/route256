package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	log := logger.New(
		os.Stdout, //zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info().Msg("config and logger init success")

	lomsConn, err := grpc.Dial(cfg.LOMS.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("error while dialing loms grpc server")
	}
	lomsClient := loms.NewClient(lomsConn)

	productConn, err := grpc.Dial(cfg.ProductService.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("error while dialing product service grpc server")
	}
	productClient := productservice.NewClient(productConn, cfg.ProductService.Token)

	model := domain.New(lomsClient, productClient, lomsClient)

	listener, err := net.Listen("tcp", net.JoinHostPort("", cfg.GRPC.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("error while creating listener")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, checkoutservice.New(model))

	go func() {
		if err = s.Serve(listener); !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal().Err(err).Msg("error while running grpc server")
		}
	}()
	log.Info().Msgf("grpc server listening at %v", listener.Addr())

	<-ctx.Done()
	cancel()

	log.Info().Msg("shutting down: checkout service")
	s.GracefulStop()
}
