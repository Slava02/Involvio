package repo

import (
	"github.com/Slava02/Involvio/bot/internal/models"
	"sync"
)

type Storage struct {
	Mutex sync.RWMutex
	Data  map[int64]*models.User
}

func New() *Storage {
	return &Storage{
		Data: make(map[int64]*models.User),
	}
}

func (s *Storage) GetUser(id int64) *models.User {
	return s.Data[id]
}

func (s *Storage) UpdateUser(u *models.User) {
	s.Data[u.TelegID] = u
}
