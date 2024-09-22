package usecase

import (
	"context"
	"fmt"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase/commands"
)

type ISpaceRepository interface {
	GetSpace(ctx context.Context, ID int) (*entity.Space, error)
	UpdateSpace(ctx context.Context, ID int, name, description string) (*entity.Space, error)
	DeleteSpace(ctx context.Context, ID int) error
	CreateSpace(ctx context.Context, name, description string, tags entity.Tags) (*entity.Space, error)
	AddUser(ctx context.Context, userId, spaceId int) error
}

func NewSpaceUseCase(ur ISpaceRepository) *SpaceUseCase {
	return &SpaceUseCase{spaceRepo: ur}
}

type SpaceUseCase struct {
	spaceRepo ISpaceRepository
}

func (sc *SpaceUseCase) UpdateSpace(ctx context.Context, cmd commands.UpdateSpaceCommand) (*entity.Space, error) {
	space, err := sc.spaceRepo.UpdateSpace(ctx, cmd.ID, cmd.Name, cmd.Description)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return space, nil
}

func (sc *SpaceUseCase) DeleteSpace(ctx context.Context, cmd commands.SpaceByIdCommand) error {
	err := sc.spaceRepo.DeleteSpace(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (sc *SpaceUseCase) CreateSpace(ctx context.Context, cmd commands.SpaceCommand) (*entity.Space, error) {
	space, err := sc.spaceRepo.CreateSpace(ctx, cmd.Name, cmd.Description, cmd.Tags)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return space, nil
}

func (sc *SpaceUseCase) GetSpace(ctx context.Context, cmd commands.SpaceByIdCommand) (*entity.Space, error) {
	space, err := sc.spaceRepo.GetSpace(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return space, nil
}

func (sc *SpaceUseCase) JoinSpace(ctx context.Context, cmd commands.JoinSpaceCommand) error {
	err := sc.spaceRepo.AddUser(ctx, cmd.SpaceID, cmd.UserID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
