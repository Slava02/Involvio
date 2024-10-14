package models

// TODO: clean up models
const (
	DefaultGroup = "Общая"
)

type User struct {
	TelegID   int64
	FullName  string
	UserName  string
	Birthday  string `validate:"date"`
	Gender    string
	City      string
	Socials   string `validate:"uri"`
	Position  string
	Interests string
	Goal      string
	Groups    []string
	Photo
}

// Photo describes a submitted photo
type Photo struct {
	FileID string
}
