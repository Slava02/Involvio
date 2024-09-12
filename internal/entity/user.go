package entity

import "google.golang.org/genproto/googleapis/type/datetime"

// User -.
type User struct {
	ID        int               `json:"id"       example:"1234"`
	FirstName string            `json:"first_name"       example:"slava"`
	LastName  string            `json:"last_name"       example:"zhuvaga"`
	UserName  string            `json:"user_name"       example:"s1av4"`
	PhotoURL  string            `json:"photo_url" example:"https://photo"`
	AuthDate  datetime.DateTime `json:"auth_date"       example:"25.09.2002 12:00"`
}
