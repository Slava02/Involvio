package repository

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/labstack/gommon/log"
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
		log.Errorf("error: %w\n", err)
		return fail(fmt.Errorf("couldn't insert data in event table"))
	}

	_, err = r.db.Pool.Exec(ctx, sqlUserEvent, argsUserEvent...)
	if err != nil {
		return fail(fmt.Errorf("couldn't insert data in user_event table"))
	}

	if err = tx.Commit(ctx); err != nil {
		return fail(fmt.Errorf("couldn't commit transaction"))
	}

	return nil
}

func (r *EventRepository) GetEvent(ctx context.Context, id int) (*entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (r *EventRepository) AddUser(ctx context.Context, eventId, userId int) error {
	//TODO implement me
	panic("implement me")
}

func (r *EventRepository) DeleteEvent(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
