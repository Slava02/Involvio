package usecase

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/Slava02/Involvio/pkg/hexid"
)

type IEventRepository interface {
	InsertEvent(ctx context.Context, userId int, event *entity.Event) error
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
	// TODO: вынести генерацию id в зависимость
	eventId, err := hexid.Generate()
	if err != nil {
		// TODO: error couldn't generate ID
		return nil, fmt.Errorf("%w", err)
	}

	event := &entity.Event{
		ID:          eventId,
		SpaceId:     cmd.SpaceId,
		Name:        cmd.Name,
		Description: cmd.Description,
		Tags:        cmd.Tags,
		BeginDate:   cmd.BeginDate,
		EndDate:     cmd.EndDate,
	}

	err = ec.eventRepo.InsertEvent(ctx, cmd.UserId, event)
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
