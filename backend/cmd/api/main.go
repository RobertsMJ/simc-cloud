package main

import (
	"github.com/RobertsMJ/simc-cloud/backend/api"
	"github.com/RobertsMJ/simc-cloud/backend/db"
	"github.com/RobertsMJ/simc-cloud/backend/gearset"
	"github.com/RobertsMJ/simc-cloud/backend/job"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/RobertsMJ/simc-cloud/backend/platform/postgresx"
	"github.com/RobertsMJ/simc-cloud/backend/profile"
	"go.uber.org/fx"
)

func AsApiHandler(h any) any {
	return fx.Annotate(
		h,
		fx.As(new(api.Handler)),
		fx.ResultTags(`group:"api"`),
	)
}

func BuildOptions() []fx.Option {
	return []fx.Option{
		natsx.Module,
		postgresx.Module,
		fx.Provide(
			api.NewConfig,
			job.NewConfig,
			fx.Annotate(
				db.NewJobRepository,
				fx.As(new(job.Repo)),
			),
			fx.Annotate(
				profile.NewService,
				fx.As(new(api.ProfileParser)),
			),
			fx.Annotate(
				gearset.NewService,
				fx.As(new(job.GearsetCounter)),
			),
			fx.Annotate(
				job.NewService,
				fx.As(new(api.JobService)),
			),
			AsApiHandler(api.NewProfileHandler),
			AsApiHandler(api.NewJobHandler),
			fx.Annotate(
				api.NewAPI,
				fx.ParamTags(``, ``, `group:"api"`),
			),
		),
		natsx.NewPublisherModule[models.GenerateGearsets]("blacksmith"),
		httpx.New("api"),
	}
}

func main() {
	platform.New(
		BuildOptions()...,
	).Run()
}
