package repository

import (
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"sync"
)

func NewRepository(once *sync.Once, pg *database.Postgres) *Repository {
	var repo *Repository
	once.Do(func() {
		repo = &Repository{pg}
	})

	return repo
}

type Repository struct {
	*database.Postgres
}

func (r Repository) CreateEvent(userId, spaceId int64, event *entity.Event) (*entity.Event, error) {
	//sqlInsertEvent, args, err := r.Builder.
	//	Insert("event").
	//	Columns("name, description, begin_date, end_date, tags").
	//	Values(event.Name, event.Description, event.BeginDate, event.EndDate).
	//	ToSql()
	//if err != nil {
	//	return nil, fmt.Errorf("%w", err)
	//}
	//
	//sqlInsertUserEvent, args, err := r.Builder.
	//	Insert("user_event").
	//	Columns("user_id, event_id").
	//	Values(userId, spaceId).
	//	ToSql()
	//if err != nil {
	//	return nil, fmt.Errorf("%w", err)
	//}
	//
	//sqlInsertUserEvent, args, err := r.Builder.
	//	Insert("user_event").
	//	Columns("user_id, event_id").
	//	Values(userId, spaceId).
	//	ToSql()
	//if err != nil {
	//	return nil, fmt.Errorf("%w", err)
	//}
	//
	//tx, err := r.Postgres.Pool.BeginTx(context.TODO(), pgx.TxOptions{})
	//if err != nil {
	//	return nil, fmt.Errorf("%w", err)
	//}
	//defer tx.Rollback(context.TODO())
	//
	//_, err = r.Pool.Exec(context.Background(), sqlInsertEvent, args...)
	//if err != nil {
	//	return nil, fmt.Errorf("%w", err)
	//}
	//
	//if err = tx.Commit(context.TODO()); err != nil {
	//	return nil, fmt.Errorf("%w", err)
	//}
	//
	//return nil
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetEventById(eventId int64) (*entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetEventUsers(eventId int64) ([]*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) AddUserToEvent(eventId, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) CreateSpace(spaceName, spaceDescription string, adminId int64, tags entity.TagOptions) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetSpaceById(spaceId int64) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) AddUserToSpace(spaceId, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) CreateUser(user *entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) DeleteUser(userId, spaceId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetUserById(userId int64) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetFormsByUserId(userId int64) ([]*entity.Form, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) UpdateUser(userId int64, form *entity.Form) error {
	//TODO implement me
	panic("implement me")
}
