package repository

import (
	"context"
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/pkg/database"
	"sync"
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

func (s *SpaceRepository) GetSpace(ctx context.Context, ID int) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceRepository) UpdateSpace(ctx context.Context, ID int, name, description string) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceRepository) DeleteSpace(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceRepository) CreateSpace(ctx context.Context, name, description string, tags entity.Tags) (*entity.Space, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SpaceRepository) AddUser(ctx context.Context, userId, spaceId int) error {
	//TODO implement me
	panic("implement me")
}
