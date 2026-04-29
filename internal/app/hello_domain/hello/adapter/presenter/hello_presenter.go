package presenter

import (
	jsonpresenter "go-clean-api-scaffold/internal/app/hello_domain/hello/adapter/presenter/json"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/adapter/presenter/types"
)

type HelloPresenter struct {
	SayHello types.SayHelloPresenterType
}

func NewHelloPresenter() *HelloPresenter {
	return &HelloPresenter{
		SayHello: jsonpresenter.NewSayHelloPresenter(),
	}
}
