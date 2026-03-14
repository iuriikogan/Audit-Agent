// Package logger provides logger.go implementation.
//
// Rationale: This module is designed to encapsulate domain-specific logic,
// ensuring strict separation of concerns within the multi-agent CRA architecture.
// Terminology: CRA (Cyber Resilience Act), GCP (Google Cloud Platform), Agent (Autonomous AI actor).
// Measurability: Ensures code maintainability and testability by isolating discrete workflow steps.
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
