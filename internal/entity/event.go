package entity

import "time"

// Event -.
type Event struct {
	ID          int       `json:"id"`
	SpaceId     int       `json:"space_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
	Tags        Tags      `json:"tags"`
}
