package entity

import (
	"time"
)

// User -.
type User struct {
	ID        int       `doc:"User ID" json:"id"       example:"1234"`
	FirstName string    `doc:"First name" json:"first_name"       example:"slava"`
	LastName  string    `doc:"Last name" json:"last_name"       example:"zhuvaga"`
	UserName  string    `doc:"Username" json:"user_name"       example:"s1av4"`
	PhotoURL  string    `doc:"Photo URL" json:"photo_url" example:"https://photo"`
	AuthDate  time.Time `doc:"Authorization date" json:"auth_date"       example:"25.09.2002 12:00"`
	Forms     []Form    `doc:"User info and preferences in certain space" json:"forms"`
}

type Form struct {
	SpaceID  int  `doc:"Space Id" json:"id"       example:"1234"`
	UserTags Tags `doc:"User's tags" json:"user_tags"`
	PairTags Tags `doc:"User's preference tags" json:"pair_tags"`
}
