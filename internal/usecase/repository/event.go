package repository

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"sync"
	"time"
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

func (e *EventRepository) InsertEvent(ctx context.Context, userId, spaceId int, name, description string, tags entity.Tags, beginDate, endDate time.Time) (*entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepository) GetEvent(ctx context.Context, id int) (*entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepository) AddUser(ctx context.Context, eventId, userId int) error {
	//TODO implement me
	panic("implement me")
}

func (e *EventRepository) DeleteEvent(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
