package jsonpresenter

import (
	"errors"

	"go-clean-api-scaffold/api"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/adapter/presenter/types"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/application/usecase"
)

type SayHelloPresenter struct{}

func NewSayHelloPresenter() types.SayHelloPresenterType {
	return &SayHelloPresenter{}
}

func (p *SayHelloPresenter) Output(output *usecase.SayHelloOutput) api.SayHelloResponseObject {
	return api.SayHello200JSONResponse{
		Success: true,
		Data:    api.HelloGreeting{Message: output.Message},
	}
}

func (p *SayHelloPresenter) Error(err error) api.SayHelloResponseObject {
	switch {
	case errors.Is(err, usecase.ErrNameRequired):
		return api.SayHello400JSONResponse{
			Success: false,
			Message: "Bad Request",
			Details: []string{err.Error()},
		}
	default:
		return api.SayHello500JSONResponse{
			Success: false,
			Message: "Internal Server Error",
		}
	}
}
