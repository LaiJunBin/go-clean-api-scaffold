package database

import (
	"go-clean-api-scaffold/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteGormDB(cfg *config.Config) *gorm.DB {
	path := cfg.Database.Sqlite.Path
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect to sqlite database: " + err.Error())
	}

	return db
}
