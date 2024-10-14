package route

import (
	"github.com/Slava02/Involvio/api/internal/entity"
	"github.com/Slava02/Involvio/api/internal/handler/rest/v1/group"
	"github.com/Slava02/Involvio/api/internal/repository"
	"github.com/Slava02/Involvio/api/internal/usecase"
	"github.com/Slava02/Involvio/api/pkg/database"
	"github.com/Slava02/Involvio/api/pkg/valid"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"sync"
)

type GroupDeps struct {
	Validator *valid.Validator
}

//nolint:funlen
func SetupGroupRoutes(api huma.API, pg *database.Postgres, deps *GroupDeps) {
	o := sync.Once{}
	groupUseCase := usecase.NewGroupUseCase(repository.NewGroupRepository(&o, pg))

	groupHandler := group.NewGroupHandler(groupUseCase)

	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)
	groupSchema := huma.SchemaFromType(registry, reflect.TypeOf(&entity.Group{}))

	huma.Register(api, huma.Operation{
		OperationID:   "CreateGroup",
		Method:        http.MethodPost,
		Path:          "/groups",
		Summary:       "create new group",
		Description:   "Create a new group record.",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "IGroupUC created",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: groupSchema,
					},
				},
				Headers: map[string]*huma.Param{
					"Location": {
						Description: "URL of the newly created group",
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
	}, groupHandler.CreateGroup)

	huma.Register(api, huma.Operation{
		OperationID: "GetGroup",
		Method:      http.MethodGet,
		Path:        "/groups/{name}",
		Summary:     "group by name",
		Description: "Get group by name.",
		Tags:        []string{"Groups"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IGroupUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: groupSchema,
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
				Description: "IGroupUC not found",
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
	}, groupHandler.GetGroup)

	huma.Register(api, huma.Operation{
		OperationID:   "DeleteGroup",
		Method:        http.MethodDelete,
		Path:          "/groups/{name}",
		Summary:       "delete group",
		Description:   "[NOT NEEDED] delete group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"204": {
				Description: "IGroupUC deleted",
				Content:     map[string]*huma.MediaType{},
			},
			"404": {
				Description: "IGroupUC not found",
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
	}, groupHandler.DeleteGroup)

	huma.Register(api, huma.Operation{
		OperationID:   "JoinGroup",
		Method:        http.MethodPost,
		Path:          "/groups/join",
		Summary:       "join group",
		Description:   "join group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "joined IGroupUC",
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
	}, groupHandler.JoinGroup)

	huma.Register(api, huma.Operation{
		OperationID:   "LeaveGroup",
		Method:        http.MethodPost,
		Path:          "/groups/leave",
		Summary:       "leave group",
		Description:   "leave group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "leaved IGroupUC",
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
	}, groupHandler.LeaveGroup)
}
