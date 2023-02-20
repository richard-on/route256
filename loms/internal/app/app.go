package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/lib/logger"
	"route256/lib/server/wrapper"
	"route256/loms/config"
	cancelorder "route256/loms/internal/handler/cancel"
	"route256/loms/internal/handler/create"
	"route256/loms/internal/handler/list"
	"route256/loms/internal/handler/pay"
	stockhandler "route256/loms/internal/handler/stock"
	"syscall"

	"github.com/pkg/errors"
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

	createOrder := create.New()
	listOrder := list.New()
	payedOrder := pay.New()
	cancelOrder := cancelorder.New()
	stock := stockhandler.New()

	http.Handle("/createOrder", wrapper.New(createOrder.Handle))
	http.Handle("/listOrder", wrapper.New(listOrder.Handle))
	http.Handle("/orderPayed", wrapper.New(payedOrder.Handle))
	http.Handle("/cancelOrder", wrapper.New(cancelOrder.Handle))
	http.Handle("/stocks", wrapper.New(stock.Handle))

	httpServer := http.Server{
		Addr:         net.JoinHostPort("", cfg.HTTP.Port),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("error while running http server")
		}
	}()

	<-ctx.Done()
	cancel()

	ctx, shutdown := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer shutdown()

	log.Info().Msg("shutting down: logistics and order management system")
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot shutdown http server")
	}
}
