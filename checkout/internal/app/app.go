package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"route256/checkout/config"
	"route256/checkout/internal/client/loms"
	"route256/checkout/internal/client/productservice"
	"route256/checkout/internal/domain/cart"
	"route256/checkout/internal/domain/list"
	purchasedomain "route256/checkout/internal/domain/purchase"
	"route256/checkout/internal/handler/addtocart"
	"route256/checkout/internal/handler/deletefromcart"
	"route256/checkout/internal/handler/listcart"
	purchasehandler "route256/checkout/internal/handler/purchase"
	"route256/lib/logger"
	"route256/lib/server/wrapper"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	log := logger.NewLogger(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		cfg.Log.Level,
		"checkout",
	)
	log.Info().Msg("config and logger init success")

	lomsClient := loms.New(cfg.LOMS.URL)
	stockCheck := cart.NewChecker(lomsClient)
	addToCart := addtocart.New(stockCheck)

	deleteFromCart := deletefromcart.New()

	productClient := productservice.New(
		cfg.ProductService.URL,
		cfg.ProductService.Token,
	)
	productInfo := list.New(productClient)
	listCart := listcart.New(productInfo)

	orderCreate := purchasedomain.New(lomsClient)
	purchase := purchasehandler.New(orderCreate)

	http.Handle("/addToCart", wrapper.New(addToCart.Handle))
	http.Handle("/deleteFromCart", wrapper.New(deleteFromCart.Handle))
	http.Handle("/listCart", wrapper.New(listCart.Handle))
	http.Handle("/purchase", wrapper.New(purchase.Handle))

	go func() {
		err := http.ListenAndServe(":"+cfg.HTTP.Port, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot listen http")
		}
	}()

	<-ctx.Done()
	cancel()
	log.Info().Msg("shutting down: checkout service")
}
