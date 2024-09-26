package commands

import (
	"github.com/Slava02/Involvio/internal/entity"
	"time"
)

// USER
type (
	UserByIdCommand struct {
		ID int
	}

	FormByIdCommand struct {
		UserID  int
		SpaceID int
	}

	UpdateFormCommand struct {
		UserID   int
		SpaceID  int
		UserTags entity.Tags
		PairTags entity.Tags
	}

	UpdateUserCommand struct {
		ID        int
		FirstName string
		LastName  string
		UserName  string
		PhotoURL  string
	}

	CreateUserCommand struct {
		FirstName string
		LastName  string
		UserName  string
		PhotoURL  string
		AuthDate  time.Time
	}
)
