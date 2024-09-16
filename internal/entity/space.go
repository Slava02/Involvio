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

type CreateSpaceRequest struct {
	AdminId int  `json:"admin_id"`
	Tags    Tags `json:"tags"`
}

type SpaceInfoReq struct {
	Space Space `json:"space"`
}

type JoinSpaceReq struct {
	SpaceId int  `json:"spaceId"`
	User    User `json:"user"`
}

type CreateSpaceResp struct {
	SpaceId int `json:"spaceId"`
}

type ListSpacesResp struct {
	Spaces []Space `json:"spaces"`
}

type Tags map[string]interface{}

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
