package sim

import (
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/config"
)

type Config struct {
	MaxRuntime time.Duration `env:"MAX_RUNTIME,required"`
}

func NewConfig() (*Config, error) {
	return config.Parse[Config]("")
}
