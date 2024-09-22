package repository

import (
	"github.com/Slava02/Involvio/pkg/database"
	"sync"
)

func NewSpaceRepository(once *sync.Once, db *database.Postgres) *SpaceRepository {
	var repo *SpaceRepository
	once.Do(func() {
		repo = &SpaceRepository{db: db}
	})

	return repo
}

type SpaceRepository struct {
	db *database.Postgres
}
