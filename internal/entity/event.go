package entity

import "time"

// Event -.
type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
	Tags        Tags      `json:"tags"`
}

type CreateEventReq struct {
	UserId  int   `json:"user_id"`
	SpaceId int   `json:"space_id"`
	Event   Event `json:"event"`
}

type JoinEventReq struct {
	EventId int `json:"event_id"`
	UserId  int `json:"user_id"`
}

type CreateEventResp struct {
	EventId int `json:"event_id"`
}

type EventInfoResp struct {
	ID    int    `json:"id"`
	Event Event  `json:"event"`
	Users []User `json:"users"`
}
