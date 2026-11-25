package postgresx

import (
	"context"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/platform/readiness"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"postgresx",
	fx.Provide(
		NewConfig,
		NewPostgresConnection,
	),
	fx.Invoke(func(lc fx.Lifecycle, pool *pgxpool.Pool) error {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return pool.Ping(ctx)
			},
			OnStop: func(ctx context.Context) error {
				pool.Close()
				return nil
			},
		})

		return nil
	}),
)

type ConnChecker struct {
	pool *pgxpool.Pool
}

func NewConnChecker(js *pgxpool.Pool) *ConnChecker {
	return &ConnChecker{
		pool: js,
	}
}

func (c *ConnChecker) Check() readiness.CheckResult {
	res := readiness.CheckResult{
		Name:   "postgres",
		Status: readiness.StatusUnhealthy,
	}

	if c.pool == nil {
		return res
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := c.pool.Ping(ctx); err != nil {
		res.Status = readiness.StatusUnhealthy
	} else {
		res.Status = readiness.StatusHealthy
	}

	return res
}
