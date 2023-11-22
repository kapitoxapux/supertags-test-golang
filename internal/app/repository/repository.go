package repository

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"supertags/internal/app/model"
	"supertags/internal/app/service"
)

type Repository interface {
	PutFact(m *service.Fact) (*service.Fact, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) PutFact(m *service.Fact) (*service.Fact, error) {
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func NewRepository(dns string) Repository {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Gorm repository failed %w", err.Error())
	}
	if exist := db.Migrator().HasTable(&model.Fact{}); !exist {
		db.Migrator().CreateTable(&model.Fact{})
	}

	return &repository{db}
}
