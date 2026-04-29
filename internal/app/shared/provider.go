package shared

import (
	"go-clean-api-scaffold/internal/app/shared/infrastructure/logger"
	"go-clean-api-scaffold/internal/app/shared/types/provider"
)

func NewProvider() *provider.Provider {
	return provider.New(provider.Config{
		Infrastructures: []interface{}{
			logger.NewZapLogger,
		},
	})
}
