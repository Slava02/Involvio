package usecase

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/usecase/commands"
)

type ISpaceRepository interface {
}

func NewSpaceUseCase(ur ISpaceRepository) *SpaceUseCase {
	return &SpaceUseCase{spaceRepo: ur}
}

type SpaceUseCase struct {
	spaceRepo ISpaceRepository
}

func (s SpaceUseCase) CreateSpace(ctx context.Context, cmd commands.SpaceCommand) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (s SpaceUseCase) GetSpace(ctx context.Context, cmd commands.SpaceByIdCommand) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (s SpaceUseCase) JoinSpace(ctx context.Context, cmd commands.JoinSpaceCommand) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}
