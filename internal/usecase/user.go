package usecase

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/pkg/errors"
)

type UserUseCase struct {
	userRepo IUserRepository
}

var ErrUserNotFound = errors.New("user not found")

func NewUserUseCase(ur IUserRepository) *UserUseCase {
	return &UserUseCase{userRepo: ur}
}

func (uc *UserUseCase) FindUserByID(ctx context.Context, cmd FindUserByIDCommand) (*entity.User, error) {
	user, err := uc.userRepo.GetUserByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, cmd CreateUpdateUserCommand) (*entity.User, error) {
	user, err := uc.userRepo.InsertUser(ctx, &cmd.User)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, cmd CreateUpdateUserCommand) (*entity.User, error) {
	userByID, err := uc.FindUserByID(ctx, FindUserByIDCommand{ID: cmd.User.ID})
	if err != nil {
		return nil, err
	}
	if userByID == nil {
		return nil, ErrUserNotFound
	}
	user, err := uc.userRepo.UpdateUser(ctx, &cmd.User)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, cmd DeleteUserByIDCommand) error {
	userByID, err := uc.FindUserByID(ctx, FindUserByIDCommand{ID: cmd.ID})
	if err != nil {
		return err
	}
	if userByID == nil {
		return ErrUserNotFound
	}
	return uc.userRepo.DeleteUser(ctx, userByID)
}
