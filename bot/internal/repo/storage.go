package repo

import (
	"github.com/Slava02/Involvio/bot/internal/models"
	"strings"
	"sync"
)

// TODO: implement Redis local Percistance
var groups = []string{"Slava", "Test", models.DefaultSpace}

type Storage struct {
	Mutex sync.RWMutex
	Data  map[string]*models.User
}

func New() *Storage {
	return &Storage{
		Data: make(map[string]*models.User),
	}
}

//  TODO: refactor storage

func (s *Storage) GetUser(username string) *models.User {
	return s.Data[username]
}

func (s *Storage) UpdateUser(u *models.User) {
	s.Data[u.UserName] = u
}

func (s *Storage) AddGroups(username string, groups string) string {
	g := strings.Split(groups, ",")

	newGroups := make([]string, 0)

	for _, group := range g {
		if ok := s.GetGroup(group); !ok {
			s.CreateGroup(group)
			newGroups = append(newGroups, group)
		} else {
			s.JoinGroup(username, group)
		}
	}

	if groups[len(groups)-1] == ',' {
		s.JoinGroup(username, models.DefaultSpace)
	}

	var groupList string
	if newGroups != nil {
		if len(newGroups) > 1 {
			groupList = strings.Join(newGroups, ",")
		} else {
			groupList = newGroups[0]
		}
	}

	return groupList
}

func (s *Storage) JoinGroup(username string, g string) {
	s.Data[username].Groups = append(s.Data[username].Groups, g)
}

func (s *Storage) CreateGroup(g string) {
	groups = append(groups, g)
}

func (s *Storage) GetGroup(g string) bool {
	m := make(map[string]struct{})
	for _, v := range groups {
		m[v] = struct{}{}
	}

	if _, ok := m[g]; !ok {
		return false
	}

	return true
}
