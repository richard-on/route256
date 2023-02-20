package app

import (
	"context"
	"github.com/rs/zerolog"
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
	"time"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	log := logger.NewLogger(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		cfg.Log.Level,
		"logistics and order management system",
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

	go func() {
		err := http.ListenAndServe(":"+cfg.HTTP.Port, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot listen http")
		}
	}()

	<-ctx.Done()
	cancel()
	log.Info().Msg("shutting down: logistics and order management system")
}
