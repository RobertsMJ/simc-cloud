package natsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type PublisherConfig struct {
	StreamID string
	Subject  string
}

func NewPublisherConfig() *PublisherConfig {
	return &PublisherConfig{
		StreamID: config.MustEnv("PUBLISHER_STREAM_NAME"),
		Subject:  config.MustEnv("PUBLISHER_SUBJECT"),
	}
}

type Publisher[T any] struct {
	js  jetstream.JetStream
	cfg *PublisherConfig
}

func NewPublisher[T any](js jetstream.JetStream, cfg *PublisherConfig) *Publisher[T] {
	return &Publisher[T]{
		js:  js,
		cfg: cfg,
	}
}

func (p *Publisher[T]) Publish(ctx context.Context, msg T) error {
	slog.Info("stub publishing result", "result", msg)

	res, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal message", "error", err)
		return errors.Join(queue.ErrFatal, fmt.Errorf("failed to marshal message: %w", err))
	}

	ack, err := p.js.Publish(ctx, p.cfg.Subject, res)
	if err != nil {
		slog.Error("failed to publish message", "error", err)
		switch {
		case errors.Is(err, nats.ErrMaxPayload):
			return errors.Join(queue.ErrFatal)
		}
		return err
	}

	slog.Info("message published successfully", "ack", ack)
	return nil
}
