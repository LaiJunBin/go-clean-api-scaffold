package controller

import (
	"go-clean-api-scaffold/api"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/adapter/presenter"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/application"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/application/usecase"
	sharedtypes "go-clean-api-scaffold/internal/app/shared/types"

	"github.com/gin-gonic/gin"
)

type HelloController struct {
	presenter *presenter.HelloPresenter
	app       *application.HelloApplication
	logger    sharedtypes.Logger
}

func NewHelloController(
	presenter *presenter.HelloPresenter,
	app *application.HelloApplication,
	logger sharedtypes.Logger,
) *HelloController {
	return &HelloController{
		presenter: presenter,
		app:       app,
		logger:    logger,
	}
}

func (ctl *HelloController) SayHello(ctx *gin.Context, request api.SayHelloRequestObject) (api.SayHelloResponseObject, error) {
	name := ""
	if request.Params.Name != nil {
		name = *request.Params.Name
	}

	ctl.logger.Debugf("SayHello request: name=%s", name)

	input := &usecase.SayHelloInput{Name: name}

	output, err := ctl.app.SayHello(input)
	if err != nil {
		return ctl.presenter.SayHello.Error(err), nil
	}

	ctl.logger.Debugf("SayHello succeeded: message=%s", output.Message)

	return ctl.presenter.SayHello.Output(output), nil
}
