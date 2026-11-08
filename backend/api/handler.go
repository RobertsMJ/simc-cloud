package api

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

const BasePath = "/api"

type Request[ReqBody any] struct {
	Body ReqBody
}

type Response[RespBody any] struct {
	Body RespBody
}

func NewLogMiddleware(log *slog.Logger) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		log.Info("handling request", "method", ctx.Method(), "path", ctx.URL().Path)
		next(ctx)

	}
}

func NewAPI(config *Config, log *slog.Logger, profileHandler *ProfileHandler) http.Handler {

	router := http.NewServeMux()
	logMiddleware := NewLogMiddleware(log)

	cfg := huma.DefaultConfig(config.Name, config.Version)
	api := humago.NewWithPrefix(router, BasePath, cfg)
	api.UseMiddleware(logMiddleware)

	RegisterRoutes(api, profileHandler)
	return router
}

func RegisterRoutes(api huma.API, profileHandler *ProfileHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "parse-profile",
		Method:      http.MethodPost,
		Path:        "/profile",
	}, profileHandler.HandleProfileParse)
}
