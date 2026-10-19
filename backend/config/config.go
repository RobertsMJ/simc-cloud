package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.uber.org/fx"
)

func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}

func OptionalEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type Config struct {
	AdminPort   string
	AdminHost   string
	StopTimeout time.Duration
}

// termGrace allows a buffer to make sure we shut down before k8s kills us
const termGrace = time.Duration(3 * time.Second)

func NewAppConfig() *Config {
	StopTimeout, err := time.ParseDuration(OptionalEnv("TERMINATION_PERIOD", "30s"))
	if err != nil {
		slog.Error("failed to parse termination period", "err", err)
		StopTimeout = 30 * time.Second
	}
	StopTimeout -= termGrace

	return &Config{
		AdminPort:   OptionalEnv("ADMIN_PORT", "8081"),
		AdminHost:   OptionalEnv("ADMIN_HOST", "0.0.0.0"),
		StopTimeout: StopTimeout,
	}
}

func WithShutdownTimeout() fx.Option {
	cfg := NewAppConfig()
	return fx.StopTimeout(cfg.StopTimeout)
}
