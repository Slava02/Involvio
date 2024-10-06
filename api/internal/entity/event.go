package entity

import "time"

// Event -.
type Event struct {
	ID          int       `json:"id" doc:"event id" example:"123"`
	Name        string    `json:"name" doc:"event name" example:"random coffee"`
	Description string    `json:"description" doc:"event description" example:"super event"`
	Date        time.Time `json:"date" doc:"event time" example:"05.09.2002"`
	Users       []User    `json:"users" doc:"event members"`
}

type Review struct {
	ID      int `json:"id" doc:"review id" example:"123"`
	EventID int `json:"event_id" doc:"event id" example:"123"`
	WhoID   int `json:"who_id" doc:"reviewer id" example:"123"`
	WhomID  int `json:"whom_id" doc:"reviewee id" example:"123"`
	Grade   int `json:"grade" doc:"grade of event" example:"123"`
}
