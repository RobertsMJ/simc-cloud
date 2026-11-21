package main

import (
	"github.com/RobertsMJ/simc-cloud/backend/db"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/result"
	"go.uber.org/fx"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		fx.Provide(
			fx.Annotate(db.NewResultRepository, fx.As(new(result.Repo))),
		),
		natsx.NewConsumerModule[models.SimResult](result.NewService),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
