package controller

import (
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/route/api/events"
	"github.com/Slava02/Involvio/internal/route/api/spaces"
	"github.com/Slava02/Involvio/internal/route/api/users"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/go-openapi/runtime/middleware"
	"log/slog"
)

type UseCase interface {
	CreateEvent(cmd *usecase.CreateEventCmd) (*entity.Event, error)
	GetEvent(cmd *usecase.GetEventCmd) (*entity.EventInfoResp, error)
	JoinEvent(cmd *usecase.JoinEventCmd) error
	CreateSpace(cmd *usecase.CreateSpaceCmd) (*entity.Space, error)
	GetSpace(cmd *usecase.GetSpaceCmd) (*entity.Space, error)
	JoinSpace(cmd *usecase.JoinSpaceCmd) error
	CreateUser(cmd *usecase.CreateUserCmd) (*entity.User, error)
	DeleteUser(cmd *usecase.DeleteUserCmd) error
	GetUser(cmd *usecase.GetUserCmd) (*entity.UserInfoResp, error)
	UpdateUser(cmd *usecase.UpdateUserCmd) error
}

type Impl struct {
	useCase UseCase
}

func NewImpl(u UseCase) *Impl {
	return &Impl{
		useCase: u,
	}
}

func (i Impl) CreateEvent(params events.CreateEventParams) middleware.Responder {
	p := params.Body

	cmd := &usecase.CreateEventCmd{
		Event:   p.Event,
		UserId:  p.UserID,
		SpaceId: p.SpaceID,
	}

	eventInfo, err := i.useCase.CreateEvent(cmd)
	if err != nil {
		slog.Error(err.Error())
		return events.NewJoinEventInternalServerError()
	}

	return events.NewCreateEventOK().WithPayload(eventInfo)
}

func (i Impl) GetEvent(params events.GetEventParams) middleware.Responder {
	cmd := &usecase.GetEventCmd{
		EventId: params.EventID,
	}

	eventInfo, err := i.useCase.GetEvent(cmd)
	if err != nil {
		slog.Error(err.Error())
		return events.NewJoinEventInternalServerError()
	}

	return events.NewGetEventOK().WithPayload(eventInfo)
}

func (i Impl) JoinEvent(p events.JoinEventParams) middleware.Responder {
	cmd := &usecase.JoinEventCmd{
		UserId:  p.Body.UserID,
		EventId: p.EventID,
	}

	err := i.useCase.JoinEvent(cmd)
	if err != nil {
		slog.Error(err.Error())
		return events.NewJoinEventInternalServerError()
	}

	return events.NewJoinEventOK()
}

func (i Impl) CreateSpace(params spaces.CreateSpaceParams) middleware.Responder {
	p := params.Body

	cmd := &usecase.CreateSpaceCmd{
		SpaceId:          p.SpaceID,
		SpaceName:        p.SpaceName,
		SpaceDescription: p.SpaceDescription,
		AdminId:          p.AdminID,
		Tags:             p.Tags,
	}

	spaceInfo, err := i.useCase.CreateSpace(cmd)
	if err != nil {
		slog.Error(err.Error())
		return spaces.NewCreateSpaceInternalServerError()
	}

	return spaces.NewCreateSpaceOK().WithPayload(spaceInfo)
}

func (i Impl) GetSpace(params spaces.GetSpaceParams) middleware.Responder {
	cmd := &usecase.GetSpaceCmd{
		SpaceId: params.SpaceID,
	}

	spaceInfo, err := i.useCase.GetSpace(cmd)
	if err != nil {
		slog.Error(err.Error())
		return spaces.NewGetSpaceInternalServerError()
	}

	return spaces.NewGetSpaceOK().WithPayload(spaceInfo)
}

func (i Impl) JoinSpace(p spaces.JoinSpaceParams) middleware.Responder {
	cmd := &usecase.JoinSpaceCmd{
		SpaceId: p.SpaceID,
		User:    p.Body,
	}

	err := i.useCase.JoinSpace(cmd)
	if err != nil {
		slog.Error(err.Error())
		return spaces.NewJoinSpaceInternalServerError()
	}

	return events.NewJoinEventOK()
}

func (i Impl) CreateUser(p users.CreateUserParams) middleware.Responder {
	cmd := &usecase.CreateUserCmd{
		User: p.Body,
	}

	userInfo, err := i.useCase.CreateUser(cmd)
	if err != nil {
		slog.Error(err.Error())
		return users.NewDeleteUserInternalServerError()
	}

	return users.NewCreateUserOK().WithPayload(userInfo)
}

func (i Impl) DeleteUser(p users.DeleteUserParams) middleware.Responder {
	cmd := &usecase.DeleteUserCmd{
		SpaceId: p.Body.SpaceID,
		UserId:  p.UserID,
	}

	err := i.useCase.DeleteUser(cmd)
	if err != nil {
		slog.Error(err.Error())
		return users.NewDeleteUserInternalServerError()
	}

	return users.NewDeleteUserOK()
}

func (i Impl) GetUser(params users.GetUserParams) middleware.Responder {
	cmd := &usecase.GetUserCmd{
		UserId: params.UserID,
	}

	userInfo, err := i.useCase.GetUser(cmd)
	if err != nil {
		slog.Error(err.Error())
		return users.NewGetUserInternalServerError()
	}

	return users.NewGetUserOK().WithPayload(userInfo)
}

func (i Impl) UpdateUser(params users.UpdateUserParams) middleware.Responder {
	cmd := &usecase.UpdateUserCmd{
		UserId: params.UserID,
		Form:   params.Body,
	}

	err := i.useCase.UpdateUser(cmd)
	if err != nil {
		slog.Error(err.Error())
		return users.NewUpdateUserInternalServerError()
	}

	return users.NewUpdateUserOK()
}
