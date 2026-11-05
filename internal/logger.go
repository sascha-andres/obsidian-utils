package internal

import (
	"log/slog"
	"os"
)

// CreateLogger initializes and returns a new slog.Logger with the specified log level and project name.
// The log level can be "debug", "info", "warn" or "error". Defaults to "info" if an unknown level is provided.
func CreateLogger(logLevel, project string) *slog.Logger {
	level := slog.LevelInfo
	// an unknown level leaves level unchanged
	_ = level.UnmarshalText([]byte(logLevel))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).With("project", project)
	slog.SetDefault(logger)
	return logger
}
