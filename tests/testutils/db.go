package testutils

import (
	"go-clean-api-scaffold/database"
	"go-clean-api-scaffold/database/migration"
	"go-clean-api-scaffold/internal/config"

	"gorm.io/gorm"
)

var db = database.DatabaseProvider(config.NewTestConfig())

func GetDB() *gorm.DB {
	return db
}

func RunMigration(db *gorm.DB) {
	migration.RunMigrate(db)
}

func ResetDB(db *gorm.DB) {
	migration.RunRollback(db)
}
