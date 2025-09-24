package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	logLevel := os.Getenv("API_LOGGER_LEVEL")

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

	config := zap.NewDevelopmentConfig()

	config.Level = level
	config.DisableStacktrace = true
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stdout"}

	return config.Build()
}
