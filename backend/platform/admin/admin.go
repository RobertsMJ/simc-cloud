package admin

import (
	"net/http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"admin",
	fx.Provide(
		NewConfig,
		NewAdminServer,
		NewLogger,
		fx.Annotate(NewServeMux, fx.ParamTags(`group:"admin"`)),
		AsAdminRoute(NewHealthCheckHandler),
		AsAdminRoute(NewReadyCheckHandler),
		AsAdminRoute(NewMetricsHandler),
	),
	fx.Invoke(func(*http.Server) {}),
)
