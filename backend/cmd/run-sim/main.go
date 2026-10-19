package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform/admin"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/sim"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type SimConfig struct {
	MaxRuntime time.Duration
}

func NewSimConfig() *SimConfig {
	maxRuntime, err := time.ParseDuration(config.MustEnv("MAX_RUNTIME"))
	if err != nil {
		panic("failed to parse maxRuntime")
	}
	return &SimConfig{
		MaxRuntime: maxRuntime,
	}
}

type SimService struct {
	simulator sim.Simulator
	publisher queue.Publisher[models.SimResult]
	cfg       *SimConfig
}

var _ queue.Handler[models.SimRequest] = (*SimService)(nil)

func NewSimService(sim sim.Simulator, publisher queue.Publisher[models.SimResult], cfg *SimConfig) *SimService {
	return &SimService{
		simulator: sim,
		publisher: publisher,
		cfg:       cfg,
	}
}

func (s *SimService) Handle(ctx context.Context, input models.SimRequest) error {
	slog.Info("handling sim request", "job_id", input.JobID, "gearset_id", input.GearsetID)
	simCtx, simCancel := context.WithTimeout(ctx, s.cfg.MaxRuntime)
	defer simCancel()
	res, err := s.simulator.Run(simCtx, &input)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("sim deadline exceeded", "job_id", input.JobID, "gearset_id", input.GearsetID)
			return errors.Join(queue.ErrFatal, err)
		}
		slog.Error("sim failed", "job_id", input.JobID, "gearset_id", input.GearsetID, "error", err)
		return err
	}

	publishCtx, publishCancel := context.WithTimeout(ctx, 1*time.Second)
	defer publishCancel()
	if err := s.publisher.Publish(publishCtx, res); err != nil {
		if ctxErr := publishCtx.Err(); ctxErr != nil {

		}
		slog.Error("failed to publish result", "job_id", input.JobID, "gearset_id", input.GearsetID, "error", err)
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

func main() {
	app := fx.New(
		config.WithShutdownTimeout(),
		fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
			logger := fxevent.SlogLogger{Logger: log}
			logger.UseLogLevel(slog.LevelDebug)
			return &logger
		}),
		admin.Module,
		natsx.Module,
		fx.Provide(
			natsx.NewConsumerConfig,
			natsx.NewPublisherConfig,
			NewSimConfig,
			sim.NewSimulator,
			fx.Annotate(natsx.NewPublisher[models.SimResult], fx.As(new(queue.Publisher[models.SimResult]))),
			fx.Annotate(NewSimService, fx.As(new(queue.Handler[models.SimRequest]))),
		),
		fx.Invoke(natsx.RunConsumer[models.SimRequest]),
	)

	app.Run()
}
