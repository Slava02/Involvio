package v1

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks.go -package=v1

type IUserUseCase interface {
	FindUserByID(ctx context.Context, cmd usecase.FindUserByIDCommand) (*entity.User, error)
	CreateUser(ctx context.Context, cmd usecase.CreateUpdateUserCommand) (*entity.User, error)
	UpdateUser(ctx context.Context, cmd usecase.CreateUpdateUserCommand) (*entity.User, error)
	DeleteUser(ctx context.Context, cmd usecase.DeleteUserByIDCommand) error
}
