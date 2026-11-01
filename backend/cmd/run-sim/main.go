package main

import (
	"go.uber.org/fx"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/sim"
	"github.com/RobertsMJ/simc-cloud/backend/sim/simcexec"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		fx.Provide(
			simcexec.NewConfig,
			sim.NewConfig,
			fx.Annotate(simcexec.New, fx.As(new(sim.Simulator))),
		),
		natsx.NewPublisherModule[models.SimResult](),
		natsx.NewConsumerModule[models.SimRequest](sim.NewService),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
