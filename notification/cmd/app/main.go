package main

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().
			Str("component", "notification-init").
			Err(err).
			Msg("config load fail")
	}

	log.Info().
		Str("component", fmt.Sprintf("%v-init", cfg.Service.Name)).
		Msgf("%v: version=%v, build=%v, go version=%v",
			cfg.Service.Name, config.Version, config.Build, runtime.Version())

	app.Run(cfg)
}
