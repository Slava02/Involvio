package event

import (
	"github.com/Slava02/Involvio/api/internal/entity"
	"time"
)

// Converters
func ToEventsOutputFromEntity(events []entity.Event) *UserEventsResponse {
	return &UserEventsResponse{
		Body: struct{ Events []entity.Event }{events},
	}
}

func ToEventOutputFromEntity(event entity.Event) *EventResponse {
	return &EventResponse{
		Body: struct{ entity.Event }{event},
	}
}

type (
	CreateEventRequest struct {
		Body struct {
			Name        string        `json:"name" doc:"event name" example:"random coffee"`
			Description string        `json:"description" doc:"event description" example:"super event"`
			Date        time.Time     `json:"date" doc:"event time RFC 3339" example:"2020-12-09T16:09:53+00:00"`
			Users       []entity.User `json:"users" doc:"event members"`
		}
	}

	EventByIdRequest struct {
		ID int `path:"id" json:"id" maxLength:"30" example:"1" doc:"event id"`
	}

	EventByUserIdRequest struct {
		ID int `path:"id" json:"id" maxLength:"30" example:"1" doc:"event id"`
	}

	JoinEventRequest struct {
		EventId int `path:"id" json:"event_id" example:"123" doc:"Event ID"`
		Body    struct {
			UserId int `json:"user_id" example:"123" doc:"User ID"`
		}
	}

	ReviewEventRequest struct {
		Body struct {
			EventID int `json:"event_id" doc:"event id" example:"123"`
			WhoID   int `json:"who_id" doc:"reviewer id" example:"123"`
			WhomID  int `json:"whom_id" doc:"reviewee id" example:"123"`
			Grade   int `json:"grade" doc:"grade of event" example:"123"`
		}
	}

	UserEventsResponse struct {
		Body struct {
			Events []entity.Event
		}
	}

	EventResponse struct {
		Body struct {
			entity.Event
		}
	}
)
