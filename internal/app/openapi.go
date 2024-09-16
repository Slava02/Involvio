package app

import (
	"github.com/Slava02/Involvio/internal/entity"
	v1 "github.com/Slava02/Involvio/internal/handler/rest/v1"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/internal/usecase/repository"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"reflect"
	"sync"
)

func setupRoutes(router *fiber.App, pg *database.Postgres) {
	openapiConfig := huma.DefaultConfig("Clean Architecture Template", "1.0.0")
	openapiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"auth": {
			Type:         "rest",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	openapiConfig.Security = []map[string][]string{
		{"auth": {""}},
	}

	api := humafiber.New(router, openapiConfig)

	setupUserRoutes(api, pg)
}

//nolint:funlen
func setupUserRoutes(api huma.API, pg *database.Postgres) {
	// Initialize use cases
	o := sync.Once{}
	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(&o, pg.Pool))

	// Initialize handlers
	userHandler := v1.NewUserHandler(userUseCase)

	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)

	createUserSchema := huma.SchemaFromType(registry, reflect.TypeOf(&entity.User{}))

	huma.Register(api, huma.Operation{
		OperationID:   "Create user",
		Method:        http.MethodPost,
		Path:          "/user",
		Summary:       "create new user",
		Description:   "Create a new user record.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "IUserUC created",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: createUserSchema,
					},
				},
				Headers: map[string]*huma.Param{
					"Location": {
						Description: "URL of the newly created user",
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
		},
	}, userHandler.CreateUser)

	userByIDSchema := huma.SchemaFromType(registry, reflect.TypeOf(&v1.UserResponse{}))

	huma.Register(api, huma.Operation{
		OperationID: "Get user by id",
		Method:      http.MethodGet,
		Path:        "/user/{id}",
		Summary:     "user by id",
		Description: "Get a user by id.",
		Tags:        []string{"Users"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IUserUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: userByIDSchema,
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
				Description: "IUserUC not found",
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
	}, userHandler.FindUserByID)

	huma.Register(api, huma.Operation{
		OperationID:   "Delete user",
		Method:        http.MethodDelete,
		Path:          "/user/{id}",
		Summary:       "delete user",
		Description:   "Delete a user by ID.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusNoContent,
		Responses: map[string]*huma.Response{
			"204": {
				Description: "IUserUC deleted",
				Content:     map[string]*huma.MediaType{},
			},
			"404": {
				Description: "IUserUC not found",
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
	}, userHandler.DeleteUser)

	huma.Register(api, huma.Operation{
		OperationID: "Update user",
		Method:      http.MethodPatch,
		Path:        "/user/{id}",
		Summary:     "update user",
		Description: "Update an existing user by ID.",
		Tags:        []string{"Users"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IUserUC updated",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: &huma.Schema{
							Type: "object",
							Properties: map[string]*huma.Schema{
								"body": {
									Type: "object",
									Properties: map[string]*huma.Schema{
										"name": {Type: "string"},
									},
									Required: []string{"name"},
								},
							},
							Required: []string{"body"},
						},
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
				Description: "IUserUC not found",
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
	}, userHandler.UpdateUser)
}
