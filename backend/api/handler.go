package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Request[ReqBody any] struct {
	Body ReqBody
}

type Response[RespBody any] struct {
	Body RespBody
}

func NewAPI(config *Config, profileHandler *ProfileHandler) http.Handler {
	router := http.NewServeMux()
	cfg := huma.DefaultConfig(config.Name, config.Version)
	api := humago.New(router, cfg)
	RegisterRoutes(api, profileHandler)
	return router
}

func RegisterRoutes(api huma.API, profileHandler *ProfileHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "parse-profile",
		Method:      http.MethodPost,
		Path:        "/api/profile",
	}, profileHandler.HandleProfileParse)
}
