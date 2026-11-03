package admin

import (
	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"github.com/RobertsMJ/simc-cloud/backend/platform/readiness"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

const adminModule = "admin"

var Module = fx.Module(
	adminModule,
	readiness.Module,
	fx.Provide(
		fx.Annotate(
			NewPrometheusRegistry,
			fx.As(new(prometheus.Registerer)),
			fx.As(new(prometheus.Gatherer)),
		),
	),
	httpx.NewWithConfig(adminModule, NewAdminHandler, NewAdminServerConfig),
)
