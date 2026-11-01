package httpx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"go.uber.org/fx"
)

const httpxModule = "httpx"

// New produces an fx module with the provided handler and the default APP configuration
func New(
	modulePrefix string,
	handler any,
) fx.Option {
	return NewWithConfig(modulePrefix, handler, NewConfig)
}

// NewWithConfig produces an fx module with the provided handler and configuration
func NewWithConfig(
	modulePrefix string,
	handler any,
	config func() (*Config, error),
) fx.Option {
	return fx.Module(
		fmt.Sprintf("%s_%s", modulePrefix, httpxModule),
		fx.Provide(
			fx.Private,
			config,
			fx.Annotate(handler, fx.As(new(http.Handler))),
			NewServer,
		),
		fx.Invoke(func(lc fx.Lifecycle, sd fx.Shutdowner, logger *slog.Logger, server *Server) {
			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						if err := server.Start(); err != nil {
							return err
						}

						go func() {
							err := <-server.Done()
							if err != nil {
								logger.Warn("server crashed, triggering shutdown", "err", err)
								if err = sd.Shutdown(fx.ExitCode(1)); err != nil {
									logger.Error("failed to trigger fx shutdown", "err", err)
								}
							}
						}()
						return nil
					},
					OnStop: server.Stop,
				},
			)
		}),
	)
}
