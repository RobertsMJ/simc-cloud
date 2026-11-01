package natsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumer[T any] struct {
	logger *slog.Logger
	config *ConsumerConfig

	onTerminalError func() error
	js              jetstream.JetStream
	sub             jetstream.ConsumeContext
	handler         queue.Handler[T]

	abort             context.CancelFunc
	heartbeatInterval time.Duration

	draining atomic.Bool
}

func NewConsumer[T any](
	logger *slog.Logger,
	config *ConsumerConfig,
	js jetstream.JetStream,
	handler queue.Handler[T],
	onTerminalError func() error,
) *Consumer[T] {
	return &Consumer[T]{
		onTerminalError: onTerminalError,
		config:          config,
		js:              js,
		handler:         handler,
		logger:          logger.With("stream", config.StreamName, "consumer", config.DurableName),
	}
}

func (c *Consumer[T]) Start(ctx context.Context) error {
	consumer, err := c.js.Consumer(ctx, c.config.StreamName, c.config.DurableName)
	if err != nil {
		return fmt.Errorf("failed to bind consumer: %w", err)
	}
	c.heartbeatInterval = consumer.CachedInfo().Config.AckWait * 30 / 100
	c.logger.Info("bound to consumer")

	var abortCtx context.Context
	abortCtx, c.abort = context.WithCancel(context.Background())

	c.sub, err = consumer.Consume(
		func(msg jetstream.Msg) {
			c.onMsg(abortCtx, msg)
		},
		jetstream.PullMaxMessages(c.config.MaxMessages),
		jetstream.ConsumeErrHandler(c.onConsumeErr),
	)
	if err != nil {
		c.abort()
		return fmt.Errorf("failed to consume messages: %w", err)
	}
	return nil
}

func (c *Consumer[T]) onConsumeErr(_ jetstream.ConsumeContext, err error) {
	if !errors.Is(err, jetstream.ErrConsumerDeleted) && !errors.Is(err, jetstream.ErrBadRequest) {
		c.logger.Warn("consume error", "error", err)
		return
	}
	c.logger.Error("terminal consume error, shutting down", "error", err)
	if err := c.onTerminalError(); err != nil {
		c.logger.Error("failed to trigger shutdown", "error", err)
	}
}

func (c *Consumer[T]) Stop(ctx context.Context) error {
	c.logger.Info("stopping consumer")
	defer c.abort()
	c.draining.Store(true)
	c.sub.Drain() // no more new messages please

	select {
	case <-c.sub.Closed():
		c.logger.Info("consumer drained cleanly")
	case <-ctx.Done():
		c.logger.Warn("stop deadline reached with message still being processed")
	}
	return nil
}

func (c *Consumer[T]) onMsg(ctx context.Context, msg jetstream.Msg) {

	// nak any pending messages while shutdown in progress
	if c.draining.Load() {
		msg.Nak()
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var payload T
	if err := json.Unmarshal(msg.Data(), &payload); err != nil {
		c.logger.Error("failed to unmarshal message", "error", err)
		msg.TermWithReason("failed to unmarshal message: " + err.Error())
		return
	}

	stopHeartbeat := c.runHeartbeat(ctx, cancel, msg)
	err := c.handler.Handle(ctx, payload)
	stopHeartbeat()

	switch {
	case err == nil:
		c.logger.Debug("message handled successfully, acknowledging")
		if err := msg.Ack(); err != nil {
			c.logger.Error("failed to acknowledge message", "error", err)
		}
	case errors.Is(err, queue.ErrFatal):
		c.logger.Error("fatal error while handling message", "error", err)
		msg.TermWithReason(err.Error())
	default:
		c.logger.Error("retriable error while handling message", "error", err)
		msg.Nak()
	}
}

func (c *Consumer[T]) runHeartbeat(ctx context.Context, cancel context.CancelFunc, msg jetstream.Msg) func() {

	done := make(chan struct{})
	go func() {
		defer close(done)
		heartbeat := time.NewTicker(c.heartbeatInterval)
		defer heartbeat.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-heartbeat.C:
				if err := msg.InProgress(); err != nil {
					c.logger.Warn("heartbeat failed", "error", err)
				}
			}
		}
	}()

	return func() {
		cancel()
		<-done
	}
}
