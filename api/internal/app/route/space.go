package route

import (
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/handler/rest/v1/space"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/Slava02/Involvio/pkg/valid"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"sync"
)

type SpaceDeps struct {
	Validator *valid.Validator
}

//nolint:funlen
func setupSpaceRoutes(api huma.API, pg *database.Postgres, deps *SpaceDeps) {
	o := sync.Once{}
	spaceUseCase := usecase.NewSpaceUseCase(repository.NewSpaceRepository(&o, pg))

	spaceHandler := space.NewSpaceHandler(spaceUseCase)

	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)
	spaceSchema := huma.SchemaFromType(registry, reflect.TypeOf(&entity.Space{}))

	huma.Register(api, huma.Operation{
		OperationID:   "CreateSpace",
		Method:        http.MethodPost,
		Path:          "/spaces",
		Summary:       "create new space",
		Description:   "Create a new space record.",
		Tags:          []string{"Spaces"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "ISpaceUC created",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: spaceSchema,
					},
				},
				Headers: map[string]*huma.Param{
					"Location": {
						Description: "URL of the newly created space",
						Schema:      &huma.Schema{Type: "string"},
						Required:    true,
					},
				},
			},
			"400": {
				Description: "Invalid request",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"message": {Type: "string"},
								"field":   {Type: "string"},
							},
						},
					},
				},
			},
			"500": {
				Description: "Internal server error",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"error": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}, spaceHandler.CreateSpace)

	huma.Register(api, huma.Operation{
		OperationID: "GetSpace",
		Method:      http.MethodGet,
		Path:        "/spaces/{id}",
		Summary:     "space by id",
		Description: "Get space by id.",
		Tags:        []string{"Spaces"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "ISpaceUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: spaceSchema,
					},
				},
			},
			"400": {
				Description: "Invalid request",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"message": {Type: "string"},
								"field":   {Type: "string"},
							},
						},
					},
				},
			},
			"404": {
				Description: "ISpaceUC not found",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"error": {Type: "string"},
							},
						},
					},
				},
			},
			"500": {
				Description: "Internal server error",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"error": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}, spaceHandler.GetSpace)

	huma.Register(api, huma.Operation{
		OperationID:   "DeleteSpace",
		Method:        http.MethodDelete,
		Path:          "/spaces/{id}",
		Summary:       "delete space",
		Description:   "delete space",
		Tags:          []string{"Spaces"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"204": {
				Description: "ISpaceUC deleted",
				Content:     map[string]*huma.MediaType{},
			},
			"404": {
				Description: "ISpaceUC not found",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"error": {Type: "string"},
							},
						},
					},
				},
			},
			"500": {
				Description: "Internal server error",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"error": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}, spaceHandler.DeleteSpace)

	huma.Register(api, huma.Operation{
		OperationID:   "JoinSpace",
		Method:        http.MethodPost,
		Path:          "/spaces/join",
		Summary:       "join space",
		Description:   "join space",
		Tags:          []string{"Spaces"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "joined ISpaceUC",
				Content:     map[string]*huma.MediaType{},
			},
			"400": {
				Description: "Invalid request",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"message": {Type: "string"},
								"field":   {Type: "string"},
							},
						},
					},
				},
			},
			"500": {
				Description: "Internal server error",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"error": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}, spaceHandler.JoinSpace)
}
