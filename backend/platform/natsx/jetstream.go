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

type NatsConfig struct {
	Url string
}

func NewNatsConfig() *NatsConfig {
	return &NatsConfig{
		Url: config.MustEnv("QUEUE_URL"),
	}
}

func NewJetstream(lc fx.Lifecycle, cfg *NatsConfig) (jetstream.JetStream, error) {
	slog.Debug("initializing jetstream connection")

	drained := make(chan bool)
	nc, err := nats.Connect(
		cfg.Url,
		nats.ClosedHandler(
			func(c *nats.Conn) {
				close(drained)
			},
		),
		nats.ConnectHandler(func(c *nats.Conn) {
			slog.Info("connected to nats")
		}),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			slog.Error("disconnect error", "err", err)
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			slog.Info("attempting reconnect")
		}),
		nats.ReconnectErrHandler(func(c *nats.Conn, err error) {
			slog.Error("failed to reconnect", "err", err)
		}),
		nats.LameDuckModeHandler(func(c *nats.Conn) {
			slog.Warn("lame duck detected")
		}),
	)
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("failed to create JetStream context", "error", err)
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			slog.Info("draining nats connection")
			if err := nc.Drain(); err != nil {
				slog.Error("failed to drain nats connection", "err", err)
			}
			select {
			case <-drained:
			case <-ctx.Done():
			}
			return ctx.Err()
		},
	})
	return js, nil
}
