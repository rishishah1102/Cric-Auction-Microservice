package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// newLogger initializes a custom zap logger
func newLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Will not give error as all the configs are defined as constants
	logger, _ := config.Build()
	globalLogger = logger
}

// WrapError logs and wraps an existing error
func WrapError(err error, msg string, fields ...zap.Field) error {
	fields = append(fields, zap.Error(err))
	globalLogger.Error(msg, fields...)
	return fmt.Errorf("%s: %w", msg, err)
}

// Get returns the initialized logger instance
func Get() *zap.Logger {
	if globalLogger == nil {
		newLogger()
	}
	return globalLogger
}
