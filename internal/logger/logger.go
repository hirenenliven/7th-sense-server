package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
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

// Init initializes the global logger with JSON output to both terminal and file.
func Init() {
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // MB before rotating
		MaxBackups: 5,  // old files to keep
		MaxAge:     30, // days to keep
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	log = slog.New(slog.NewJSONHandler(multiWriter, opts))
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
