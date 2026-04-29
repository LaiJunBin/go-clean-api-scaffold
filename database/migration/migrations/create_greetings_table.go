package migrations

import (
	"go-clean-api-scaffold/internal/app/shared/model"

	"gorm.io/gorm"
)

type GreetingMigrator struct{}

func (m GreetingMigrator) Up(db *gorm.DB) error {
	return db.AutoMigrate(&model.Greeting{})
}

func (m GreetingMigrator) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.Greeting{})
}
