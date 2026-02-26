package b

import "log/slog"

func demo() {
	slog.Info("Starting server") // want "log messages should start with a lowercase letter"
}