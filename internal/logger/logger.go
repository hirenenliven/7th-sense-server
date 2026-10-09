package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

// Level represents the log level.
type Level string

const (
	LevelInfo  Level = "info"
	LevelError Level = "error"
	LevelDebug Level = "debug"
	LevelWarn  Level = "warn"
)

// Init initializes the global logger with JSON output and INFO level.
func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(log)
}

// Log logs a message with the specified level and key-value pairs.
func Log(level Level, msg string, args ...any) {
	switch level {
	case LevelInfo:
		slog.Info(msg, args...)
	case LevelError:
		slog.Error(msg, args...)
	case LevelDebug:
		slog.Debug(msg, args...)
	case LevelWarn:
		slog.Warn(msg, args...)
	default:
		slog.Info(msg, args...)
	}
}
