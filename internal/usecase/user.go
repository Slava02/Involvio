package usecase

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/Slava02/Involvio/pkg/hexid"
	"log/slog"
)

type IUserRepository interface {
	GetUserData(ctx context.Context, id int) (*entity.User, error)
	GetUserForms(ctx context.Context, userId int) ([]*entity.Form, error)
	InsertUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, id int, firstName, lastName, userName, photoURL string) (*entity.User, error)
	DeleteUser(ctx context.Context, userId, spaceId int) error
	GetForm(ctx context.Context, userId, spaceId int) (*entity.Form, error)
	UpdateForm(ctx context.Context, userId, spaceId int, userTags, pairTags entity.Tags) error
}

func NewUserUseCase(ur IUserRepository) *UserUseCase {
	return &UserUseCase{userRepo: ur}
}

type UserUseCase struct {
	userRepo IUserRepository
}

func (uc *UserUseCase) GetUser(ctx context.Context, cmd commands.UserByIdCommand) (*entity.User, []*entity.Form, error) {
	const op = "Usecase:GetForm"

	fail := func(err error) (*entity.User, []*entity.Form, error) {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	userData, err := uc.userRepo.GetUserData(ctx, cmd.ID)
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	userForms, err := uc.userRepo.GetUserForms(ctx, cmd.ID)
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	return userData, userForms, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, cmd commands.FormByIdCommand) error {
	const op = "Usecase:DeleteUser"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, _, err := uc.GetUser(ctx, commands.UserByIdCommand{ID: cmd.UserID})
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	err = uc.userRepo.DeleteUser(ctx, cmd.UserID, cmd.SpaceID)
	if err != nil {
		log.Debug("couldn't delete user: ", err.Error())
		return fail(err)
	}

	return nil
}

func (uc *UserUseCase) GetForm(ctx context.Context, cmd commands.FormByIdCommand) (*entity.Form, error) {
	const op = "Usecase:GetForm"

	fail := func(err error) (*entity.Form, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	form, err := uc.userRepo.GetForm(ctx, cmd.UserID, cmd.SpaceID)
	if err != nil {
		return fail(err)
	}

	return form, nil
}

func (uc *UserUseCase) UpdateForm(ctx context.Context, cmd commands.UpdateFormCommand) (*entity.User, []*entity.Form, error) {
	const op = "Usecase:UpdateForm"

	fail := func(err error) (*entity.User, []*entity.Form, error) {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, err := uc.GetForm(ctx, commands.FormByIdCommand{UserID: cmd.UserID, SpaceID: cmd.SpaceID})
	if err != nil {
		log.Debug("couldn't get form: ", err.Error())
		return fail(err)
	}

	err = uc.userRepo.UpdateForm(ctx, cmd.UserID, cmd.SpaceID, cmd.UserTags, cmd.PairTags)
	if err != nil {
		return fail(err)
	}

	user, forms, err := uc.GetUser(ctx, commands.UserByIdCommand{ID: cmd.UserID})
	if err != nil {
		log.Debug("couldn't get form: ", err.Error())
		return fail(err)
	}

	return user, forms, nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, cmd commands.CreateUserCommand) (*entity.User, error) {
	const op = "Usecase:CreateUser"

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	userId, err := hexid.Generate()
	if err != nil {
		log.Error("couldn't generate id: ", err.Error())
		return fail(err)
	}

	user := &entity.User{
		ID:        userId,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		UserName:  cmd.UserName,
		PhotoURL:  cmd.PhotoURL,
		AuthDate:  cmd.AuthDate,
	}

	err = uc.userRepo.InsertUser(ctx, user)
	if err != nil {
		return fail(err)
	}

	return user, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entity.User, error) {
	const op = "Usecase:UpdateUser"

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, _, err := uc.GetUser(ctx, commands.UserByIdCommand{ID: cmd.ID})
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	user, err := uc.userRepo.UpdateUser(ctx, cmd.ID, cmd.FirstName, cmd.LastName, cmd.UserName, cmd.PhotoURL)
	if err != nil {
		log.Debug("couldn't update user: ", err.Error())
		return fail(err)
	}

	return user, nil
}
