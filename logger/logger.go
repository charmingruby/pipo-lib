// Package logger provides a logger instrumentation.
package logger

import (
	"log/slog"
	"os"
)

// Logger is a slog.Logger.
type Logger = slog.Logger

// New constructs a new logger.
//
// Returns a new logger.
func New() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	return logger
}
