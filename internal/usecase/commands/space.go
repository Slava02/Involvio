package commands

import "github.com/Slava02/Involvio/internal/entity"

// SPACES
type (
	SpaceCommand struct {
		Name        string
		Description string
		Tags        entity.Tags
	}

	SpaceByIdCommand struct {
		ID int
	}

	JoinSpaceCommand struct {
		SpaceID int
		UserID  int
	}

	UpdateSpaceCommand struct {
		ID          int
		Name        string
		Description string
	}
)
