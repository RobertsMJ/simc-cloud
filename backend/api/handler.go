package api

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Handler interface {
	RegisterHandler(api huma.API)
}

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

func NewAPI(config *Config, log *slog.Logger, handlers []Handler) http.Handler {

	router := http.NewServeMux()
	logMiddleware := NewLogMiddleware(log)

	cfg := huma.DefaultConfig(config.Name, config.Version)
	api := humago.NewWithPrefix(router, BasePath, cfg)
	api.UseMiddleware(logMiddleware)

	for _, handler := range handlers {
		handler.RegisterHandler(api)
	}

	return router
}
