package api

import "github.com/RobertsMJ/simc-cloud/backend/config"

type Config struct {
	Name    string `env:"NAME" envDefault:"SimC API"`
	Version string `env:"VERSION" envDefault:"dev"`
}

func NewConfig() (*Config, error) {
	return config.Parse[Config]("API_")
}
