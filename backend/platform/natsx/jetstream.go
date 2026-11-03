package natsx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

func NewJetstream(lc fx.Lifecycle, cfg *NatsConfig) (jetstream.JetStream, error) {
	slog.Debug("initializing jetstream connection")

	drained := make(chan bool)
	nc, err := nats.Connect(
		cfg.Url,
		nats.MaxReconnects(-1),
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
		slog.Error("failed to connect to NATS", "error", err, "url", cfg.Url)
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

func registerConnMetrics(reg prometheus.Registerer, js jetstream.JetStream) {
	conn := js.Conn()
	reg.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Namespace: "natsx", Name: "connected"}, func() float64 {
			if conn.IsConnected() {
				return 1
			}
			return 0
		}),
		prometheus.NewCounterFunc(prometheus.CounterOpts{Namespace: "natsx", Name: "reconnects_total"}, func() float64 {
			return float64(conn.Stats().Reconnects)
		}),
	)
}
