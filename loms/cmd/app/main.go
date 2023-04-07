// LOMS (Logistics and Order Management System) handles order management and logistics.
package main

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().
			Str("component", "loms-init").
			Err(err).
			Msg("config load fail")
	}

	log.Info().
		Str("component", fmt.Sprintf("%v-init", cfg.Service.Name)).
		Msgf("%v: version=%v, build=%v, protoVersion=%v, go version=%v",
			cfg.Service.Name, config.Version, config.Build, config.ProtoVersion, runtime.Version())

	app.Run(cfg)
}
