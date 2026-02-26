package a

import "log/slog"

// diagnostics-only: все правила
func demo() {
	slog.Info("Starting server")    // want "log messages should start with a lowercase letter"
	slog.Info("запуск сервера")     // want "log messages must be written using English characters only"
	slog.Info("server started!!!")  // want "avoid punctuation, symbols, or emoji in log messages"
	slog.Info("token: 123") 		// want "avoid punctuation, symbols, or emoji in log messages" "possible sensitive information detected in log message"
}