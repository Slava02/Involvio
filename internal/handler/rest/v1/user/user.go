package user

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type IUserUseCase interface {
	GetUser(ctx context.Context, cmd usecase.UserByIdCommand) (*entity.User, []*entity.Form, error)
	CreateUser(ctx context.Context, cmd usecase.CreateUserCommand) (*entity.User, error)
	UpdateUser(ctx context.Context, cmd usecase.UpdateUserCommand) (*entity.User, error)
	DeleteUser(ctx context.Context, cmd usecase.FormByIdCommand) error
	GetForm(ctx context.Context, cmd usecase.FormByIdCommand) (*entity.Form, error)
	UpdateForm(ctx context.Context, cmd usecase.UpdateFormCommand) (*entity.Form, error)
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
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "FindUserByID", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.UserByIdCommand{ID: req.ID}

	user, forms, err := uh.userUC.GetUser(ctx, cmd)
	if user == nil && err == nil {
		return nil, huma.Error404NotFound("user not found")
	} else if err != nil {
		return nil, huma.Error500InternalServerError("getting user by id error. ", err)
	}

	resp := ToUserWithFormsOutputFromEntity(user, forms)

	return resp, nil
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "CreateUser", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.CreateUserCommand{
		FirstName: req.Body.FirstName,
		LastName:  req.Body.LastName,
		UserName:  req.Body.UserName,
		PhotoURL:  req.Body.UserName,
		AuthDate:  time.Now(),
	}

	user, err := uh.userUC.CreateUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := ToUserOutputFromEntity(user)

	return resp, nil
}

func (uh *UserHandler) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "UpdateUser", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.UpdateUserCommand{
		ID:        req.ID,
		FirstName: req.Body.FirstName,
		LastName:  req.Body.LastName,
		UserName:  req.Body.UserName,
		PhotoURL:  req.Body.UserName,
	}

	user, err := uh.userUC.UpdateUser(ctx, cmd)

	if err != nil {
		return nil, huma.Error500InternalServerError("update user error. ", err)
	}

	resp := ToUserOutputFromEntity(user)

	return resp, nil
}

func (uh *UserHandler) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*struct{}, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "DeleteUser", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.FormByIdCommand{UserID: req.UserId, SpaceID: req.SpaceId}

	err := uh.userUC.DeleteUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error400BadRequest("user not found", err)
	}

	return &struct{}{}, nil
}

func (uh *UserHandler) GetForm(ctx context.Context, req *FormByIdRequest) (*FormResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "FindUserFormByID", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.FormByIdCommand{UserID: req.UserID, SpaceID: req.SpaceID}

	form, err := uh.userUC.GetForm(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("getting user by id error. ", err)
	}

	resp := ToFormOutputFromEntity(form)

	return resp, nil
}

func (uh *UserHandler) UpdateForm(ctx context.Context, req *UpdateFormRequest) (*FormResponse, error) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "FindUserFormByID", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	cmd := usecase.UpdateFormCommand{Form: &entity.Form{
		UserID:   req.UserID,
		SpaceID:  req.SpaceID,
		Admin:    req.Body.Admin,
		Creator:  req.Body.Creator,
		UserTags: req.Body.UserTags,
		PairTags: req.Body.PairTags,
	}}

	// TODO: в таких случаях надо будет по-человечески ошибки обрабатывать свитчом
	form, err := uh.userUC.UpdateForm(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("getting user by id error. ", err)
	}

	resp := ToFormOutputFromEntity(form)

	return resp, nil
}
