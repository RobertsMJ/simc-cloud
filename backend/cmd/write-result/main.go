package main

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform/admin"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/sim"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type SimResultService struct {
}

func NewSimResultService() *SimResultService {
	return &SimResultService{}
}

var _ queue.Handler[models.SimResult] = (*SimResultService)(nil)

func (s *SimResultService) Handle(ctx context.Context, result models.SimResult) error {
	slog.Info("handling sim result", "result", result)
	return nil
}

func main() {
	fx.New(
		fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
			logger := fxevent.SlogLogger{Logger: log}
			logger.UseLogLevel(slog.LevelDebug)
			return &logger
		}),
		admin.Module,
		natsx.Module,
		fx.Provide(
			natsx.NewConsumerConfig,
			sim.NewSimulator,
			fx.Annotate(NewSimResultService, fx.As(new(queue.Handler[models.SimResult]))),
		),
		fx.Invoke(natsx.RunConsumer[models.SimResult]),
	).Run()
}
