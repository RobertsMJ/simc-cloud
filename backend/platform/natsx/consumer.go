package natsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

type ConsumerConfig struct {
	StreamID   string
	ConsumerID string
}

func NewConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		StreamID:   config.MustEnv("CONSUMER_STREAM_NAME"),
		ConsumerID: config.MustEnv("CONSUMER_NAME"),
	}
}

func RunConsumer[T any](lc fx.Lifecycle, js jetstream.JetStream, cfg *ConsumerConfig, handler queue.Handler[T]) error {
	var draining atomic.Bool
	jobCtx, jobCancel := context.WithCancel(context.Background())

	// Root consumer handling context - indicates if this consumer is still accepting messages
	var consumerCtx jetstream.ConsumeContext

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			consumer, err := js.Consumer(ctx, cfg.StreamID, cfg.ConsumerID)
			if err != nil {
				jobCancel()
				return fmt.Errorf("failed to get consumer: %w", err)
			}
			slog.Info("created consumer", "stream", cfg.StreamID, "consumer", cfg.ConsumerID)
			pulseInterval := consumer.CachedInfo().Config.AckWait / 3

			consumerCtx, err = consumer.Consume(
				func(msg jetstream.Msg) {
					// nak any pending messages while shutdown in progress
					if draining.Load() {
						msg.Nak()
						return
					}

					var data T
					if err := json.Unmarshal(msg.Data(), &data); err != nil {
						slog.Error("failed to unmarshal message", "error", err)
						msg.Term()
						return
					}

					err := invokeHandler(
						jobCtx,
						func(ctx context.Context) error {
							return handler.Handle(ctx, data)
						},
						invokeParams{
							pulseInterval: pulseInterval,
							onPulse:       msg.InProgress,
						},
					)

					switch {
					case err == nil:
						slog.Debug("message handled successfully, acknowledging")
						if err := msg.Ack(); err != nil {
							slog.Error("failed to acknowledge message", "error", err)
						}
					case jobCtx.Err() != nil:
						slog.Info("shutdown interrupted message, nak")
						msg.Nak()
					case errors.Is(err, queue.ErrRetriable):
						slog.Error("retriable error while handling message", "error", err)
						msg.NakWithDelay(5 * time.Second)
					default:
						slog.Error("fatal error while handling message", "error", err)
						msg.Term()
					}
				},
				jetstream.PullMaxMessages(1),
			)
			if err != nil {
				jobCancel()
				return fmt.Errorf("failed to consume messages: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("stopping consumer")
			defer jobCancel()
			draining.Store(true)
			consumerCtx.Drain() // no more new messages please

			select {
			case <-consumerCtx.Closed():
				slog.Info("consumer drained cleanly")
			case <-ctx.Done():
				slog.Warn("stop deadline reached with message still being processed")
			}
			return nil
		},
	})
	return nil
}

type invokeParams struct {
	pulseInterval time.Duration
	onPulse       func() error
}

func invokeHandler(ctx context.Context, handler func(context.Context) error, params invokeParams) error {
	// context termination reasons:
	//   lost connection with nats
	//   application is shutting down
	// on done, kills the pulse
	invokeCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	done := make(chan struct{})
	go func() {
		defer close(done)
		pulse := time.NewTicker(params.pulseInterval)
		defer pulse.Stop()
		for {
			select {
			case <-invokeCtx.Done():
				return
			case <-pulse.C:
				if err := params.onPulse(); err != nil {
					slog.Warn("heartbeat failed", "error", err)
					cancel(err)
					continue
				}
			}
		}
	}()
	err := handler(invokeCtx)
	cancel(nil) // stop the pulse...
	<-done      // ...and wait for it, so no InProgress lands after the Ack

	return err
}
