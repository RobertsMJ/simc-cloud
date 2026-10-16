package main

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/fx"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform/admin"
	natsx "github.com/RobertsMJ/simc-cloud/backend/platform/nats"
	"github.com/RobertsMJ/simc-cloud/backend/sim"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type StubPublisher struct{}

func NewStubPublisher() *StubPublisher {
	return &StubPublisher{}
}

func (p *StubPublisher) Publish(ctx context.Context, msg models.SimResult) error {
	slog.Info("stub publishing result", "result", msg)
	return nil
}

type SimService struct {
	simulator sim.Simulator
	publisher sim.ResultPublisher
}

var _ queue.Consumer[models.SimRequest] = (*SimService)(nil)

func NewSimService(sim sim.Simulator, publisher sim.ResultPublisher) *SimService {
	return &SimService{
		simulator: sim,
		publisher: publisher,
	}
}

func (s *SimService) Handle(ctx context.Context, input models.SimRequest) error {
	slog.Info("handling sim request", "input", input)
	res, err := s.simulator.Run(ctx, &input)
	if err != nil {
		slog.Error("simulation failed", "error", err)
		return fmt.Errorf("simulation failed: %w", err)
	}

	slog.Info("simulation completed", "result", res)
	if err := s.publisher.Publish(ctx, res); err != nil {
		slog.Error("failed to publish result", "error", err)
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

func NewConsumerConfig() *natsx.ConsumerConfig {
	return &natsx.ConsumerConfig{
		StreamID:   config.MustEnv("STREAM_ID"),
		ConsumerID: config.MustEnv("CONSUMER_ID"),
	}
}

func main() {
	app := fx.New(
		// fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
		// 	return &fxevent.SlogLogger{Logger: log}
		// }),
		admin.Module,
		natsx.Module,
		fx.Provide(
			NewConsumerConfig,
			sim.NewSimulator,
			fx.Annotate(NewStubPublisher, fx.As(new(sim.ResultPublisher))),
			fx.Annotate(NewSimService, fx.As(new(queue.Consumer[models.SimRequest]))),
		),
		fx.Invoke(natsx.NewConsumer[models.SimRequest]),
	)

	app.Run()
}
