package zerolog

import (
	"fmt"
	"github.com/jackc/pgconn"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/logger"
	"google.golang.org/grpc"
	"io"
	"time"

	"github.com/rs/zerolog"
)

// Log  is a wrapper for a logging library.
type Log struct {
	log *zerolog.Logger
}

// New creates a new Logger.
func New(out io.Writer, level string, componentName string) logger.Logger {

	var logLevel zerolog.Level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	skipFrameCount := 1
	log := zerolog.New(out).
		Level(logLevel).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount+skipFrameCount).
		Str("component", componentName).
		Logger()

	return &Log{
		log: &log,
	}
}

func (l *Log) Debug(i ...interface{}) {
	l.log.Debug().Msgf(fmt.Sprint(i...))
}

func (l *Log) Debugf(format string, i ...interface{}) {
	l.log.Debug().Msgf(format, i...)
}

func (l *Log) Info(i ...interface{}) {
	l.log.Info().Msgf(fmt.Sprint(i...))
}

func (l *Log) Infof(format string, i ...interface{}) {
	l.log.Info().Msgf(format, i...)
}

func (l *Log) Warn(i ...interface{}) {
	l.log.Info().Msgf(fmt.Sprint(i...))
}

func (l *Log) Warnf(format string, i ...interface{}) {
	l.log.Info().Msgf(format, i...)
}

func (l *Log) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
}

func (l *Log) Errorf(err error, format string, i ...interface{}) {
	l.log.Error().Err(err).Msgf(format, i...)
}

func (l *Log) Fatal(err error, msg string) {
	l.log.Fatal().Err(err).Msg(msg)
}

func (l *Log) Fatalf(err error, format string, i ...interface{}) {
	l.log.Fatal().Err(err).Msgf(format, i...)
}

func (l *Log) DebugGRPC(req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error) {
	l.handleGRPC(l.log.Debug(), req, resp, info, now, errors...)
}

func (l *Log) WarnGRPC(req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error) {
	l.handleGRPC(l.log.Warn(), req, resp, info, now, errors...)
}

func (l *Log) ErrorGRPC(req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error) {
	l.handleGRPC(l.log.Error(), req, resp, info, now, errors...)
}

func (l *Log) handleGRPC(event *zerolog.Event, req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error) {
	for _, err := range errors {
		event.Err(err)
	}
	event.
		Interface("server", info.Server).
		Str("method", info.FullMethod).
		Interface("request", req).
		Interface("response", resp).
		Str("latency", time.Since(now).String()).
		Msg("handling gRPC request/response")
}

func (l *Log) RawSQL(method, sql string, args any) {
	l.log.Debug().
		Str("method", method).
		Str("sql", sql).
		Interface("args", args).
		Msg("raw sql query")
}

func (l *Log) PGTag(method string, tag pgconn.CommandTag, errors ...error) {
	var op string
	switch {
	case tag.Update() == true:
		op = "update"
	case tag.Delete() == true:
		op = "delete"
	case tag.Select() == true:
		op = "select"
	case tag.Insert() == true:
		op = "insert"
	default:
		op = "no_op"
	}

	event := l.log.Debug()
	for _, err := range errors {
		event.Err(err)
	}
	event.
		Str("method", method).
		Str("sql_op", op).
		Int64("rows_affected", tag.RowsAffected()).
		Msg("postgres execute result")
}
