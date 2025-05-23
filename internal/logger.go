package internal

import (
	"log/slog"
	"os"
	"strings"
)

// CreateLogger initializes and returns a new slog.Logger with the specified log level and project name.
// The log level can be "warn", "info", or "debug". Defaults to "info" if an unknown level is provided.
func CreateLogger(logLevel, project string) *slog.Logger {
	var handlerOpts *slog.HandlerOptions
	switch strings.ToLower(logLevel) {
	case "warn":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "info":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelInfo}
	case "debug":
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelDebug}
	default:
		handlerOpts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts)).With("project", project)
	slog.SetDefault(logger)
	return logger
}
