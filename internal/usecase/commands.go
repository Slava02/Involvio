package usecase

import "github.com/Slava02/Involvio/internal/entity"

type CreateEventCmd struct {
	Event           *entity.Event
	UserId, SpaceId int64
}

type GetEventCmd struct {
	EventId int64
}

type JoinEventCmd struct {
	UserId, EventId int64
}

type CreateSpaceCmd struct {
	SpaceId          int64
	SpaceDescription string
	SpaceName        string
	AdminId          int64
	Tags             entity.TagOptions
}

type GetSpaceCmd struct {
	SpaceId int64
}

type JoinSpaceCmd struct {
	SpaceId int64
	User    *entity.User
}

type CreateUserCmd struct {
	User *entity.User
}

type DeleteUserCmd struct {
	UserId  int64
	SpaceId int64
}

type GetUserCmd struct {
	UserId int64
}

type UpdateUserCmd struct {
	UserId int64
	Form   *entity.Form
}
