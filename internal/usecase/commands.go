package usecase

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
		*entity.Form
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
