package usecase

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/usecase/commands"
	"log/slog"
	"time"
)

type IUserRepository interface {
	BlockUser(ctx context.Context, who, whom int) error
	SetHoliday(ctx context.Context, id int, tillDate time.Time) (*entity.User, error)
	CancelHoliday(ctx context.Context, id int) error
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, id int, fullName, city, position, interests, photoURL string) (*entity.User, error)
}

func NewUserUseCase(ur IUserRepository) *UserUseCase {
	return &UserUseCase{userRepo: ur}
}

type UserUseCase struct {
	userRepo IUserRepository
}

func (uc *UserUseCase) BlockUser(ctx context.Context, cmd commands.BlockUserCommand) error {
	const op = "UseCase:BlockUser"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("who", cmd.WhoID),
		slog.Int("whom", cmd.WhomID),
	)
	log.Debug(op)

	err := uc.userRepo.BlockUser(ctx, cmd.WhoID, cmd.WhomID)
	if err != nil {
		log.Debug("couldn't block user: ", err.Error())
		return fail(err)
	}

	return nil
}

func (uc *UserUseCase) SetHoliday(ctx context.Context, cmd commands.SetHolidayCommand) (*entity.User, error) {
	const op = "UseCase:SetHoliday"

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", cmd.ID),
		slog.String("tillDate", cmd.TillDate.String()),
	)
	log.Debug(op)

	user, err := uc.userRepo.SetHoliday(ctx, cmd.ID, cmd.TillDate)
	if err != nil {
		log.Debug("couldn't set holiday: ", err.Error())
		return fail(err)
	}

	return user, nil
}

func (uc *UserUseCase) CancelHoliday(ctx context.Context, cmd commands.CancelHolidayCommand) error {
	const op = "UseCase:CancelHoliday"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", cmd.ID),
	)
	log.Debug(op)

	err := uc.userRepo.CancelHoliday(ctx, cmd.ID)
	if err != nil {
		log.Debug("couldn't cancel holiday: ", err.Error())
		return fail(err)
	}

	return nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, cmd commands.UserByUsernameCommand) (*entity.User, error) {
	const op = "UseCase:GetUser"

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	user, err := uc.userRepo.GetUserByUsername(ctx, cmd.Username)
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	return user, nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, cmd commands.CreateUserCommand) error {
	const op = "UseCase:CreateUser"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.String("username", cmd.User.UserName),
		slog.Int("userID", cmd.User.ID),
	)
	log.Debug(op)

	err := uc.userRepo.InsertUser(ctx, cmd.User)
	if err != nil {
		return fail(err)
	}

	return nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entity.User, error) {
	const op = "UseCase:UpdateUser"

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", cmd.ID),
	)
	log.Debug(op)

	_, err := uc.userRepo.GetUserByID(ctx, cmd.ID)
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	user, err := uc.userRepo.UpdateUser(ctx, cmd.ID, cmd.FullName, cmd.City, cmd.Position, cmd.Interests, cmd.PhotoURL)
	if err != nil {
		log.Debug("couldn't update user: ", err.Error())
		return fail(err)
	}

	return user, nil
}
