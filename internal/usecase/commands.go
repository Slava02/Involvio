package usecase

import "github.com/Slava02/Involvio/internal/entity"

// USER
type (
	FindUserByIDCommand struct {
		ID int
	}

	CreateUpdateUserCommand struct {
		User entity.User
	}

	DeleteUserByIDCommand struct {
		ID int
	}
)
