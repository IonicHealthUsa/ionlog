package internal

import "log/slog"

func createDefaultLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		&logger.writerHandler,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}
