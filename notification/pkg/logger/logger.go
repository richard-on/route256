package logger

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
}
