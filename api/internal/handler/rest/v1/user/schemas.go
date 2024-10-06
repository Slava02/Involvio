package user

import (
	"github.com/Slava02/Involvio/internal/entity"
)

// USER

// Converters
func ToUserOutputFromEntity(user *entity.User) *UserResponse {
	return &UserResponse{
		Body: struct{ *entity.User }{user},
	}
}

func ToFormOutputFromEntity(form *entity.Form) *FormResponse {
	return &FormResponse{
		Body: struct{ *entity.Form }{form},
	}
}

func ToUserWithFormsOutputFromEntity(user *entity.User, forms []*entity.Form) *UserWithFormsResponse {
	return &UserWithFormsResponse{
		Body: struct {
			*entity.User
			Forms []*entity.Form
		}{user, forms},
	}
}

// Schemas
type (
	UserByIdRequest struct {
		ID int `path:"id" maxLength:"30" example:"1" doc:"user id"`
	}

	DeleteUserRequest struct {
		UserId  int `path:"userId" maxLength:"30" example:"123" doc:"user id"`
		SpaceId int `path:"spaceId" maxLength:"30" example:"123" doc:"space id"`
	}

	CreateUserRequest struct {
		Body struct {
			FirstName string `json:"first_name" example:"ivan" doc:"User first name"`
			LastName  string `json:"last_name" example:"ivanov" doc:"User last nam"`
			Username  string `json:"username" example:"ivanko228" doc:"Username"`
			PhotoURL  string `json:"photo_url" example:"https://telegram/photo_ivan.png" doc:"User photo url"`
		}
	}

	UpdateUserRequest struct {
		ID   int `path:"id" maxLength:"30" example:"1" doc:"user id"`
		Body struct {
			FirstName string `json:"first_name" example:"ivan" doc:"User first name"`
			LastName  string `json:"last_name" example:"ivanov" doc:"User last nam"`
			Username  string `json:"username" example:"ivanko228" doc:"Username"`
			PhotoURL  string `json:"photo_url" example:"https://telegram/photo_ivan.png" doc:"User photo url"`
		}
	}

	FormByIdRequest struct {
		UserID  int `path:"userId" maxLength:"30" example:"1" doc:"user id"`
		SpaceID int `path:"spaceId" maxLength:"30" example:"1" doc:"space id"`
	}

	UpdateFormRequest struct {
		UserID  int `path:"userId" maxLength:"30" example:"1" doc:"user id"`
		SpaceID int `path:"spaceId" maxLength:"30" example:"1" doc:"space id"`
		Body    struct {
			UserTags entity.Tags
			PairTags entity.Tags
		}
	}

	UserResponse struct {
		Body struct {
			*entity.User
		}
	}

	UserWithFormsResponse struct {
		Body struct {
			*entity.User
			Forms []*entity.Form
		}
	}

	FormResponse struct {
		Body struct {
			*entity.Form
		}
	}
)
