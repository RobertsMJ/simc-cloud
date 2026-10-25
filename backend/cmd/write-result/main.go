package main

import (
	"github.com/RobertsMJ/simc-cloud/backend/db"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/result"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		fx.Provide(
			natsx.NewConsumerConfig,
			// db.NewResultRepository,
			fx.Annotate(db.NewResultRepository, fx.As(new(result.Repo))),
			fx.Annotate(result.NewService, fx.As(new(queue.Handler[models.SimResult]))),
		),
		fx.Invoke(natsx.RunConsumer[models.SimResult]),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
