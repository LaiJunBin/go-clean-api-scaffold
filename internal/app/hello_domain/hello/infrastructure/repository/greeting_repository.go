package repository

import (
	"errors"

	domaingreeting "go-clean-api-scaffold/internal/app/hello_domain/hello/domain/entity/greeting"
	domainrepo "go-clean-api-scaffold/internal/app/hello_domain/hello/domain/repository"
	"go-clean-api-scaffold/internal/app/shared/model"

	"gorm.io/gorm"
)

type greetingRepository struct {
	db *gorm.DB
}

func NewGreetingRepository(db *gorm.DB) domainrepo.GreetingRepository {
	return &greetingRepository{db: db}
}

func (r *greetingRepository) Save(g *domaingreeting.Greeting) error {
	record := &model.Greeting{
		Name:    g.GetName(),
		Message: g.GetMessage(),
	}
	if err := r.db.Create(record).Error; err != nil {
		return err
	}
	g.SetID(int(record.ID))
	return nil
}

func (r *greetingRepository) FindByName(name string) (*domaingreeting.Greeting, error) {
	var record model.Greeting
	if err := r.db.Where("name = ?", name).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	g := domaingreeting.Reconstruct(int(record.ID), record.Name, record.Message)
	return g, nil
}
