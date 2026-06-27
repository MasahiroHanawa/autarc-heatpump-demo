package platform

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates a structured JSON logger at the given level. It uses the
// standard library's log/slog (Go 1.21+), so there is no external logging
// dependency. JSON output is trivial for log aggregators (Loki, CloudWatch,
// Datadog) to parse in production.
func NewLogger(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})
	return slog.New(handler)
}

// parseLevel maps a human-friendly level string to a slog.Level, defaulting to
// Info for unknown or empty values.
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
