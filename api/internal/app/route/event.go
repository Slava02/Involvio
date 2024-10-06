package route

import (
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/handler/rest/v1/event"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/Slava02/Involvio/pkg/valid"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"sync"
)

type EventDeps struct {
	Validator *valid.Validator
}

//nolint:funlen
func setupEventRoutes(api huma.API, pg *database.Postgres, deps *EventDeps) {
	o := sync.Once{}
	eventUseCase := usecase.NewEventUseCase(repository.NewEventRepository(&o, pg))

	eventHandler := event.NewEventHandler(eventUseCase)

	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)
	eventSchema := huma.SchemaFromType(registry, reflect.TypeOf(&entity.Event{}))

	huma.Register(api, huma.Operation{
		OperationID:   "CreateEvent",
		Method:        http.MethodPost,
		Path:          "/events",
		Summary:       "create new event",
		Description:   "Create a new event record.",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "IEventUC created",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: eventSchema,
					},
				},
				Headers: map[string]*huma.Param{
					"Location": {
						Description: "URL of the newly created event",
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
	}, eventHandler.CreateEvent)

	huma.Register(api, huma.Operation{
		OperationID: "GetEvent",
		Method:      http.MethodGet,
		Path:        "/events/{id}",
		Summary:     "event by id",
		Description: "Get event by id.",
		Tags:        []string{"Events"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IEventUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: eventSchema,
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
				Description: "IEventUC not found",
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
	}, eventHandler.GetEvent)

	huma.Register(api, huma.Operation{
		OperationID:   "DeleteEvent",
		Method:        http.MethodDelete,
		Path:          "/events/{id}",
		Summary:       "delete event",
		Description:   "delete event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "deleted IEventUC",
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
			"404": {
				Description: "IEventUC not found",
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
	}, eventHandler.DeleteEvent)

	huma.Register(api, huma.Operation{
		OperationID:   "ReviewEvent",
		Method:        http.MethodPost,
		Path:          "/events/review",
		Summary:       "join event",
		Description:   "join event",
		Tags:          []string{"Events"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "reviewed IEventUC",
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
			"404": {
				Description: "IEventUC not found",
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
	}, eventHandler.ReviewEvent)
}
