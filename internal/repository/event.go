package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"log/slog"
	"sync"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

func NewEventRepository(once *sync.Once, db *database.Postgres) *EventRepository {
	var repo *EventRepository
	once.Do(func() {
		repo = &EventRepository{db: db}
	})

	return repo
}

type EventRepository struct {
	db *database.Postgres
}

func (r *EventRepository) InsertEvent(ctx context.Context, userId int, event *entity.Event) error {
	const op = "Repo:InsertEvent"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	queryEvent, argsEvent, err := r.db.Builder.
		Insert("event").
		Columns("id, space_id, name, description, begin_date, end_date, tags").
		Values(event.ID, event.SpaceId, event.Name, event.Description, event.BeginDate, event.EndDate, event.Tags).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	queryUserEvent, argsUserEvent, err := r.db.Builder.
		Insert("user_event").
		Columns("user_id, event_id").
		Values(userId, event.ID).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Pool.Exec(ctx, queryEvent, argsEvent...)
	if err != nil {
		log.Debug("couldn't insert data in event: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, queryUserEvent, argsUserEvent...)
	if err != nil {
		log.Debug("couldn't insert data in event: ", err.Error())
		return fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Debug("couldn't commit transaction: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *EventRepository) GetEvent(ctx context.Context, id int) (*entity.Event, error) {
	const op = "Repo:GetEvent"

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", id),
	)
	log.Debug(op)

	fail := func(err error) (*entity.Event, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id, space_id, name, description, begin_date, end_date, tags").
		From("event").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	event := new(entity.Event)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&event.ID, &event.SpaceId, &event.Name, &event.Description, &event.BeginDate, &event.EndDate, &event.Tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("event not found: ", err.Error())
			return fail(ErrEventNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return event, nil
}

// TODO: Err already joined
func (r *EventRepository) AddUser(ctx context.Context, eventId, userId int) error {
	const op = "Repo:AddUserToEvent"

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", eventId),
		slog.Int("user id", userId),
	)

	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("user_event").
		Columns("user_id, event_id").
		Values(userId, eventId).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't insert data in user_event: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, id int) error {
	const op = "Repo:DeleteEvent"

	log := slog.With(
		slog.String("op", op),
		slog.Int("event id", id),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	queryEvent, argsEvent, err := r.db.Builder.
		Delete("event").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	queryUserEvent, argsUserEvent, err := r.db.Builder.
		Delete("user_event").
		Where("event_id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Pool.Exec(ctx, queryUserEvent, argsUserEvent...)
	if err != nil {
		log.Debug("couldn't delete data from user_event: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, queryEvent, argsEvent...)
	if err != nil {
		log.Debug("couldn't delete data from event: ", err.Error())
		return fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Debug("couldn't commit transaction: ", err.Error())
		return fail(err)
	}

	return nil
}
