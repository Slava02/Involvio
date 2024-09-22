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
			SpaceId int `json:"spaceId" example:"123" doc:"Space ID"`
			UserId  int `json:"userId" example:"123" doc:"Space ID"`
		}
	}

	CreateSpaceRequest struct {
		Body struct {
			Name        string      `json:"name" example:"MAI" doc:"Space Name"`
			Description string      `json:"description" example:"university" doc:"Space description"`
			Tags        entity.Tags `json:"tags" doc:"Tags options for this space"`
		}
	}

	UpdateSpaceRequest struct {
		ID   int `path:"id" maxLength:"30" example:"1" doc:"space id"`
		Body struct {
			Name        string `json:"name" example:"MAI" doc:"Space Name"`
			Description string `json:"description" example:"university" doc:"Space description"`
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
