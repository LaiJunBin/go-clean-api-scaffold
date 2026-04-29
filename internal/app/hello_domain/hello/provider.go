package hello

import (
	"go-clean-api-scaffold/internal/app/hello_domain/hello/adapter/controller"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/adapter/presenter"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/application"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/application/usecase"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/infrastructure/repository"
	"go-clean-api-scaffold/internal/app/shared/types/provider"
)

func NewProvider() *provider.Provider {
	return provider.New(provider.Config{
		Infrastructures: []interface{}{
			repository.NewGreetingRepository,
		},
		Controllers: []interface{}{
			controller.NewHelloController,
		},
		Presenters: []interface{}{
			presenter.NewHelloPresenter,
		},
		Usecases: []interface{}{
			usecase.NewSayHelloUseCase,
		},
		Applications: []interface{}{
			application.NewHelloApplication,
		},
	})
}
