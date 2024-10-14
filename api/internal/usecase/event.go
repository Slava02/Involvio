package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/repository"
	"github.com/Slava02/Involvio/api/internal/usecase/commands"
	"github.com/Slava02/Involvio/api/pkg/hexid"
	"log/slog"
)

type IEventRepository interface {
	InsertEvent(ctx context.Context, event *entity.Event) error
	GetEvent(ctx context.Context, id int) (*entity.Event, error)
	GetUserEvents(ctx context.Context, id int) ([]entity.Event, error)
	AddUser(ctx context.Context, eventId, userId int) error
	DeleteEvent(ctx context.Context, id int) error
	AddReview(ctx context.Context, eventId, who, whom, grade int) error
}

func NewEventUseCase(ur IEventRepository) *EventUseCase {
	return &EventUseCase{eventRepo: ur}
}

type EventUseCase struct {
	eventRepo IEventRepository
}

func (ec *EventUseCase) ReviewEvent(ctx context.Context, cmd commands.ReviewEventCommand) error {
	const op = "UseCase:ReviewEvent"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", cmd.EventID),
	)
	log.Debug(op)

	_, err := ec.eventRepo.GetEvent(ctx, cmd.EventID)
	if err != nil {
		log.Debug("couldn't get event: ", err.Error())
		return fail(err)
	}

	err = ec.eventRepo.AddReview(ctx, cmd.EventID, cmd.WhoID, cmd.WhomID, cmd.Grade)
	if err != nil {
		log.Debug("couldn't add user review to event: ", err.Error())
		return fail(err)
	}

	return nil
}

func (ec *EventUseCase) CreateEvent(ctx context.Context, cmd commands.CreateEventCommand) (*entity.Event, error) {
	const op = "UseCase:CreateEvent"

	fail := func(err error) (*entity.Event, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.String("name", cmd.Name),
	)
	log.Debug(op)

	// TODO: вынести генератор id в зависимость
	eventId, err := hexid.Generate()
	if err != nil {
		log.Error("couldn't generate id: ", err.Error())
		return fail(err)
	}

	event := &entity.Event{
		ID:          eventId,
		Name:        cmd.Name,
		Description: cmd.Description,
		Date:        cmd.Date,
		Users:       cmd.Users,
	}

	err = ec.eventRepo.InsertEvent(ctx, event)
	if err != nil {
		log.Debug("couldn't insert event: ", err.Error())
		return fail(err)
	}

	return event, nil
}

func (ec *EventUseCase) GetUserEvents(ctx context.Context, cmd commands.EventByUserIdCommand) ([]entity.Event, error) {
	const op = "UseCase:GetEvents"

	fail := func(err error) ([]entity.Event, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", cmd.ID),
	)
	log.Debug(op)

	events, err := ec.eventRepo.GetUserEvents(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Info("couldn't get events: ", err.Error())
		} else {
			log.Debug("couldn't get events: ", err.Error())
		}
		return fail(err)
	}

	return events, nil
}

func (ec *EventUseCase) JoinEvent(ctx context.Context, cmd commands.JoinEventCommand) error {
	const op = "UseCase:JoinEvent"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", cmd.EventId),
		slog.Int("user id", cmd.UserId),
	)
	log.Debug(op)

	_, err := ec.eventRepo.GetEvent(ctx, cmd.EventId)
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
	const op = "UseCase:DeleteEvent"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("eventID", cmd.ID),
	)
	log.Debug(op)

	_, err := ec.eventRepo.GetEvent(ctx, cmd.ID)
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
