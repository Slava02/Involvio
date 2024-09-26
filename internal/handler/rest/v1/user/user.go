package user

import (
	"context"
	"errors"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type IUserUseCase interface {
	GetUser(ctx context.Context, cmd commands.UserByIdCommand) (*entity.User, []*entity.Form, error)
	CreateUser(ctx context.Context, cmd commands.CreateUserCommand) (*entity.User, error)
	UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entity.User, error)
	DeleteUser(ctx context.Context, cmd commands.FormByIdCommand) error
	GetForm(ctx context.Context, cmd commands.FormByIdCommand) (*entity.Form, error)
	UpdateForm(ctx context.Context, cmd commands.UpdateFormCommand) (*entity.User, []*entity.Form, error)
}

var _ IUserUseCase = (*usecase.UserUseCase)(nil)

const tracerName = "user handler"

type UserHandler struct {
	userUC IUserUseCase
}

func NewUserHandler(uc IUserUseCase) *UserHandler {
	return &UserHandler{userUC: uc}
}

func (uh *UserHandler) GetUserWithForms(ctx context.Context, req *UserByIdRequest) (*UserWithFormsResponse, error) {
	const op = "Handler:GetUserWithForms"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", req.ID),
	)
	log.Debug(op)

	cmd := commands.UserByIdCommand{ID: req.ID}

	user, forms, err := uh.userUC.GetUser(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			log.Info("couldn't get user: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get user: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToUserWithFormsOutputFromEntity(user, forms)

	return resp, nil
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	const op = "Handler:CreateUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.CreateUserCommand{
		FirstName: req.Body.FirstName,
		LastName:  req.Body.LastName,
		UserName:  req.Body.Username,
		PhotoURL:  req.Body.PhotoURL,
		AuthDate:  time.Now(),
	}

	user, err := uh.userUC.CreateUser(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			log.Info("couldn't create user: ", err.Error())
			return nil, huma.Error400BadRequest("user in space already exists")
		case errors.Is(err, repository.ErrUserNotFound):
			log.Info("couldn't create user: ", err.Error())
			return nil, huma.Error404NotFound("space not found")
		default:
			log.Error("couldn't create user: ", err.Error())
			return nil, huma.Error500InternalServerError("internal service error")
		}
	}

	resp := ToUserOutputFromEntity(user)

	return resp, nil
}

func (uh *UserHandler) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	const op = "Handler:UpdateUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.UpdateUserCommand{
		ID:        req.ID,
		FirstName: req.Body.FirstName,
		LastName:  req.Body.LastName,
		UserName:  req.Body.Username,
		PhotoURL:  req.Body.Username,
	}

	user, err := uh.userUC.UpdateUser(ctx, cmd)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
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

func (uh *UserHandler) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*struct{}, error) {
	const op = "Handler:DeleteUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.FormByIdCommand{UserID: req.UserId, SpaceID: req.SpaceId}

	err := uh.userUC.DeleteUser(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			log.Info("couldn't delete user: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't delete user: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	return &struct{}{}, nil
}

func (uh *UserHandler) GetForm(ctx context.Context, req *FormByIdRequest) (*FormResponse, error) {
	const op = "Handler:DeleteUser"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.FormByIdCommand{UserID: req.UserID, SpaceID: req.SpaceID}

	form, err := uh.userUC.GetForm(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			log.Info("couldn't get form: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't get form: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToFormOutputFromEntity(form)

	return resp, nil
}

func (uh *UserHandler) UpdateForm(ctx context.Context, req *UpdateFormRequest) (*UserWithFormsResponse, error) {
	const op = "Handler:UpdateForm"

	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, op, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	cmd := commands.UpdateFormCommand{
		UserID:   req.UserID,
		SpaceID:  req.SpaceID,
		UserTags: req.Body.UserTags,
		PairTags: req.Body.PairTags,
	}

	user, forms, err := uh.userUC.UpdateForm(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			log.Info("couldn't update user: ", err.Error())
			return nil, huma.Error404NotFound(err.Error())
		default:
			log.Error("couldn't update user: ", err.Error())
			return nil, huma.Error500InternalServerError(err.Error())
		}
	}

	resp := ToUserWithFormsOutputFromEntity(user, forms)

	return resp, nil
}
