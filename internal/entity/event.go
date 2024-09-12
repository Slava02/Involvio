package entity

import "google.golang.org/genproto/googleapis/type/datetime"

// Event -.
type Event struct {
	ID          int               `json:"id"       example:"1234"`
	Name        string            `json:"name"       example:"Гости"`
	Description string            `json:"description"       example:"Приглашаю всех в гости"`
	BeginDate   datetime.DateTime `json:"begin_date"       example:"25.09.2002 12:00"`
	EndDate     datetime.DateTime `json:"end_date"       example:"25.18.2002 18:00"`
}
