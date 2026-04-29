package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	hellomocks "go-clean-api-scaffold/internal/app/hello_domain/hello/mocks"
	sharedmocks "go-clean-api-scaffold/internal/app/shared/mocks"

	"go-clean-api-scaffold/internal/app/hello_domain/hello/application/usecase"
	"go-clean-api-scaffold/internal/app/hello_domain/hello/domain/entity/greeting"
)

// bypassLogger configures the mock to silently accept logger calls from use cases.
func bypassLogger(t *testing.T) *sharedmocks.Logger {
	l := sharedmocks.NewLogger(t)
	for _, method := range []string{"Debugf", "Infof", "Warnf", "Errorf", "Fatalf", "Panicf"} {
		l.On(method, mock.Anything).Maybe()
		l.On(method, mock.Anything, mock.Anything).Maybe()
	}
	return l
}

func TestSayHelloUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := hellomocks.NewGreetingRepository(t)

		mockRepo.On("Save", mock.MatchedBy(func(g *greeting.Greeting) bool {
			return g.GetName() == "Alice" && g.GetMessage() == "Hello, Alice!"
		})).Return(nil).Once()

		uc := usecase.NewSayHelloUseCase(mockRepo, bypassLogger(t))
		output, err := uc.Execute(&usecase.SayHelloInput{Name: "Alice"})

		assert.NoError(t, err)
		assert.Equal(t, "Hello, Alice!", output.Message)
	})

	t.Run("empty_name_returns_ErrNameRequired", func(t *testing.T) {
		mockRepo := hellomocks.NewGreetingRepository(t)

		uc := usecase.NewSayHelloUseCase(mockRepo, bypassLogger(t))
		output, err := uc.Execute(&usecase.SayHelloInput{Name: ""})

		assert.Nil(t, output)
		assert.ErrorIs(t, err, usecase.ErrNameRequired)
		mockRepo.AssertNotCalled(t, "Save", mock.Anything)
	})

	t.Run("repository_save_fails_returns_ErrGreetingFailed", func(t *testing.T) {
		mockRepo := hellomocks.NewGreetingRepository(t)

		mockRepo.On("Save", mock.MatchedBy(func(g *greeting.Greeting) bool {
			return g.GetName() == "Bob" && g.GetMessage() == "Hello, Bob!"
		})).Return(errors.New("db error")).Once()

		uc := usecase.NewSayHelloUseCase(mockRepo, bypassLogger(t))
		output, err := uc.Execute(&usecase.SayHelloInput{Name: "Bob"})

		assert.Nil(t, output)
		assert.ErrorIs(t, err, usecase.ErrGreetingFailed)
	})
}
