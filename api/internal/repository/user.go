package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/pkg/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"log/slog"
	"sync"
	"time"
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

func (r *UserRepository) InsertUser(ctx context.Context, user *entity.User) error {
	const op = "Repo:InsertUser"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", user.ID),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("\"user\"").
		// TODO: birthday пока строка
		// TODO: поменять sex в БД yf gender
		Columns("id, full_name, username, birthday, photo_url, city, socials, position, sex, interests, goal").
		Values(user.ID, user.FullName, user.UserName, user.Birthday, user.PhotoURL, user.City, user.Socials, user.Position, user.Gender, user.Interests, user.Goal).
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
				return fail(ErrAlreadyExists)
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

func (r *UserRepository) BlockUser(ctx context.Context, who, whom int) error {
	const op = "Repo:BlockUser"

	log := slog.With(
		slog.String("op", op),
		slog.Int("who", who),
		slog.Int("who", whom),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("blocks").
		// TODO: change to who and whom
		Columns("id, user_id").
		Values(who, whom).
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
				log.Debug("couldn't insert data in blocks: ", err.Error())
				return fail(ErrAlreadyExists)
			default:
				log.Debug("couldn't insert data in blocks: ", err.Error())
				return fail(err)
			}
		} else {
			log.Debug("couldn't insert data in blocks: ", err.Error())
			return fail(err)
		}
	}

	return nil
}

func (r *UserRepository) SetHoliday(ctx context.Context, id int, tillDate time.Time) (*entity.User, error) {
	const op = "Repo:SetHoliday"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", id),
		slog.String("tillDate", tillDate.String()),
	)
	log.Debug(op)

	fail := func(err error) (*entity.User, error) {
		return &entity.User{}, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Insert("holiday_status").
		Columns("id, status, till_date").
		Values(id, true, tillDate).
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
				log.Debug("couldn't insert data in holiday_status: ", err.Error())
				return fail(ErrAlreadyExists)
			default:
				log.Debug("couldn't insert data in holiday_status: ", err.Error())
				return fail(err)
			}
		} else {
			log.Debug("couldn't insert data in holiday_status: ", err.Error())
			return fail(err)
		}
	}

	return nil, nil
}

// TODO: надо добавить id самого hpoliday, поменять UserID и добавить дату когда его поставили
func (r *UserRepository) CancelHoliday(ctx context.Context, id int) error {
	const op = "Repo:CancelHoliday"

	log := slog.With(
		slog.String("op", op),
		slog.Int("user id", id),
	)
	log.Debug(op)

	fail := func(err error) error {
		return fmt.Errorf("%r: %w", op, err)
	}

	//  Tут легче руками написать - UPDATE holiday_status SET status WHERE holiday_id = (SELECT holiday_id FROM holidays_status WHERE user_id = ? and set_date = (SELECT MAX(set_date) FROM holiday_status WHERE user_id = ?) FOR UPDATE)

	mainQuery, args, err := r.db.Builder.
		Update("\"holiday_status\"").
		Set("status", false).
		Where("user_id=?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	subQuery := r.db.Builder.
		Select("holiday_id, MAX(date)").
		Where("")

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't update group: ", err.Error())
		return fail(err)
	}

	return nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	const op = "Repo:GetUserByUsername"

	log := slog.With(
		slog.String("op", op),
		slog.String("username", username),
	)
	log.Debug(op)

	fail := func(err error) (*entity.User, error) {
		return &entity.User{}, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id, username, full_name, birthday, gender, city, socials, position, sex, photo_url, interests, goal").
		From("\"user\"").
		Where("username = ?", username).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	user := new(entity.User)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.UserName, &user.FullName, &user.Birthday, &user.Gender, &user.City, &user.Socials, &user.Position, &user.Gender, &user.PhotoURL, &user.Interests, &user.Goal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found: ", err.Error())
			return fail(ErrNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return nil, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	const op = "Repo:GetUserByID"

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", id),
	)
	log.Debug(op)

	fail := func(err error) (*entity.User, error) {
		return &entity.User{}, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id, username, full_name, birthday, gender, city, socials, position, sex, photo_url, interests, goal").
		From("\"user\"").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	user := new(entity.User)

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.UserName, &user.FullName, &user.Birthday, &user.Gender, &user.City, &user.Socials, &user.Position, &user.Gender, &user.PhotoURL, &user.Interests, &user.Goal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found: ", err.Error())
			return fail(ErrNotFound)
		} else {
			log.Debug("error: ", err.Error())
			return fail(err)
		}
	}

	return nil, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, fullName, city, position, interests, photoURL string) (*entity.User, error) {
	const op = "Repo:UpdateUser"

	log := slog.With(
		slog.String("op", op),
		slog.Int("userID", id),
	)
	log.Debug(op)

	fail := func(err error) (*entity.User, error) {
		return &entity.User{}, fmt.Errorf("%r: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Update("\"user\"").
		Set("full_name", fullName).
		Set("city", city).
		Set("position", position).
		Set("interests", interests).
		Set("photo_url", photoURL).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Debug("couldn't create SQL statement: ", err.Error())
		return fail(err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		log.Debug("couldn't update group: ", err.Error())
		return fail(err)
	}

	//  TODO: тупость - это надо в юзкейз, а тут ничего не возвращаем
	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		log.Debug("couldn't get user: ", err.Error())
		return fail(err)
	}

	return user, nil
}
