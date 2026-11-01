package natsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher[T any] struct {
	logger *slog.Logger
	js     jetstream.JetStream
	cfg    *PublisherConfig
}

func NewPublisher[T any](logger *slog.Logger, cfg *PublisherConfig, js jetstream.JetStream) *Publisher[T] {
	return &Publisher[T]{
		logger: logger.With("stream", cfg.StreamName, "subject", cfg.Subject),
		js:     js,
		cfg:    cfg,
	}
}

func (p *Publisher[T]) Publish(ctx context.Context, msg T) error {
	data, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error("failed to marshal message", "error", err)
		return errors.Join(queue.ErrFatal, fmt.Errorf("failed to marshal message: %w", err))
	}

	ack, err := p.js.Publish(ctx, p.cfg.Subject, data)
	if err != nil {
		p.logger.Error("failed to publish message", "error", err)
		if errors.Is(err, nats.ErrMaxPayload) {
			return errors.Join(queue.ErrFatal, err)
		}
		return err
	}

	p.logger.Debug("message published", "sequence", ack.Sequence)
	return nil
}
