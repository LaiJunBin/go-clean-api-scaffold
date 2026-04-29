package migration

import (
	"fmt"
	"go-clean-api-scaffold/database/migration/migrations"

	"gorm.io/gorm"
)

type Migrator interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

var Migrations = []Migrator{
	migrations.GreetingMigrator{},
}

func RunMigrate(db *gorm.DB) {
	fmt.Println("Running migrations...")
	for _, m := range Migrations {
		fmt.Printf("Migrating: %T\n", m)
		if err := m.Up(db); err != nil {
			panic(err)
		}
	}
	fmt.Println("Migrations completed.")
}

func RunRollback(db *gorm.DB) {
	fmt.Println("Rolling back migrations...")
	for i := len(Migrations) - 1; i >= 0; i-- {
		fmt.Printf("Rolling back: %T\n", Migrations[i])
		if err := Migrations[i].Down(db); err != nil {
			panic(err)
		}
	}
	fmt.Println("Rollback completed.")
}
