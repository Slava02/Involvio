package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/Slava02/Involvio/pkg/hexid"
	"log/slog"
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
	const op = "Usecase:CreateEvent"

	fail := func(err error) (*entity.Event, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)

	log.Debug(op)

	// TODO: вынести генерацию id в зависимость
	eventId, err := hexid.Generate()
	if err != nil {
		log.Error("couldn't generate id: ", err.Error())
		return fail(err)
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
		log.Debug("couldn't insert event: ", err.Error())
		return fail(err)
	}

	return event, nil
}

func (ec *EventUseCase) GetEvent(ctx context.Context, cmd commands.EventByIdCommand) (*entity.Event, error) {
	const op = "Usecase:GetEvent"

	fail := func(err error) (*entity.Event, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", cmd.ID),
	)

	log.Debug(op)

	event, err := ec.eventRepo.GetEvent(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			log.Info("couldn't get event: ", err.Error())
		} else {
			log.Debug("couldn't get event: ", err.Error())
		}
		return fail(err)
	}

	return event, nil
}

func (ec *EventUseCase) JoinEvent(ctx context.Context, cmd commands.JoinEventCommand) error {
	const op = "Usecase:JoinEvent"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", cmd.EventId),
		slog.Int("user id", cmd.UserId),
	)

	log.Debug(op)

	_, err := ec.GetEvent(ctx, commands.EventByIdCommand{ID: cmd.EventId})
	if err != nil {
		log.Debug("couldn't get event: ", err.Error())
		return fail(err)
	}

	err = ec.eventRepo.AddUser(ctx, cmd.EventId, cmd.UserId)
	if err != nil {
		log.Debug("couldn't add user to event: ", err.Error())
		return fail(err)
	}

	return nil
}

func (ec *EventUseCase) DeleteEvent(ctx context.Context, cmd commands.EventByIdCommand) error {
	const op = "Usecase:DeleteEvent"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", cmd.ID),
	)

	log.Debug(op)

	_, err := ec.GetEvent(ctx, commands.EventByIdCommand{ID: cmd.ID})
	if err != nil {
		log.Debug("couldn't get event: ", err.Error())
		return fail(err)
	}

	err = ec.eventRepo.DeleteEvent(ctx, cmd.ID)
	if err != nil {
		log.Debug("couldn't delete event: ", err.Error())
		return fail(err)
	}

	return nil
}
