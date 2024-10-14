package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/pkg/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"sync"
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

func (r *EventRepository) GetUserEvents(ctx context.Context, id int) ([]entity.Event, error) {
	const op = "Repo:GetUserEvents"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) ([]entity.Event, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mainQuery := r.db.Builder.
		Select("id, date, name, description").
		From("event")

	subQuery := r.db.Builder.
		Select("event_id").
		From("event_members").
		Where(sq.Eq{"user_id": id})

	mainQuery = mainQuery.Where(subQuery.Prefix("id IN (").Suffix(")"))

	query, args, err := mainQuery.ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	events := make([]entity.Event, 0)

	rows, err := r.db.Pool.Query(ctx, query, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found: ", err.Error())
			return fail(ErrNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	for rows.Next() {
		event := new(entity.Event)

		err = rows.Scan(&event.ID, &event.Date, &event.Name, &event.Description)
		if err != nil {
			return fail(err)
		}

		//  TODO: fix
		events = append(events, *event)
	}

	return events, nil
}

// TODO: make reviews auto-increment and maybe return the review
func (r *EventRepository) AddReview(ctx context.Context, eventID, who, whom, grade int) error {
	const op = "Repo:AddReview"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("reviews").
		Columns("event_id, who_id, about_whom_id, grade").
		Values(eventID, who, whom, grade).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't insert data in event: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *EventRepository) InsertEvent(ctx context.Context, event *entity.Event) error {
	const op = "Repo:InsertEvent"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("event").
		Columns("id, date, name, description").
		Values(event.ID, event.Date, event.Name, event.Description).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't insert data in event: ", err.Error())
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
		Select("id, date, name, description").
		From("event").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	event := new(entity.Event)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&event.ID, &event.Date, &event.Name, &event.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("event not found: ", err.Error())
			return fail(ErrNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return event, nil
}

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				log.Debug("couldn't insert data in user_event: ", err.Error())
				return fail(ErrAlreadyExists)
			default:
				log.Debug("couldn't insert data in user_event: ", err.Error())
				return fail(err)
			}
		} else {
			log.Debug("couldn't insert data in user_event: ", err.Error())
			return fail(err)
		}
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
