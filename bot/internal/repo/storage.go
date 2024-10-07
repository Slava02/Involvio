package repo

import (
	"github.com/Slava02/Involvio/bot/internal/models"
	"sync"
)

type Storage struct {
	mutex sync.RWMutex
	data  map[int64]models.User
}
