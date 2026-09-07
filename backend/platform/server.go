package platform

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"go.uber.org/fx"
)

type AdminServerParams struct {
	fx.In

	Config *config.Config
	Mux    *http.ServeMux
}

type Result struct {
	fx.Out

	Server *http.Server
}

func NewAdminServer(lc fx.Lifecycle, p AdminServerParams) (Result, error) {

	server := &http.Server{
		Addr:    p.Config.AdminHost + ":" + p.Config.AdminPort,
		Handler: p.Mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			slog.Info("Starting HTTP server at", "addr", server.Addr)
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
	return Result{Server: server}, nil
}

type Route interface {
	http.Handler
	Path() string
}

func AsAdminRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"admin"`),
	)
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Path(), route)
	}
	return mux
}
