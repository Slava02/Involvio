package usecase

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"time"
)

type IEventRepository interface {
	InsertEvent(ctx context.Context, userId, spaceId int, name, description string, tags entity.Tags, beginDate, endDate time.Time) (*entity.Event, error)
	GetEvent(ctx context.Context, id int) (*entity.Event, error)
	AddUser(ctx context.Context, eventId, userId int) error
	DeleteEvent(ctx context.Context, id int) error
}

func NewEventUseCase(ur IEventRepository) *EventUseCase {
	return &EventUseCase{eventRepo: ur}
}

type EventUseCase struct {
	eventRepo IEventRepository
}

func (ec *EventUseCase) CreateEvent(ctx context.Context, cmd commands.CreateEventCommand) (*entity.Event, error) {
	event, err := ec.eventRepo.InsertEvent(ctx, cmd.UserId, cmd.SpaceId, cmd.Name, cmd.Description, cmd.Tags, cmd.BeginDate, cmd.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return event, nil
}

func (ec *EventUseCase) GetEvent(ctx context.Context, cmd commands.EventByIdCommand) (*entity.Event, error) {
	event, err := ec.eventRepo.GetEvent(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return event, nil
}

func (ec *EventUseCase) JoinEvent(ctx context.Context, cmd commands.JoinEventCommand) error {
	err := ec.eventRepo.AddUser(ctx, cmd.EventId, cmd.UserId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (ec *EventUseCase) DeleteEvent(ctx context.Context, cmd commands.EventByIdCommand) error {
	err := ec.eventRepo.DeleteEvent(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
