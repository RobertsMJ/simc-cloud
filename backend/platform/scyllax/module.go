package scyllax

import (
	"context"

	"github.com/scylladb/gocqlx/v3"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"scyllax",
	fx.Provide(
		NewScyllaConfig,
		NewScyllaConnection,
	),
	fx.Invoke(func(lc fx.Lifecycle, session *gocqlx.Session) error {
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				session.Close()
				return nil
			},
		})
		return nil
	}),
)
