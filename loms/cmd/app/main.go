package main

import (
	"runtime"

	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/app"
)

var (
	Version string // Version of this app.
	Build   string // Build date and time.
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
