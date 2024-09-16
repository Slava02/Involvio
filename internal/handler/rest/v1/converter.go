package v1

import "github.com/Slava02/Involvio/internal/entity"

func ToUserOutputFromEntity(user *entity.User) UserResponse {
	return UserResponse{
		Body: struct{ *entity.User }{user},
	}
}
