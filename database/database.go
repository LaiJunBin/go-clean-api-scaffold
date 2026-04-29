package database

import (
	"go-clean-api-scaffold/internal/config"

	"gorm.io/gorm"
)

func DatabaseProvider(cfg *config.Config) *gorm.DB {
	return NewSqliteGormDB(cfg)
}
