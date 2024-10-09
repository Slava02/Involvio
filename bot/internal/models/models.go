package models

// TODO: clean up models
const (
	DefaultSpace = "Общая"
)

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
	Groups      []string
	Photo
}

// Photo describes a submitted photo
type Photo struct {
	FileID string
}

func (u *User) GetSpaces() string {
	res := ""
	for i, v := range u.Groups {
		if i != 0 {
			res += "," + v
		} else {
			res += v
		}
	}
	return res
}
