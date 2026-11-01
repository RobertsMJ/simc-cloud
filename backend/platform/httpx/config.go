package httpx

import (
	"strings"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/config"
)

type Config struct {
	ID      string        `env:"-"`
	Host    string        `env:"HOST" envDefault:"0.0.0.0"`
	Port    string        `env:"PORT"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"10s"`
}

func NewConfig() (*Config, error) {
	return NewConfigWithPrefix("app", "8080")
}

func NewConfigWithPrefix(prefix string, defaultPort string) (*Config, error) {
	cfg, err := config.Parse[Config](strings.ToUpper(prefix) + "_")
	if err != nil {
		return nil, err
	}
	cfg.ID = prefix
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}
	return cfg, nil
}
