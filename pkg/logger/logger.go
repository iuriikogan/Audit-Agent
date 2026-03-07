package logger

import (
	"log/slog"
	"os"
)

// Setup initializes the global logger with the specified options.
// It configures structured logging (JSON) and sets the log level.
func Setup(level string) {
	var logLevel slog.Level
	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
		// Add source file information for easier debugging
		AddSource: true,
	}

	// Use JSON handler for structured logging, suitable for cloud aggregation
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}
