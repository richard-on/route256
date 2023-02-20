package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/checkout/config"
	"route256/checkout/internal/client/loms"
	"route256/checkout/internal/client/productservice"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handler/addtocart"
	"route256/checkout/internal/handler/deletefromcart"
	"route256/checkout/internal/handler/listcart"
	purchasehandler "route256/checkout/internal/handler/purchase"
	"route256/lib/logger"
	"route256/lib/server/wrapper"
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

	lomsClient := loms.New(cfg.LOMS.URL)
	productClient := productservice.New(
		cfg.ProductService.URL,
		cfg.ProductService.Token,
	)

	model := domain.New(lomsClient, productClient, lomsClient)

	addToCart := addtocart.New(model)
	deleteFromCart := deletefromcart.New()
	listCart := listcart.New(model)
	purchase := purchasehandler.New(model)

	http.Handle("/addToCart", wrapper.New(addToCart.Handle))
	http.Handle("/deleteFromCart", wrapper.New(deleteFromCart.Handle))
	http.Handle("/listCart", wrapper.New(listCart.Handle))
	http.Handle("/purchase", wrapper.New(purchase.Handle))

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

	log.Info().Msg("shutting down: checkout service")
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot shutdown http server")
	}
}
