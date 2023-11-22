package mysql

import (
	"log"
	// "time"

	"supertags/internal/app/config"
	"supertags/internal/app/repository"
	"supertags/internal/app/service"
)

type DB struct {
	repo repository.Repository
}

func NewDB() *DB {
	repo := repository.NewRepository(config.GetDB())

	return &DB{
		repo: repo,
	}
}

func (db *DB) InsertFact(m *service.Fact) (*service.Fact, error) {
	m, err := db.repo.PutFact(m)
	if err != nil {
		log.Fatal("Произошла ошибка добавления: %w", err)
	}

	return m, nil
}
