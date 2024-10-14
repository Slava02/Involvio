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

func NewGroupRepository(once *sync.Once, db *database.Postgres) *GroupRepository {
	var repo *GroupRepository
	once.Do(func() {
		repo = &GroupRepository{db: db}
	})

	return repo
}

type GroupRepository struct {
	db *database.Postgres
}

// TODO: change space to group
func (r *GroupRepository) RemoveUser(ctx context.Context, userId int, groupName string) error {
	const op = "Repo:RemoveUser"

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", userId),
		slog.String("groupName", groupName),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	mainQuery := r.db.Builder.
		Delete("user_space").
		Where(sq.Eq{"user_id": userId})

	subQuery := r.db.Builder.
		Select("id").
		From("space").
		Where(sq.Eq{"name": groupName})

	mainQuery = mainQuery.Where(subQuery.Prefix("id IN (").Suffix(")"))

	query, args, err := mainQuery.ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't remove user from group: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *GroupRepository) GetGroup(ctx context.Context, name string) (*entity.Group, error) {
	const op = "Repo:GetGroup"

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", name),
	)
	log.Debug(op)

	fail := func(err error) (*entity.Group, error) {
		return nil, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id, name").
		From("space").
		Where("name = ?", name).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	group := new(entity.Group)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&group.ID, &group.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("group not found: ", err.Error())
			return fail(ErrNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return group, nil
}

func (r *GroupRepository) InsertGroup(ctx context.Context, group *entity.Group) error {
	const op = "Repo:InsertGroup"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("space").
		Columns("id, name").
		Values(group.ID, group.Name).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't insert data in group: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *GroupRepository) AddUser(ctx context.Context, userId int, groupName string) error {
	const op = "Repo:AddUserToGroup"

	log := slog.With(
		slog.String("op", op),
		slog.String("group id", groupName),
		slog.Int("user id", userId),
	)

	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "INSERT INTO user_space (user_id, space_id) SELECT $1 user_id, id space_id FROM space WHERE name = $2"

	_, err := r.db.Pool.Exec(ctx, query, userId, groupName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				log.Debug("couldn't insert data in user_group: ", err.Error())
				return fail(ErrAlreadyExists)
			default:
				log.Debug("couldn't insert data in user_group: ", err.Error())
				return fail(err)
			}
		} else {
			log.Debug("couldn't insert data in user_group: ", err.Error())
			return fail(err)
		}
	}

	return nil
}

func (r *GroupRepository) DeleteGroup(ctx context.Context, name string) error {
	const op = "Repo:DeleteGroup"

	log := slog.With(
		slog.String("op", op),
		slog.String("groupName", name),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	mainUserSpaceQuery := r.db.Builder.
		Delete("user_space")

	subUserSpaceQuery := r.db.Builder.
		Select("id").
		From("space").
		Where(sq.Eq{"name": name})

	mainUserSpaceQuery = mainUserSpaceQuery.Where(subUserSpaceQuery.Prefix("id ="))

	queryUserSpace, argsUserSpace, err := mainUserSpaceQuery.ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	querySpace, argsSpace, err := r.db.Builder.
		Delete("space").
		Where("name = ?", name).
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
		log.Debug("couldn't delete data from user_group: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, querySpace, argsSpace...)
	if err != nil {
		log.Debug("couldn't delete data from user_group: ", err.Error())
		return fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Debug("couldn't commit transaction: ", err.Error())
		return fail(err)
	}

	return nil
}
