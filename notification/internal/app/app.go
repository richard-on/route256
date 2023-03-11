package app

import (
	"os"

	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
)

func Run(cfg *config.Config) {
	log := logger.New(
		os.Stdout,
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info("config and logger init success")

	// Blank service

	log.Info("shutting down: notification service")
}
