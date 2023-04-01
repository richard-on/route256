package main

import (
	"runtime"

	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config load fail")
	}

	log.Info().Msgf("notification service: version=%v, build=%v, go version=%v",
		config.Version, config.Build, runtime.Version())

	app.Run(cfg)
}
