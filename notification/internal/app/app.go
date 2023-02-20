package app

import (
	"os"
	"route256/lib/logger"
	"route256/notification/config"
)

func Run(cfg *config.Config) {
	log := logger.New(
		os.Stdout, //zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info().Msg("config and logger init success")

	// Add services

	log.Info().Msg("shutting down: notification service")
}
