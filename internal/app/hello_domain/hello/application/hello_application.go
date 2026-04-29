package application

import "go-clean-api-scaffold/internal/app/hello_domain/hello/application/usecase"

type HelloApplication struct {
	sayHello usecase.SayHelloUseCase
}

func NewHelloApplication(sayHello usecase.SayHelloUseCase) *HelloApplication {
	return &HelloApplication{sayHello: sayHello}
}

func (a *HelloApplication) SayHello(input *usecase.SayHelloInput) (*usecase.SayHelloOutput, error) {
	return a.sayHello.Execute(input)
}
