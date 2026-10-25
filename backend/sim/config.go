package sim

import (
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/config"
)

type Config struct {
	MaxRuntime time.Duration
}

func NewConfig() *Config {
	maxRuntime, err := time.ParseDuration(config.MustEnv("MAX_RUNTIME"))
	if err != nil {
		panic("failed to parse maxRuntime")
	}
	return &Config{
		MaxRuntime: maxRuntime,
	}
}
