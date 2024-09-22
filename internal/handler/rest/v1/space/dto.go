package space

import "github.com/Slava02/Involvio/internal/entity"

// Converters
func ToSpaceOutputFromEntity(space *entity.Space) *SpaceResponse {
	return &SpaceResponse{
		Body: struct{ entity.Space }{*space},
	}
}

type (
	JoinSpaceRequest struct {
		Body struct {
			SpaceId int
			UserId  int
		}
	}

	CreateSpaceRequest struct {
		Body struct {
			Name        string
			Description string
			Tags        entity.Tags
		}
	}

	UpdateSpaceRequest struct {
		ID   int `path:"id" maxLength:"30" example:"1" doc:"space id"`
		Body struct {
			Name        string
			Description string
		}
	}

	SpaceByIdRequest struct {
		ID int `path:"id" maxLength:"30" example:"1" doc:"space id"`
	}

	SpaceResponse struct {
		Body struct {
			entity.Space
		}
	}
)
