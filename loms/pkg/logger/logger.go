package logger

import (
	"github.com/jackc/pgconn"
	"google.golang.org/grpc"
	"time"
)

type Logger interface {
	Debug(i ...interface{})
	Debugf(format string, i ...interface{})

	Info(i ...interface{})
	Infof(format string, i ...interface{})

	Warn(i ...interface{})
	Warnf(format string, i ...interface{})

	Error(err error, msg string)
	Errorf(err error, format string, i ...interface{})

	Fatal(err error, msg string)
	Fatalf(err error, format string, i ...interface{})

	GRPCLog
	DBLog
}

type GRPCLog interface {
	DebugGRPC(req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error)
	WarnGRPC(req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error)
	ErrorGRPC(req, resp interface{}, info *grpc.UnaryServerInfo, now time.Time, errors ...error)
}

type DBLog interface {
	RawSQL(method, sql string, args any)
	PGTag(method string, tag pgconn.CommandTag, errors ...error)
}
