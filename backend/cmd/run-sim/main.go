package main

import (
	"go.uber.org/fx"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/sim"
	"github.com/RobertsMJ/simc-cloud/backend/sim/simcexec"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		fx.Provide(
			natsx.NewConsumerConfig,
			natsx.NewPublisherConfig,
			simcexec.NewConfig,
			sim.NewConfig,
			fx.Annotate(simcexec.New, fx.As(new(sim.Simulator))),
			fx.Annotate(natsx.NewPublisher[models.SimResult], fx.As(new(queue.Publisher[models.SimResult]))),
			fx.Annotate(sim.NewService, fx.As(new(queue.Handler[models.SimRequest]))),
		),
		fx.Invoke(natsx.RunConsumer[models.SimRequest]),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
