// LOMS (Logistics and Order Management System) handles order management and logistics.
package main

import (
	"runtime"

	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config load fail")
	}

	log.Info().Msgf("loms service: version=%v, build=%v, protoVersion=%v, go version=%v",
		config.Version, config.Build, config.ProtoVersion, runtime.Version())

	app.Run(cfg)
}
