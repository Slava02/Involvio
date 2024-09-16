package repository

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/pkg/errors"
	"sync"
)

var (
	ErrInvalidInputData = errors.New("invalid input data")
)

type UserRepository struct {
	db database.Database
}

func NewUserRepository(once *sync.Once, db database.Database) *UserRepository {
	var repo *UserRepository
	once.Do(func() {
		repo = &UserRepository{db: db}
	})

	return repo
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	return &user, nil
}

func (r *UserRepository) InsertUser(ctx context.Context, input *entity.User) (*entity.User, error) {

	return input, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, input *entity.User) (*entity.User, error) {

	return input, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, input *entity.User) error {

	return nil
}
