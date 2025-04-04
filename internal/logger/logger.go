package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// take the env variable from the config file and map it into the value zap uses
func New(level string) *zap.Logger {
	config := zap.NewProductionConfig()
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, _ := config.Build()
	return logger
}
