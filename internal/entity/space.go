package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Space -.
type Space struct {
	ID          int    `json:"id"       example:"1234"`
	Name        string `json:"name"       example:"mai"`
	Description string `json:"description"       example:"university space"`
	Tags        Tags   `json:"tags"`
}

type Tags []map[string]interface{}

func (a Tags) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Tags) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
