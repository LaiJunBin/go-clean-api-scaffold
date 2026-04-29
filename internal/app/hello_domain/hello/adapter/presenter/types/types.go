package types

import (
	"go-clean-api-scaffold/api"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/application/usecase"
	sharedpresenter "go-clean-api-scaffold/internal/app/shared/types/presenter"
)

type SayHelloPresenterType = sharedpresenter.Type[*usecase.SayHelloOutput, api.SayHelloResponseObject]
