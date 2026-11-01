package natsx

import (
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"go.uber.org/fx"
)

func NewPublisherModule[T any]() fx.Option {
	return fx.Module(
		"natsx_publisher",
		fx.Provide(fx.Private, NewPublisherConfig),
		fx.Provide(
			fx.Annotate(NewPublisher[T], fx.As(new(queue.Publisher[T]))),
		),
	)
}
