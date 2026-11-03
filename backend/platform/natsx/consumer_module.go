package natsx

import (
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go/jetstream"
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
			func(
				logger *slog.Logger,
				config *ConsumerConfig,
				js jetstream.JetStream,
				handler queue.Handler[T],
				shutdowner fx.Shutdowner,
				metrics *Metrics,
			) *Consumer[T] {
				return NewConsumer(logger, config, js, handler, func() error {
					return shutdowner.Shutdown(fx.ExitCode(1))
				}, metrics)
			},
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
