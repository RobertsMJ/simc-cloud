package main

import (
	"context"
	"net/http"

	genapi "github.com/RobertsMJ/simc-cloud/backend/api/gen"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"go.uber.org/fx"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		httpx.New(
			"api",
			// TODO: Obviously this gets pulled out to the api module
			func() http.Handler {
				handler := genapi.NewStrictHandler(Server{}, []genapi.StrictMiddlewareFunc{})
				r := http.NewServeMux()
				return genapi.HandlerFromMux(handler, r)
			},
		),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}

// TODO: Obviously this gets removed
type Server struct{}

var _ genapi.StrictServerInterface = (*Server)(nil)

func (h Server) GetApiTest(ctx context.Context, request genapi.GetApiTestRequestObject) (genapi.GetApiTestResponseObject, error) {
	return genapi.GetApiTest200JSONResponse(123), nil
}
