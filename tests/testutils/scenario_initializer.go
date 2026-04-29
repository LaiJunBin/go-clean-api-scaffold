package testutils

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func InitContextSteps(ctx *Context) {
	api := ctx.API

	ctx.Context.Step(`^I send a "(GET|POST|PUT|PATCH|DELETE)" request to "([^"]*)":$`, api.ISendrequestTo)
	ctx.Context.Step(`^the response code should be (\d+)$`, api.TheResponseCodeShouldBe)
	ctx.Context.Step(`^the response should match json:$`, api.TheResponseShouldMatchJSON)
}

type defaultScenarioInitializer struct{}

func (s *defaultScenarioInitializer) InitializeScenario(ctx *Context) {
	api := ctx.API
	db := ctx.DB

	ctx.Context.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		fmt.Printf("[feature][hook][before] scenario=%q id=%s initializer=default\n", sc.Name, sc.Id)
		ResetDB(db)
		RunMigration(db)
		api.ResetResponse(sc)
		return ctx, nil
	})

	ctx.Context.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		fmt.Printf("[feature][hook][after] scenario=%q id=%s initializer=default err=%v\n", sc.Name, sc.Id, err)
		return ctx, err
	})

	InitContextSteps(ctx)
}

func NewDefaultScenarioInitializer() ScenarioInitializer {
	return &defaultScenarioInitializer{}
}
