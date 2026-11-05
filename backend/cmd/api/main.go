package main

import (
	"github.com/RobertsMJ/simc-cloud/backend/api"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"github.com/RobertsMJ/simc-cloud/backend/profile"
	"go.uber.org/fx"
)

func BuildOptions() []fx.Option {
	return []fx.Option{
		fx.Provide(
			profile.NewService,
		),
		httpx.New(
			"api",
			api.NewHandler,
		),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
