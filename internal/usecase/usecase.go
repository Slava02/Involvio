package usecase

import (
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
)

type Repository interface {
	CreateEvent(userId, spaceId int64, event *entity.Event) (*entity.Event, error)
	GetEventById(eventId int64) (*entity.Event, error)
	GetEventUsers(eventId int64) ([]*entity.User, error)
	AddUserToEvent(eventId, userId int64) error
	CreateSpace(spaceName, spaceDescription string, adminId int64, tags entity.TagOptions) (*entity.Space, error)
	GetSpaceById(spaceId int64) (*entity.Space, error)
	AddUserToSpace(spaceId, userId int64) error
	CreateUser(user *entity.User) (*entity.User, error) // TODO: не нужно возвращать пользователя
	DeleteUser(userId, spaceId int64) error
	GetUserById(userId int64) (*entity.User, error)
	GetFormsByUserId(userId int64) ([]*entity.Form, error)
	UpdateUser(userId int64, form *entity.Form) error
}

func NewUseCase(r Repository) *UseCase {
	return &UseCase{repo: r}
}

type UseCase struct {
	repo Repository
}

func (u UseCase) CreateEvent(cmd *CreateEventCmd) (*entity.Event, error) {
	event, err := u.repo.CreateEvent(cmd.UserId, cmd.SpaceId, cmd.Event)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return event, nil
}

func (u UseCase) GetEvent(cmd *GetEventCmd) (*entity.EventInfoResp, error) {
	event, err := u.repo.GetEventById(cmd.EventId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	users, err := u.repo.GetEventUsers(cmd.EventId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &entity.EventInfoResp{
		Event: event,
		Users: users,
	}, nil
}

// TODO: сделать ID у event
func (u UseCase) JoinEvent(cmd *JoinEventCmd) error {
	_, err := u.repo.GetEventById(cmd.EventId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = u.repo.AddUserToEvent(cmd.EventId, cmd.UserId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// TODO: сделать DTO для создания спейса без ID
func (u UseCase) CreateSpace(cmd *CreateSpaceCmd) (*entity.Space, error) {
	space, err := u.repo.CreateSpace(cmd.SpaceName, cmd.SpaceDescription, cmd.AdminId, cmd.Tags)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return space, nil
}

func (u UseCase) GetSpace(cmd *GetSpaceCmd) (*entity.Space, error) {
	space, err := u.repo.GetSpaceById(cmd.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return space, nil
}

// TODO: оставить только UserId в JoinSpaceCmd
func (u UseCase) JoinSpace(cmd *JoinSpaceCmd) error {
	_, err := u.repo.GetSpaceById(cmd.SpaceId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = u.repo.AddUserToSpace(cmd.SpaceId, cmd.User.UserID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (u UseCase) CreateUser(cmd *CreateUserCmd) (*entity.User, error) {
	user, err := u.repo.CreateUser(cmd.User)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return user, nil
}

func (u UseCase) DeleteUser(cmd *DeleteUserCmd) error {
	err := u.repo.DeleteUser(cmd.UserId, cmd.SpaceId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (u UseCase) GetUser(cmd *GetUserCmd) (*entity.UserInfoResp, error) {
	user, err := u.repo.GetUserById(cmd.UserId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	forms, err := u.repo.GetFormsByUserId(cmd.UserId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp := &entity.UserInfoResp{
		User:  user,
		Forms: forms,
	}

	return resp, nil
}

// TODO: из формы, наверное, стоит убрать space_id
func (u UseCase) UpdateUser(cmd *UpdateUserCmd) error {
	err := u.repo.UpdateUser(cmd.UserId, cmd.Form)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
