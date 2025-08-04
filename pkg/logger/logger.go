// internal/logger/logger.go
package logger

import "go.uber.org/zap"

var globalLogger *zap.Logger

func Init(production bool) *zap.Logger {
	var logger *zap.Logger
	var err error

	if production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	globalLogger = logger
	return logger
}

func Get() *zap.Logger {
	if globalLogger == nil {
		panic("logger not initialized")
	}
	return globalLogger
}
