package natsx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"natsx",
	fx.Provide(
		NewNatsConfig,
		NewJetstream,
	),
)
