package v1

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/danielgtaylor/huma/v2"
)

var _ IUserUseCase = (*usecase.UserUseCase)(nil)

const tracerName = "user handler"

type UserHandler struct {
	userUC IUserUseCase
}

func NewUserHandler(uc IUserUseCase) *UserHandler {
	return &UserHandler{userUC: uc}
}

func (uh *UserHandler) FindUserByID(ctx context.Context, req *FindUserRequest) (*UserResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "FindUserByID", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.FindUserByIDCommand{ID: req.ID}

	user, err := uh.userUC.FindUserByID(ctx, cmd)
	if user == nil && err == nil {
		return nil, huma.Error404NotFound("user not found")
	}

	if err != nil {
		return nil, huma.Error500InternalServerError("getting user by id error. ", err)
	}

	resp := &UserResponse{
		Body: struct{ *entity.User }{User: user},
	}

	return resp, nil
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "CreateUser", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.CreateUpdateUserCommand{User: req.Body.User}

	user, err := uh.userUC.CreateUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := &UserResponse{
		Body: struct{ *entity.User }{User: user},
	}

	return resp, nil
}

func (uh *UserHandler) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "UpdateUser", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.CreateUpdateUserCommand{User: req.Body.User}

	user, err := uh.userUC.UpdateUser(ctx, cmd)

	if err != nil {
		return nil, huma.Error500InternalServerError("update user error. ", err)
	}

	resp := &UserResponse{
		Body: struct{ *entity.User }{User: user},
	}

	return resp, nil
}

func (uh *UserHandler) DeleteUser(ctx context.Context, req *FindUserRequest) (*struct{}, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "DeleteUser", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.DeleteUserByIDCommand{ID: req.ID}

	err := uh.userUC.DeleteUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error400BadRequest("user not found", err)
	}

	return &struct{}{}, nil
}
