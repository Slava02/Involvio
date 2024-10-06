package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"log/slog"
	"sync"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func NewUserRepository(once *sync.Once, db *database.Postgres) *UserRepository {
	var repo *UserRepository
	once.Do(func() {
		repo = &UserRepository{db: db}
	})

	return repo
}

type UserRepository struct {
	db *database.Postgres
}

func (r *UserRepository) GetUserData(ctx context.Context, id int) (*entity.User, error) {
	const op = "Repo:GetUserData"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", id),
	)
	log.Debug(op)

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id, first_name, last_name, username, photo_url, auth_date").
		From("\"user\"").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	user := new(entity.User)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.PhotoURL, &user.AuthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found: ", err.Error())
			return fail(ErrUserNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return user, nil
}

func (r *UserRepository) GetUserForms(ctx context.Context, userId int) ([]*entity.Form, error) {
	const op = "Repo:GetUserForms"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", userId),
	)
	log.Debug(op)

	fail := func(err error) ([]*entity.Form, error) {
		return nil, fmt.Errorf("%r: %w", op, err)
	}

	// TODO: fix bug

	//query, args, err := r.db.Builder.
	//	Select("user_id, space_id, is_admin, is_creator, user_tags, pair_tags").
	//	From("user_space").
	//	Where("user_id = ?::int", userId).
	//	ToSql()
	//if err != nil {
	//	log.Debug("couldn't create SQL statement: ", err.Error())
	//	return fail(err)
	//}

	query := `SELECT space_id, is_admin, is_creator, user_tags, pair_tags FROM user_space WHERE user_id = $1`

	forms := make([]*entity.Form, 0)

	rows, err := r.db.Pool.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found: ", err.Error())
			return fail(ErrUserNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	for rows.Next() {
		form := new(entity.Form)

		err = rows.Scan(&form.SpaceID, &form.Admin, &form.Creator, &form.UserTags, &form.PairTags)
		if err != nil {
			return fail(err)
		}

		forms = append(forms, form)
	}

	return forms, nil
}

func (r *UserRepository) InsertUser(ctx context.Context, user *entity.User) error {
	const op = "Repo:InsertUser"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", user.ID),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("\"user\"").
		Columns("id, first_name, last_name, username, photo_url, auth_date").
		Values(user.ID, user.FirstName, user.LastName, user.UserName, user.PhotoURL, user.AuthDate).
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
				log.Debug("couldn't insert data in user: ", err.Error())
				return fail(ErrUserAlreadyExists)
			default:
				log.Debug("couldn't insert data in user: ", err.Error())
				return fail(err)
			}
		} else {
			log.Debug("couldn't insert data in user: ", err.Error())
			return fail(err)
		}
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, firstName, lastName, userName, photoURL string) (*entity.User, error) {
	const op = "Repo:UpdateUser"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Update("\"user\"").
		Set("first_name", firstName).
		Set("last_name", lastName).
		Set("username", userName).
		Set("photo_url", photoURL).
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

	user, err := r.GetUserData(ctx, id)
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId, spaceId int) error {
	const op = "Repo:DeleteUser"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", userId),
		slog.Int("space id", spaceId),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Delete("user_space").
		Where("user_id = ? AND space_id = ?", userId, spaceId).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't delete data from user_space: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *UserRepository) GetForm(ctx context.Context, userId, spaceId int) (*entity.Form, error) {
	const op = "Repo:GetForm"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", userId),
		slog.Int("space id", spaceId),
	)
	log.Debug(op)

	fail := func(err error) (*entity.Form, error) {
		return nil, fmt.Errorf("%r: %w", op, err)
	}

	//query, args, err := r.db.Builder.
	//	Select("user_id, space_id, is_admin, is_creator, user_tags, pair_tags").
	//	From("user_space").
	//	Where("user_id = ? AND space_id = ?", userId, spaceId).
	//	ToSql()
	//if err != nil {
	//	log.Debug("couldn't create SQL statement: ", err.Error())
	//	return fail(err)
	//}

	query := "SELECT space_id, is_admin, is_creator, user_tags, pair_tags FROM user_space WHERE user_id = $1 AND space_id = $2"

	form := new(entity.Form)

	err := r.db.Pool.QueryRow(ctx, query, userId, spaceId).Scan(&form.SpaceID, &form.Admin, &form.Creator, &form.UserTags, &form.PairTags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("form not found: ", err.Error())
			return fail(ErrUserNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return form, nil
}

func (r *UserRepository) UpdateForm(ctx context.Context, userId, spaceId int, userTags, pairTags entity.Tags) error {
	const op = "Repo:UpdateUser"

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Update("user_space").
		Set("user_tags", userTags).
		Set("pair_tags", pairTags).
		Where("user_id = ? AND space_id = ?", userId, spaceId).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't update form: ", err.Error())
		return fail(err)
	}

	return nil
}
