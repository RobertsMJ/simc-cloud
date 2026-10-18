package natsx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"natsx",
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

type PublisherConfig struct {
	StreamID string
	Subject  string
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
