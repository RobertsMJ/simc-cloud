package config

import (
	"fmt"
	"os"
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
	AdminPort string
	AdminHost string
}

func NewConfig() *Config {
	return &Config{
		AdminPort: OptionalEnv("ADMIN_PORT", "8081"),
		AdminHost: OptionalEnv("ADMIN_HOST", "127.0.0.1"),
	}
}
