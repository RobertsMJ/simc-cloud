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

	base := DefaultOptions()
	return fx.New(append(base, ops...)...)
}

func DefaultOptions() []fx.Option {
	newLogger := func(log *slog.Logger) fxevent.Logger {
		logger := fxevent.SlogLogger{Logger: log}
		logger.UseLogLevel(slog.LevelDebug)
		return &logger
	}

	return []fx.Option{
		stopTimeout(),
		fx.Provide(NewLogger),
		fx.WithLogger(newLogger),
		admin.Module,
	}
}

// termGrace allows a buffer to make sure we shut down before k8s kills us
const termGrace = 3 * time.Second

type Config struct {
	TerminationPeriod time.Duration `env:"TERMINATION_PERIOD" envDefault:"30s"`
}

func stopTimeout() fx.Option {
	cfg, err := config.Parse[Config]("")
	if err != nil {
		return fx.Error(err)
	}
	return fx.StopTimeout(max(cfg.TerminationPeriod-termGrace, 0))
}
