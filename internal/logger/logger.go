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

	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	log = slog.New(slog.NewJSONHandler(multiWriter, handlerOptions))
	slog.SetDefault(log)
}

// Log logs a message with the specified level and key-value pairs.
func Log(level Level, message string, args ...any) {
	switch level {
	case LevelInfo:
		slog.Info(message, args...)
	case LevelError:
		slog.Error(message, args...)
	case LevelDebug:
		slog.Debug(message, args...)
	case LevelWarn:
		slog.Warn(message, args...)
	default:
		slog.Info(message, args...)
	}
}
