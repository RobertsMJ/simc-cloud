package postgresx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(cfg *Config) (*pgxpool.Pool, error) {
	pgxConfig, err := cfg.ToPgxConfig()
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(context.Background(), pgxConfig)
}
