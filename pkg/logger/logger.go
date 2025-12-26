package logger

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	var h slog.Handler
	var lvl slog.Level

	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}

	h = slog.NewTextHandler(os.Stdout, opts)

	return slog.New(h)
}
