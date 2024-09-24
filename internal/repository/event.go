package repository

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"log"
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

func (r *EventRepository) InsertEvent(ctx context.Context, userId int, event *entity.Event) error {
	log.Println("event: ", event.ID)

	fail := func(err error) error {
		return fmt.Errorf("InsertEvent: %w", err)
	}

	sqlEvent, argsEvent, err := r.db.Builder.
		Insert("event").
		Columns("id, space_id, name, description, begin_date, end_date, tags").
		Values(event.ID, event.SpaceId, event.Name, event.Description, event.BeginDate, event.EndDate, event.Tags).
		ToSql()
	if err != nil {
		return fail(fmt.Errorf("couldn't create SQL statement"))
	}

	sqlUserEvent, argsUserEvent, err := r.db.Builder.
		Insert("user_event").
		Columns("user_id, event_id").
		Values(userId, event.ID).
		ToSql()
	if err != nil {
		return fail(fmt.Errorf("couldn't create SQL statement"))
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Pool.Exec(ctx, sqlEvent, argsEvent...)
	if err != nil {
		return fail(fmt.Errorf("couldn't insert data in event"))
	}

	_, err = r.db.Pool.Exec(ctx, sqlUserEvent, argsUserEvent...)
	if err != nil {
		return fail(fmt.Errorf("couldn't insert data in user_event"))
	}

	if err = tx.Commit(ctx); err != nil {
		return fail(fmt.Errorf("couldn't commit transaction"))
	}

	return nil
}

func (r *EventRepository) GetEvent(ctx context.Context, id int) (*entity.Event, error) {
	fail := func(err error) (*entity.Event, error) {
		return nil, fmt.Errorf("GetEvent: %w", err)
	}

	sql, _, err := r.db.Builder.
		Select("id, space_id, name, description, begin_date, end_date, tags").
		From("event").
		ToSql()
	if err != nil {
		return fail(fmt.Errorf("couldn't create SQL statement"))
	}

	event := new(entity.Event)

	err = r.db.Pool.QueryRow(ctx, sql).Scan(&event.ID, &event.SpaceId, &event.Name, &event.Description, &event.BeginDate, &event.EndDate, &event.Tags)
	if err != nil {
		log.Println("error: ", err.Error())
		return fail(fmt.Errorf("couldn't get event"))
	}

	return event, nil
}

func (r *EventRepository) AddUser(ctx context.Context, eventId, userId int) error {
	log.Println("event: ", eventId)

	fail := func(err error) error {
		return fmt.Errorf("InsertEvent: %w", err)
	}

	sql, args, err := r.db.Builder.
		Insert("user_event").
		Columns("user_id, event_id").
		Values(userId, eventId).
		ToSql()
	if err != nil {
		return fail(fmt.Errorf("couldn't create SQL statement"))
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Println("error: ", err.Error())
		return fail(fmt.Errorf("couldn't insert data in user_event"))
	}

	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, id int) error {
	fail := func(err error) error {
		return fmt.Errorf("InsertEvent: %w", err)
	}

	sqlEvent, argsEvent, err := r.db.Builder.
		Delete("event").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fail(fmt.Errorf("couldn't create SQL statement"))
	}

	sqlUserEvent, argsUserEvent, err := r.db.Builder.
		Delete("user_event").
		Where("event_id = ?", id).
		ToSql()
	if err != nil {
		return fail(fmt.Errorf("couldn't create SQL statement"))
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Pool.Exec(ctx, sqlUserEvent, argsUserEvent...)
	if err != nil {
		return fail(fmt.Errorf("couldn't delete data from user_event"))
	}

	_, err = r.db.Pool.Exec(ctx, sqlEvent, argsEvent...)
	if err != nil {
		log.Println("error: ", err)
		return fail(fmt.Errorf("couldn't delete data from event"))
	}

	if err = tx.Commit(ctx); err != nil {
		return fail(fmt.Errorf("couldn't commit transaction"))
	}

	return nil
}
