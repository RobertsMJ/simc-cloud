package postgresx

import (
	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL string `env:"DATABASE_URL,required"`
}

func (c Config) ToPgxConfig() (*pgxpool.Config, error) {
	return pgxpool.ParseConfig(c.URL)
}

func NewConfig() (*Config, error) {
	return config.Parse[Config]("")
}
