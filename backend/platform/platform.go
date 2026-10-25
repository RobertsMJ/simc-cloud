package platform

import (
	"log/slog"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/RobertsMJ/simc-cloud/backend/platform/admin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func New(ops ...fx.Option) *fx.App {

	newLogger := func(log *slog.Logger) fxevent.Logger {
		logger := fxevent.SlogLogger{Logger: log}
		logger.UseLogLevel(slog.LevelDebug)
		return &logger
	}

	base := []fx.Option{
		fx.StopTimeout(stopTimeout()),
		fx.WithLogger(newLogger),
		admin.Module,
	}
	return fx.New(append(base, ops...)...)
}

// termGrace allows a buffer to make sure we shut down before k8s kills us
const termGrace = time.Duration(3 * time.Second)

func stopTimeout() time.Duration {
	timeout, err := time.ParseDuration(config.OptionalEnv("TERMINATION_PERIOD", "30s"))
	if err != nil {
		slog.Error("failed to parse termination period", "err", err)
		timeout = 30 * time.Second
	}
	timeout -= termGrace
	return timeout
}
