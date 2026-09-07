package main

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/sim"
)

type resultPublisher interface {
	Publish(ctx context.Context, result models.SimResult) error
}

var simulator sim.Simulator
var publisher resultPublisher

func init() {
	// applog.Init()
	// cfg := LoadConfig(context.Background())
	// publisher = transport.NewPublisher[models.SimResult](context.Background(), cfg.resultsQueueURL)
	// simulator = sim.NewSimulator()
}

func handler(ctx context.Context, input models.SimRequest) error {
	res, err := simulator.Run(ctx, &input)
	if err != nil {
		slog.Error("simulation failed", "error", err)
		return fmt.Errorf("simulation failed: %w", err)
	}

	if err := publisher.Publish(ctx, res); err != nil {
		slog.Error("failed to publish result", "error", err)
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

func main() {
	fx.New(
		fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: log}
		}),
		platform.Module,
		fx.Provide(
			context.Background,
			NewSimConfig,
			sim.NewSimulator,
		),
		fx.Invoke(func(ctx context.Context, cfg SimConfig, sim sim.Simulator) {
			slog.Info("running run-sim app")
		}),
	).Run()
}
