package natsx

import (
	"github.com/RobertsMJ/simc-cloud/backend/platform/readiness"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"natsx",
	fx.Provide(
		NewNatsConfig,
		NewJetstream,
		NewMetrics,
		readiness.AsChecker(NewConnChecker),
	),
	fx.Invoke(registerConnMetrics),
)

type ConnChecker struct {
	conn *nats.Conn
}

func NewConnChecker(js jetstream.JetStream) *ConnChecker {
	return &ConnChecker{
		conn: js.Conn(),
	}
}

func (c *ConnChecker) Check() readiness.CheckResult {
	status := readiness.StatusUnhealthy
	if c.conn != nil && c.conn.IsConnected() {
		status = readiness.StatusHealthy
	}
	return readiness.CheckResult{
		Name:   "nats",
		Status: status,
	}
}
