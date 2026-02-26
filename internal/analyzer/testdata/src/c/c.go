package c

import (
	"go.uber.org/zap"
)

func testZap() {
	logger, _ := zap.NewProduction()

	logger.Info("Starting server")   // want "log messages should start with a lowercase letter"
	logger.Info("запуск сервера")    // want "log messages must be written using English characters only"
	logger.Info("server started!!!") // want "avoid punctuation, symbols, or emoji in log messages"
	logger.Info("apikey 123")        // want "possible sensitive information detected in log message"
}