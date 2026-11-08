package simcexec

import (
	"github.com/RobertsMJ/simc-cloud/backend/config"
)

type Config struct {
	BinPath string `env:"SIM_BINARY_PATH" envDefault:"/app/simc"`
	Threads int    `env:"-"`
}

func NewConfig() (*Config, error) {
	cfg, err := config.Parse[Config]("")
	if err != nil {
		return nil, err
	}
	cfg.Threads = 1
	return cfg, nil
}
