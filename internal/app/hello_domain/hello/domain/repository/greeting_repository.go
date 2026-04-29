package repository

import "go-clean-api-scaffold/internal/app/hello_domain/hello/domain/entity/greeting"

type GreetingRepository interface {
	Save(g *greeting.Greeting) error
	FindByName(name string) (*greeting.Greeting, error)
}
