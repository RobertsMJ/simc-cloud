package natsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher[T any] struct {
	logger  *slog.Logger
	js      jetstream.JetStream
	cfg     *PublisherConfig
	metrics publisherMetrics
}

func NewPublisher[T any](logger *slog.Logger, cfg *PublisherConfig, js jetstream.JetStream, metrics *Metrics) *Publisher[T] {
	return &Publisher[T]{
		logger:  logger.With("stream", cfg.StreamName, "subject", cfg.Subject),
		js:      js,
		cfg:     cfg,
		metrics: metrics.forPublisher(cfg),
	}
}

func (p *Publisher[T]) Publish(ctx context.Context, msg T) error {
	start := time.Now()
	err := p.publish(ctx, msg)
	p.metrics.published.WithLabelValues(string(handlerOutcome(err))).Observe(time.Since(start).Seconds())
	return err
}

func (p *Publisher[T]) publish(ctx context.Context, msg T) error {
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
