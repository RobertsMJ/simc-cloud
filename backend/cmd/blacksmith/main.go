package main

import (
	"go.uber.org/fx"

	"github.com/RobertsMJ/simc-cloud/backend/gearset"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		fx.Provide(
			fx.Annotate(
				gearset.NewService,
				fx.As(new(gearset.GearsetGenerator)),
			),
		),
		natsx.NewConsumerModule[models.GenerateGearsetsRequest](gearset.NewHandler),
		natsx.NewPublisherModule[models.SimRequest]("target-dummy"),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
