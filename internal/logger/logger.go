package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Init makes slog the default logger, writing JSON to the terminal and to
// logs/app.log (rotated). After Init, log with slog.Info, slog.Error, etc.
func Init() {
	InitFile("logs/app.log")
}

// InitFile is Init with another log file. Every process needs its own file:
// lumberjack rotation cannot be shared between processes.
func InitFile(path string) {
	logFile := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // MB before rotating
		MaxBackups: 5,  // old files to keep
		MaxAge:     30, // days to keep
		Compress:   true,
	}

	output := io.MultiWriter(os.Stdout, logFile)
	slog.SetDefault(slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}
