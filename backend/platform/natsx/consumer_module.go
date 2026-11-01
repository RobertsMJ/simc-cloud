package natsx

import (
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
)

func NewConsumerModule[T any](
	handler any,
) fx.Option {
	return fx.Module(
		"natsx_consumer",
		fx.Provide(
			fx.Private,
			NewConsumerConfig,
			fx.Annotate(handler, fx.As(new(queue.Handler[T]))),
			NewConsumer[T],
		),
		fx.Invoke(
			func(lc fx.Lifecycle, consumer *Consumer[T]) {
				lc.Append(fx.Hook{
					OnStart: consumer.Start,
					OnStop:  consumer.Stop,
				})
			},
		),
	)
}
