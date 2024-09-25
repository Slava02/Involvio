package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase/commands"
	"github.com/Slava02/Involvio/pkg/hexid"
	"log/slog"
)

type ISpaceRepository interface {
	GetSpace(ctx context.Context, id int) (*entity.Space, error)
	UpdateSpace(ctx context.Context, id int, name, description string) error
	DeleteSpace(ctx context.Context, id int) error
	InsertSpace(ctx context.Context, userId int, space *entity.Space) error
	AddUser(ctx context.Context, userId, spaceId int) error
}

func NewSpaceUseCase(ur ISpaceRepository) *SpaceUseCase {
	return &SpaceUseCase{spaceRepo: ur}
}

type SpaceUseCase struct {
	spaceRepo ISpaceRepository
}

func (sc *SpaceUseCase) UpdateSpace(ctx context.Context, cmd commands.UpdateSpaceCommand) (*entity.Space, error) {
	const op = "Usecase:UpdateSpace"

	fail := func(err error) (*entity.Space, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, err := sc.GetSpace(ctx, commands.SpaceByIdCommand{ID: cmd.ID})
	if err != nil {
		log.Debug("couldn't get space: ", err.Error())
		return fail(err)
	}

	err = sc.spaceRepo.UpdateSpace(ctx, cmd.ID, cmd.Name, cmd.Description)
	if err != nil {
		return fail(err)
	}

	return &entity.Space{
		ID:          cmd.ID,
		Name:        cmd.Name,
		Description: cmd.Description,
	}, nil
}

func (sc *SpaceUseCase) DeleteSpace(ctx context.Context, cmd commands.SpaceByIdCommand) error {
	const op = "Usecase:DeleteSpace"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, err := sc.GetSpace(ctx, commands.SpaceByIdCommand{ID: cmd.ID})
	if err != nil {
		log.Debug("couldn't get space: ", err.Error())
		return fail(err)
	}

	err = sc.spaceRepo.DeleteSpace(ctx, cmd.ID)
	if err != nil {
		log.Debug("couldn't delete space: ", err.Error())
		return fail(err)
	}

	return nil
}

func (sc *SpaceUseCase) CreateSpace(ctx context.Context, cmd commands.CreateSpaceCommand) (*entity.Space, error) {
	const op = "Usecase:CreateSpace"

	fail := func(err error) (*entity.Space, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	spaceId, err := hexid.Generate()
	if err != nil {
		log.Error("couldn't generate id: ", err.Error())
		return fail(err)
	}

	space := &entity.Space{
		ID:          spaceId,
		Name:        cmd.Name,
		Description: cmd.Description,
		Tags:        cmd.Tags,
	}

	err = sc.spaceRepo.InsertSpace(ctx, cmd.UserID, space)
	if err != nil {
		log.Debug("couldn't insert space: ", err.Error())
		return fail(err)
	}

	return space, nil
}

func (sc *SpaceUseCase) GetSpace(ctx context.Context, cmd commands.SpaceByIdCommand) (*entity.Space, error) {
	const op = "Usecase:GetSpace"

	fail := func(err error) (*entity.Space, error) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
		slog.Int("space id", cmd.ID),
	)
	log.Debug(op)

	space, err := sc.spaceRepo.GetSpace(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, repository.ErrSpaceNotFound) {
			log.Info("couldn't get event: ", err.Error())
		} else {
			log.Debug("couldn't get event: ", err.Error())
		}
		return fail(err)
	}

	return space, nil
}

// TODO: already joined
// TODO: fix
func (sc *SpaceUseCase) JoinSpace(ctx context.Context, cmd commands.JoinSpaceCommand) error {
	const op = "Usecase:JoinSpace"

	fail := func(err error) error {
		return fmt.Errorf("%s: %w", op, err)
	}

	log := slog.With(
		slog.String("op", op),
	)
	log.Debug(op)

	_, err := sc.GetSpace(ctx, commands.SpaceByIdCommand{ID: cmd.SpaceID})
	if err != nil {
		log.Debug("couldn't get space: ", err.Error())
		return fail(err)
	}

	err = sc.spaceRepo.AddUser(ctx, cmd.UserID, cmd.SpaceID)
	if err != nil {
		return fail(err)
	}

	return nil
}
