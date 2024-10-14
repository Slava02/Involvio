package commands

import (
	"github.com/Slava02/Involvio/api/internal/entity"
	"time"
)

// EVENTS
type (
	CreateEventCommand struct {
		Name        string
		Description string
		Date        time.Time
		Users       []entity.User
	}

	EventByUserIdCommand struct {
		ID int
	}

	EventByIdCommand struct {
		ID int
	}

	JoinEventCommand struct {
		EventId int
		UserId  int
	}

	ReviewEventCommand struct {
		EventID int
		WhoID   int
		WhomID  int
		Grade   int
	}
)
