// Package platform holds cross-cutting concerns (configuration, logging)
// that the rest of the application depends on but that contain no business
// logic. It must not import any other internal package.
package platform

import (
	"os"
	"time"
)

// Config holds all runtime configuration for the API server. Values are read
// from environment variables with sensible defaults, so the app runs out of
// the box locally but is fully configurable in any deployment environment.
type Config struct {
	Port            string        // HTTP port to listen on
	LogLevel        string        // "debug" | "info" | "warn" | "error"
	ShutdownTimeout time.Duration // grace period for in-flight requests on shutdown
}

// Load reads configuration from the environment, applying a default for any
// variable that is unset or empty.
func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

// getEnv returns the value of the environment variable named by key, or
// fallback if the variable is unset or empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvDuration parses an environment variable as a Go duration (e.g. "5s",
// "2m"), falling back to the default if it is unset or invalid.
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
