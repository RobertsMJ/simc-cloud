package readiness

import "go.uber.org/fx"

const readinessGroup = `group:"readiness"`

var Module = fx.Provide(fx.Annotate(NewHandler, fx.ParamTags(readinessGroup)))

func AsChecker(c any) any {
	return fx.Annotate(c, fx.As(new(Checker)), fx.ResultTags(readinessGroup))
}
