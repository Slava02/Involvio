package repository

import (
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/pkg/errors"
	"sync"
)

var (
	ErrInvalidInputData = errors.New("invalid input data")
)

type Repository struct {
	db database.Database
}

func NewUserRepository(once *sync.Once, db database.Database) *Repository {
	var repo *Repository
	once.Do(func() {
		repo = &Repository{db: db}
	})

	return repo
}
