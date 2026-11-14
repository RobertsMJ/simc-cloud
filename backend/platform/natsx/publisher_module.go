package natsx

import (
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
)

func NewPublisherModule[T any](id string) fx.Option {
	return fx.Module(
		"natsx_publisher_"+id,
		fx.Provide(fx.Private, func() (*PublisherConfig, error) {
			return NewPublisherConfig(id)
		}),
		fx.Provide(
			fx.Annotate(NewPublisher[T], fx.As(new(queue.Publisher[T]))),
		),
	)
}
