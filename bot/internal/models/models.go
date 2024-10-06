package models

type User struct {
	TelegID     int
	FullName    string
	UserName    string
	PhotoURL    string
	Birthday    string
	Description string
	Gender      string
	City        string
	Socials     string
	Position    string
	Interests   string
	Goal        string
	Spaces      []string
}
