package setup

import (
	"fmt"

	"go.uber.org/fx"

	"go-clean-api-scaffold/app"
	"go-clean-api-scaffold/app/provider"
	"go-clean-api-scaffold/internal/config"
)

func (s *setup) setupServer() {
	fmt.Println("Setting up server...")

	provides := append([]interface{}{
		config.NewTestConfig,
	}, provider.GetProvides()...)
	invokes := append(provider.GetInvokes(), app.SetupHandler)

	fx.New(
		fx.Provide(provides...),
		fx.Invoke(invokes...),
	).Run()
}
