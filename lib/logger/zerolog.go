package logger

import (
	"github.com/rs/zerolog"
	"io"
)

func NewLogger(out io.Writer, level string, serviceName string) zerolog.Logger {
	var logLevel zerolog.Level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.NoLevel
	}

	return zerolog.New(out).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()
}
