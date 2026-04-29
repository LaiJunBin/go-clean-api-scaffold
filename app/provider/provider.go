package provider

import (
	"go-clean-api-scaffold/api"
	"go-clean-api-scaffold/app"
	"go-clean-api-scaffold/database"
	hello "go-clean-api-scaffold/internal/app/hello_domain/hello"
	"go-clean-api-scaffold/internal/app/shared"
	"go-clean-api-scaffold/internal/app/shared/types/provider"

	"go.uber.org/fx"
)

var (
	Databases = []interface{}{
		database.DatabaseProvider,
	}

	Providers = []*provider.Provider{
		shared.NewProvider(),
		hello.NewProvider(),
	}

	Invokes = []interface{}{}
)

func collectMiddlewares(middlewares []api.MiddlewareFunc) []api.MiddlewareFunc {
	return middlewares
}

func GetProvides() []interface{} {
	provides := []interface{}{
		fx.Annotate(collectMiddlewares, fx.ParamTags(`group:"middleware"`)),
		fx.Annotate(app.NewServer, fx.ParamTags(`group:"controller"`)),
	}

	provides = append(provides, Databases...)

	for _, p := range Providers {
		controllers := []interface{}{}
		for _, c := range p.Config.Controllers {
			controllers = append(controllers, fx.Annotate(c, fx.As(new(app.Controller)), fx.ResultTags(`group:"controller"`)))
		}

		provides = append(provides, p.Config.Infrastructures...)
		provides = append(provides, p.Config.Presenters...)
		provides = append(provides, p.Config.Services...)
		provides = append(provides, p.Config.Usecases...)
		provides = append(provides, p.Config.Applications...)
		provides = append(provides, p.Config.Factories...)
		provides = append(provides, p.Config.Providers...)
		provides = append(provides, controllers...)

		Invokes = append(Invokes, p.Config.Handlers...)
		Invokes = append(Invokes, p.Config.Invokes...)
	}

	return provides
}

func GetInvokes() []interface{} {
	return Invokes
}
