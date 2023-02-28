package main

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/app"
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

	log.Info().Msgf("notification service: version=%v, build=%v, go version=%v",
		Version, Build, runtime.Version())

	app.Run(cfg)
}
