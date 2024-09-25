package event

import (
	"github.com/Slava02/Involvio/internal/entity"
	"time"
)

// Converters
func ToEventOutputFromEntity(event *entity.Event) *EventResponse {
	return &EventResponse{
		Body: struct{ entity.Event }{*event},
	}
}

type (
	CreateEventRequest struct {
		Body struct {
			SpaceId   int `json:"spaceId" example:"123" doc:"Space ID"`
			UserId    int `json:"userId" example:"123" doc:"Event ID"`
			EventInfo struct {
				Name        string      `json:"name" example:"fun event" doc:"Event name"`
				Description string      `json:"description" example:"enormously fun event" doc:"Event description"`
				BeginDate   time.Time   `json:"beginDate" example:"2007-03-01T13:00:00" doc:"Event start date and time"`
				EndDate     time.Time   `json:"endDate" example:"2007-03-01T13:00:00" doc:"Event end date and time"`
				Tags        entity.Tags `json:"tags" doc:"Tags for this event"`
			}
		}
	}

	EventByIdRequest struct {
		ID int `path:"id" json:"id" maxLength:"30" example:"1" doc:"event id"`
	}

	JoinEventRequest struct {
		EventId int `path:"id" json:"eventId" example:"123" doc:"Event ID"`
		Body    struct {
			UserId int `json:"userId" example:"123" doc:"User ID"`
		}
	}

	EventResponse struct {
		Body struct {
			entity.Event
		}
	}
)
