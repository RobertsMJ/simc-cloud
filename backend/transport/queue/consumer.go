package queue

import "context"

type Consumer[T any] interface {
	Handle(ctx context.Context, msg T) error
}
