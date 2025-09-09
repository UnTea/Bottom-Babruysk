package logger

import "go.uber.org/zap"

func New() (*zap.Logger, error) {
	EnvironmentApiLogLevel := zap.DebugLevel.String()

	logLevel := EnvironmentApiLogLevel

	var level zap.AtomicLevel

	switch logLevel {
	case `debug`:
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case `warning`:
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case `error`:
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case `panic`:
		level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case `fatal`:
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	conf := zap.NewDevelopmentConfig()

	conf.Level = level
	conf.DisableStacktrace = true
	conf.OutputPaths = []string{"stdout"}
	conf.ErrorOutputPaths = []string{"stdout"}

	return conf.Build()
}
