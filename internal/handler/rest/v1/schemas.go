package v1

import "github.com/Slava02/Involvio/internal/entity"

// USER
type (
	FindUserRequest struct {
		ID int `path:"id" maxLength:"30" example:"1" doc:"user id"`
	}

	UserRequest struct {
		Body struct {
			entity.User
		}
	}

	UpdateUserRequest struct {
		ID   int `path:"id" maxLength:"30" example:"1" doc:"user id"`
		Body struct {
			entity.User
		}
	}

	UpdateUserPrivilegesRequest struct {
		ID   int `path:"id" maxLength:"30" example:"1" doc:"user id"`
		Body struct {
			SpaceId int
		}
	}

	UserResponse struct {
		Body struct {
			*entity.User
		}
	}
)
