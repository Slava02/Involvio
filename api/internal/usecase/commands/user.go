package commands

import (
	"github.com/Slava02/Involvio/api/internal/entity"
	"time"
)

// USER
type (
	UserByIdCommand struct {
		ID int
	}

	UserByUsernameCommand struct {
		Username string
	}

	FormByIdCommand struct {
		UserID  int
		GroupID int
	}

	BlockUserCommand struct {
		WhoID  int
		WhomID int
	}

	SetHolidayCommand struct {
		ID       int
		TillDate time.Time
	}

	CancelHolidayCommand struct {
		ID int
	}

	UpdateUserCommand struct {
		ID        int
		FullName  string
		PhotoURL  string
		City      string
		Position  string
		Interests string
	}

	CreateUserCommand struct {
		User *entity.User
	}
)
