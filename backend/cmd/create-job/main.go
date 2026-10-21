package main

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform/admin"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type CreateJobService struct {
}

func NewSimResultService() *CreateJobService {
	return &CreateJobService{}
}

var _ queue.Publisher[models.SimRequest] = (*CreateJobService)(nil)

func (s *CreateJobService) Publish(ctx context.Context, sim models.SimRequest) error {
	slog.Info("publishing sim request", "result", sim)
	return nil
}

// TODO: stand up http POST handler
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
			natsx.NewPublisherConfig,
			fx.Annotate(natsx.NewPublisher[models.SimRequest], fx.As(new(queue.Publisher[models.SimRequest]))),
		),
	).Run()
}
