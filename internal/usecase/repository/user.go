package repository

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	ErrInvalidInputData = errors.New("invalid input data")
	ErrUserNotFound     = errors.New("user not found")
)

func NewUserRepository(once *sync.Once, db *database.Postgres) *UserRepository {
	var repo *UserRepository
	once.Do(func() {
		repo = &UserRepository{db: db}
	})

	return repo
}

type UserRepository struct {
	db *database.Postgres
}

func (r *UserRepository) GetUserData(ctx context.Context, id int) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetUserForms(ctx context.Context, userId int) ([]*entity.Form, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) InsertUser(ctx context.Context, firstName, lastName, userName, photoURL string, authDate time.Time) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, firstName, lastName, userName, photoURL string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId, spaceId int) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetForm(ctx context.Context, userId, spaceId int) (*entity.Form, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) UpdateForm(ctx context.Context, form *entity.Form) (*entity.Form, error) {
	//TODO implement me
	panic("implement me")
}
