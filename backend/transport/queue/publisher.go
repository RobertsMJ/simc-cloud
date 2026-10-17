package queue

import "context"

type Publisher[T any] interface {
	Publish(ctx context.Context, msg T) error
}
