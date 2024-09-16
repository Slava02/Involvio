package usecase

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
)

type IUserRepository interface {
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	InsertUser(ctx context.Context, input *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, input *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, input *entity.User) error
}
