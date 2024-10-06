package commands

import (
	"github.com/Slava02/Involvio/internal/entity"
	"time"
)

// EVENTS
type (
	CreateEventCommand struct {
		SpaceId     int
		UserId      int
		Name        string
		Description string
		BeginDate   time.Time
		EndDate     time.Time
		Tags        entity.Tags
	}

	EventByIdCommand struct {
		ID int `path:"id" maxLength:"30" example:"1" doc:"event id"`
	}

	JoinEventCommand struct {
		EventId int `json:"eventId" example:"123" doc:"Event ID"`
		UserId  int `json:"userId" example:"123" doc:"User ID"`
	}
)
