package main

import (
	"go-clean-api-scaffold/app"
	"go-clean-api-scaffold/app/provider"
	"go-clean-api-scaffold/internal/config"

	"go.uber.org/fx"
)

func main() {
	provides := append([]interface{}{
		config.NewConfig,
	}, provider.GetProvides()...)
	invokes := append(provider.GetInvokes(), app.SetupHandler)

	fx.New(
		fx.Provide(provides...),
		fx.Invoke(invokes...),
	).Run()
}
