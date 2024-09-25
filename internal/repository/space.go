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
	ErrSpaceNotFound = errors.New("space not found")
)

func NewSpaceRepository(once *sync.Once, db *database.Postgres) *SpaceRepository {
	var repo *SpaceRepository
	once.Do(func() {
		repo = &SpaceRepository{db: db}
	})

	return repo
}

type SpaceRepository struct {
	db *database.Postgres
}

func (r *SpaceRepository) UpdateSpace(ctx context.Context, id int, name, description string) error {
	const op = "Repo:UpdateSpace"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Update("space").
		Set("name", name).
		Set("description", description).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't update space: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *SpaceRepository) GetSpace(ctx context.Context, id int) (*entity.Space, error) {
	const op = "Repo:GetSpace"

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", id),
	)
	log.Debug(op)

	fail := func(err error) (*entity.Space, error) {
		return nil, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id, name, description, tags").
		From("space").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	space := new(entity.Space)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&space.ID, &space.Name, &space.Description, &space.Tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("space not found: ", err.Error())
			return fail(ErrSpaceNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return space, nil
}

func (r *SpaceRepository) InsertSpace(ctx context.Context, userId int, space *entity.Space) error {
	const op = "Repo:InsertSpace"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	querySpace, argsSpace, err := r.db.Builder.
		Insert("space").
		Columns("id, name, description, tags").
		Values(space.ID, space.Name, space.Description, space.Tags).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	queryUserSpace, argsUserSpace, err := r.db.Builder.
		Insert("user_Space").
		Columns("user_id, space_id, is_admin, is_creator").
		Values(userId, space.ID, true, true).
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

	_, err = r.db.Pool.Exec(ctx, querySpace, argsSpace...)
	if err != nil {
		log.Debug("couldn't insert data in space: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, queryUserSpace, argsUserSpace...)
	if err != nil {
		log.Debug("couldn't insert data in space: ", err.Error())
		return fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Debug("couldn't commit transaction: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *SpaceRepository) DeleteSpace(ctx context.Context, id int) error {
	const op = "Repo:DeleteSpace"

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", id),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	querySpace, argsSpace, err := r.db.Builder.
		Delete("space").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	queryUserSpace, argsUserSpace, err := r.db.Builder.
		Delete("user_space").
		Where("space_id = ?", id).
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

	_, err = r.db.Pool.Exec(ctx, queryUserSpace, argsUserSpace...)
	if err != nil {
		log.Debug("couldn't delete data from user_space: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, querySpace, argsSpace...)
	if err != nil {
		log.Debug("couldn't delete data from space: ", err.Error())
		return fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Debug("couldn't commit transaction: ", err.Error())
		return fail(err)
	}

	return nil
}

// TODO: err ALready Exists
func (r *SpaceRepository) AddUser(ctx context.Context, userId, spaceId int) error {
	const op = "Repo:AddUserToSpace"

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", spaceId),
		slog.Int("user id", userId),
	)

	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("user_space").
		Columns("user_id, space_id, is_admin, is_creator").
		Values(userId, spaceId, false, false).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't insert data in user_space: ", err.Error())
		return fail(err)
	}

	return nil
}
