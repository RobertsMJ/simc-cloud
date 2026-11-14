package main

import (
	"go.uber.org/fx"

	"github.com/RobertsMJ/simc-cloud/backend/gearset"
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
			fx.Annotate(simcexec.New, fx.As(new(sim.Simulator))),
		),
		natsx.NewConsumerModule[models.GenerateGearsetRequest](gearset.NewService),
		natsx.NewPublisherModule[models.SimRequest]("run-sim"),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
