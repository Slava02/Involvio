package user

import (
	"context"
	"errors"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/repository"
	"github.com/Slava02/Involvio/api/internal/usecase"
	"github.com/Slava02/Involvio/api/internal/usecase/commands"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type IUserUseCase interface {
	GetUser(ctx context.Context, cmd commands.UserByUsernameCommand) (*entity.User, error)
	CreateUser(ctx context.Context, cmd commands.CreateUserCommand) error
	UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entity.User, error)
	BlockUser(ctx context.Context, cmd commands.BlockUserCommand) error
	SetHoliday(ctx context.Context, cmd commands.SetHolidayCommand) (*entity.User, error)
	CancelHoliday(ctx context.Context, cmd commands.CancelHolidayCommand) error
}

var _ IUserUseCase = (*usecase.UserUseCase)(nil)

const tracerName = "user handler"

type UserHandler struct {
	userUC IUserUseCase
}

func NewUserHandler(uc IUserUseCase) *UserHandler {
	return &UserHandler{userUC: uc}
}

func (uh *UserHandler) GetUser(ctx context.Context, req *UserByUsernameRequest) (*UserRequestResponse, error) {
	const op = "Handler:GetUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("username", req.Username),
	)
	log.Debug(op)

	cmd := commands.UserByUsernameCommand{Username: req.Username}

	user, err := uh.userUC.GetUser(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			log.Info("couldn't get user: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get user: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToUserOutputFromEntity(user)

	return resp, nil
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *UserRequestResponse) (*struct{}, error) {
	const op = "Handler:CreateUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.String("username", req.Body.UserName),
		slog.Int("id", req.Body.ID),
	)
	log.Debug(op)

	cmd := commands.CreateUserCommand{
		User: req.Body.User,
	}

	err := uh.userUC.CreateUser(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAlreadyExists):
			log.Info("couldn't create user: ", err.Error())
			return &struct{}{}, huma.Error400BadRequest("user in group already exists")
		case errors.Is(err, repository.ErrNotFound):
			log.Info("couldn't create user: ", err.Error())
			return &struct{}{}, huma.Error404NotFound("group not found")
		default:
			log.Error("couldn't create user: ", err.Error())
			return &struct{}{}, huma.Error500InternalServerError("internal service error")
		}
	}

	return &struct{}{}, nil
}

func (uh *UserHandler) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserRequestResponse, error) {
	const op = "Handler:UpdateUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.UpdateUserCommand{
		ID:        req.Body.ID,
		FullName:  req.Body.FullName,
		PhotoURL:  req.Body.PhotoURL,
		City:      req.Body.City,
		Position:  req.Body.Position,
		Interests: req.Body.Interests,
	}

	user, err := uh.userUC.UpdateUser(ctx, cmd)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			log.Info("couldn't update user: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't update user: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToUserOutputFromEntity(user)

	return resp, nil
}

func (uh *UserHandler) BlockUser(ctx context.Context, req *BlockUserRequest) (*struct{}, error) {
	const op = "Handler:BlockUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.BlockUserCommand{
		WhoID:  req.Body.WhoID,
		WhomID: req.Body.WhomID,
	}

	err := uh.userUC.BlockUser(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't block user: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	return &struct{}{}, nil
}

func (uh *UserHandler) SetHoliday(ctx context.Context, req *SetHolidayRequest) (*UserRequestResponse, error) {
	const op = "Handler:SetHoliday"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.SetHolidayCommand{
		ID:       req.Body.ID,
		TillDate: req.Body.TillDate,
	}

	user, err := uh.userUC.SetHoliday(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't set holiday: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	resp := ToUserOutputFromEntity(user)

	return resp, nil
}

func (uh *UserHandler) CancelHoliday(ctx context.Context, req *UserByIDRequest) (*struct{}, error) {
	const op = "Handler:CancelHoliday"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.CancelHolidayCommand{
		ID: req.Body.ID,
	}

	err := uh.userUC.CancelHoliday(ctx, cmd)
	if err != nil {
		switch {
		default:
			log.Error("couldn't cancel holiday: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	return &struct{}{}, nil
}
