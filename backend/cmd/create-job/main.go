package main

import (
	"github.com/RobertsMJ/simc-cloud/backend/db"
	"github.com/RobertsMJ/simc-cloud/backend/job"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		fx.Provide(
			natsx.NewPublisherConfig,
			fx.Annotate(db.NewJobRepository, fx.As(new(job.Repo))),
			fx.Annotate(natsx.NewPublisher[models.SimRequest], fx.As(new(queue.Publisher[models.SimRequest]))),
		),
	}
}

// TODO: stand up http POST handler
func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
