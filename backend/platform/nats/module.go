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
	"go.uber.org/fx"
)

var Module = fx.Module(
	"nats",
	fx.Provide(
		NewNatsConfig,
		NewStream,
	),
)

type NatsConfig struct {
	URL string
}

func NewNatsConfig() *NatsConfig {
	return &NatsConfig{
		URL: config.MustEnv("QUEUE_URL"),
	}
}

type ConsumerConfig struct {
	StreamID   string
	ConsumerID string
}

func NewStream(lc fx.Lifecycle, cfg *NatsConfig) (jetstream.JetStream, error) {

	// var natsConn *nats.Conn
	// var js jetstream.JetStream
	// var err error
	slog.Debug("initializing jetstream connection")

	natsConn, err := nats.Connect(cfg.URL)
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	js, err := jetstream.New(natsConn)
	if err != nil {
		slog.Error("failed to create JetStream context", "error", err)
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			slog.Info("draining nats connection")
			natsConn.Drain()
			return nil
		},
	})
	return js, nil
}

func NewConsumer[T any](lc fx.Lifecycle, js jetstream.JetStream, cfg *ConsumerConfig, handler queue.Consumer[T]) error {
	slog.Info("creating consumer", "stream", cfg.StreamID, "consumer", cfg.ConsumerID)
	stream, err := js.Stream(context.TODO(), cfg.StreamID)
	if err != nil {
		slog.Error("failed to get stream", "error", err)
		return fmt.Errorf("failed to get stream: %w", err)
	}

	consumer, err := stream.Consumer(context.TODO(), cfg.ConsumerID)
	if err != nil {
		return fmt.Errorf("failed to get consumer: %w", err)
	}

	jobCtx, jobCancel := context.WithCancel(context.Background())
	var consumerCtx jetstream.ConsumeContext

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			consumerCtx, err = consumer.Consume(func(msg jetstream.Msg) {
				var data T
				if err := json.Unmarshal(msg.Data(), &data); err != nil {
					slog.Error("failed to unmarshal message", "error", err)
					msg.Term()
					return
				}
				if err := handler.Handle(jobCtx, data); err != nil {
					switch {
					case errors.Is(err, queue.ErrRetriable):
						slog.Error("retriable error while handling message", "error", err)
						msg.Nak()
						return
					case errors.Is(err, queue.ErrFatal):
						slog.Error("fatal error while handling message", "error", err)
						msg.Term()
						return
					default:
						slog.Error("unexpected error while handling message", "error", err)
						msg.Term()
						return
					}
				}
				slog.Info("message handled successfully, acknowledging")
				if err := msg.Ack(); err != nil {
					slog.Error("failed to acknowledge message", "error", err)
				}
			}, jetstream.PullMaxMessages(1))
			if err != nil {
				jobCancel()
				return fmt.Errorf("failed to consume messages: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("stopping consumer")
			consumerCtx.Drain()
			<-consumerCtx.Closed()
			jobCancel()
			return nil
		},
	})
	return nil
}
