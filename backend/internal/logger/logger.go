package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}

// GetLogger returns the singleton logger instance
func GetLogger() *slog.Logger {
	return logger
}
