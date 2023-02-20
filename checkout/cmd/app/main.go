// Checkout service is responsible for cart and order creation.
package main

import (
	"route256/checkout/config"
	"route256/checkout/internal/app"
	"runtime"

	"github.com/rs/zerolog/log"
)

var (
	Version string // Version of this app
	Build   string // Build date and time
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config load fail")
	}

	log.Info().Msgf("checkout service: version=%v, build=%v, go version=%v",
		Version, Build, runtime.Version())

	app.Run(cfg)
}
