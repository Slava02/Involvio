package user

import (
	"github.com/Slava02/Involvio/api/internal/entity"
	"time"
)

// USER

// Converters
func ToUserOutputFromEntity(user *entity.User) *UserRequestResponse {
	return &UserRequestResponse{
		Body: struct{ *entity.User }{user},
	}
}

// Schemas
type (
	UserByUsernameRequest struct {
		Username string `path:"username" json:"username" example:"ivanko228" doc:"Username"`
	}

	UserByIDRequest struct {
		Body struct {
			ID int `json:"id" example:"123" doc:"ID"`
		}
	}

	//  TODO: как когда несколько параметров, а изменены могут быть только нескольо или один
	UpdateUserRequest struct {
		Body struct {
			ID        int    `doc:"Telegram ID" json:"id"       example:"1234"`
			FullName  string `doc:"First name" json:"full_name"       example:"ivan popkins"`
			PhotoURL  string `doc:"Photo URL" json:"photo_url" example:"https://photo"`
			City      string `doc:"User's city" json:"city" example:"Moscow"`
			Position  string `doc:"User's position in organization" json:"position" example:"student"`
			Interests string `doc:"User's interests" json:"interests" example:"Programming,math"`
		}
	}

	BlockUserRequest struct {
		Body struct {
			WhoID  int `doc:"Telegram ID of user who is blocking" json:"who_id"       example:"1234"`
			WhomID int `doc:"Telegram ID of user who is being blocked" json:"whom_id"       example:"4321"`
		}
	}

	SetHolidayRequest struct {
		Body struct {
			ID       int       `doc:"Telegram ID" json:"id"       example:"1234"`
			TillDate time.Time `doc:"When holiday ends RFC 3339" json:"till_date" example:"2020-12-09T16:09:53+00:00"`
		}
	}

	UserRequestResponse struct {
		Body struct {
			*entity.User
		}
	}
)
