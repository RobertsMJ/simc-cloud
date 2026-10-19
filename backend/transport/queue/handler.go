package queue

import "context"

type Handler[T any] interface {
	Handle(ctx context.Context, msg T) error
}
