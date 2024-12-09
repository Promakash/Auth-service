package log

import (
	"log/slog"
	"os"
)

type Config struct {
	level  string `yaml:"level"`
	format string `yaml:"format"`
}

func NewLogger(level string, format string) *slog.Logger {
	var handler slog.Handler

	switch format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, nil)
	default:
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}

	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return slog.New(handler).With(logLevel)
}
