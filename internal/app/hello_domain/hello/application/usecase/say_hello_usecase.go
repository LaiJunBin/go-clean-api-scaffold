package usecase

import (
	"go-clean-api-scaffold/internal/app/hello_domain/hello/domain/entity/greeting"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/domain/repository"
	sharedtypes "go-clean-api-scaffold/internal/app/shared/types"
)

type SayHelloInput struct {
	Name string
}

type SayHelloOutput struct {
	Message string
}

type SayHelloUseCase interface {
	Execute(input *SayHelloInput) (*SayHelloOutput, error)
}

type sayHelloUseCase struct {
	greetingRepository repository.GreetingRepository
	logger             sharedtypes.Logger
}

func NewSayHelloUseCase(
	greetingRepository repository.GreetingRepository,
	logger sharedtypes.Logger,
) SayHelloUseCase {
	return &sayHelloUseCase{
		greetingRepository: greetingRepository,
		logger:             logger,
	}
}

func (uc *sayHelloUseCase) Execute(input *SayHelloInput) (*SayHelloOutput, error) {
	uc.logger.Infof("SayHelloUseCase: name=%s", input.Name)

	if input.Name == "" {
		return nil, ErrNameRequired
	}

	g, err := greeting.New(input.Name)
	if err != nil {
		uc.logger.Errorf("SayHelloUseCase: failed to create greeting entity: %v", err)
		return nil, ErrNameRequired
	}

	if err := uc.greetingRepository.Save(g); err != nil {
		uc.logger.Errorf("SayHelloUseCase: failed to save greeting: %v", err)
		return nil, ErrGreetingFailed
	}

	return &SayHelloOutput{Message: g.GetMessage()}, nil
}
