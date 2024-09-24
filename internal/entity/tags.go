package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Tags []map[string]interface{}

func (a *Tags) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Tags) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal([]byte(b), &a)
}
