package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
)

// Logger is a wrapper for a logging library.
type Logger struct {
	log zerolog.Logger
}

// New creates a new Logger.
func New(out io.Writer, level string, serviceName string) Logger {
	var logLevel zerolog.Level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.NoLevel
	}

	return Logger{
		log: zerolog.New(out).
			Level(logLevel).
			With().
			Timestamp().
			Str("service", serviceName).
			Logger(),
	}

}

func (l Logger) Println(v ...interface{}) {
	l.log.Print(fmt.Sprint(v...))
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

func (l Logger) Debug(i ...interface{}) {
	l.log.Debug().Msgf(fmt.Sprint(i...))
}

func (l Logger) Debugf(format string, i ...interface{}) {
	l.log.Debug().Msgf(format, i...)
}

func (l Logger) Info(i ...interface{}) {
	l.log.Info().Msgf(fmt.Sprint(i...))
}

func (l Logger) Infof(format string, i ...interface{}) {
	l.log.Info().Msgf(format, i...)
}

func (l Logger) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
}

func (l Logger) Errorf(err error, format string, i ...interface{}) {
	l.log.Error().Err(err).Msgf(format, i...)
}

func (l Logger) Fatal(err error, msg string) {
	l.log.Fatal().Err(err).Msg(msg)
}

func (l Logger) Fatalf(err error, format string, i ...interface{}) {
	l.log.Fatal().Err(err).Msgf(format, i...)
}
