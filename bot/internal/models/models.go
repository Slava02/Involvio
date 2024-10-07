package models

type User struct {
	TelegID     int64
	FullName    string
	UserName    string
	Birthday    string `validate:"date"`
	Description string
	Gender      string
	City        string
	Socials     string `validate:"uri"`
	Position    string
	Interests   string
	Goal        string
	Spaces      []string
	Photo
}

// Photo describes a submitted photo
type Photo struct {
	ID          int
	FileID      string
	Description string
}
