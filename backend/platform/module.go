package platform

import (
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"admin",
	fx.Provide(
		config.NewConfig,
		NewAdminServer,
		NewLogger,
		fx.Annotate(NewServeMux, fx.ParamTags(`group:"admin"`)),
		AsAdminRoute(NewHealthCheckHandler),
		AsAdminRoute(NewReadyCheckHandler),
		AsAdminRoute(NewMetricsHandler),
	),
	fx.Invoke(func(*http.Server) {}),
)
