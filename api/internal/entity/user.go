package entity

import "time"

// User -.
type User struct {
	TelegID       int    `doc:"Telegram ID" json:"teleg_id"       example:"1234"`
	FirstName     string `doc:"First name" json:"first_name"       example:"slava"`
	LastName      string `doc:"Last name" json:"last_name"       example:"zhuvaga"`
	UserName      string `doc:"Username" json:"user_name"       example:"s1av4"`
	PhotoURL      string `doc:"Photo URL" json:"photo_url" example:"https://photo"`
	Birthday      string `doc:"Birthday" json:"birthday"       example:"25.09.2002"`
	Description   string `doc:"User description" json:"description" example:"I am user"`
	Gender        string `doc:"User's gender" json:"gender" example:"male"`
	City          string `doc:"User's city" json:"city" example:"Moscow"`
	Socials       string `doc:"User's account" json:"socials" exmaple:"https://vk.com"`
	Position      string `doc:"User's position in organization" json:"position" example:"student"`
	Interests     string `doc:"User's interests" json:"interests" example:"Programming,math"`
	MeetingFormat string `doc:"Preferable meeting format" json:"meeting_format" example:"online"`
	Goal          string `doc:"User's goal for the meetings" json:"goal" example:"fun"`
}

type Holiday struct {
	Status   bool      `doc:"Whether holiday is active or not" json:"status" example:"true"`
	TillDate time.Time `doc:"When holiday ends" json:"till_date" example:"25.09.2002"`
}
