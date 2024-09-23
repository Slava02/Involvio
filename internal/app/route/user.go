package route

import (
	"github.com/Slava02/Involvio/internal/entity"
	"github.com/Slava02/Involvio/internal/handler/rest/v1/user"
	"github.com/Slava02/Involvio/internal/repository"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"sync"
)

//nolint:funlen
func setupUserRoutes(api huma.API, pg *database.Postgres) {
	// Initialize use cases
	o := sync.Once{}
	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(&o, pg))

	// Initialize handlers
	userHandler := user.NewUserHandler(userUseCase)

	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)

	userSchema := huma.SchemaFromType(registry, reflect.TypeOf(&entity.User{}))
	formSchema := huma.SchemaFromType(registry, reflect.TypeOf(&entity.Form{}))
	userWithFormsSchema := huma.SchemaFromType(registry, reflect.TypeOf(&user.UserWithFormsResponse{}))

	huma.Register(api, huma.Operation{
		OperationID:   "CreateUser",
		Method:        http.MethodPost,
		Path:          "/users",
		Summary:       "create new user",
		Description:   "Create a new user record.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
		Responses: map[string]*huma.Response{
			"201": {
				Description: "IUserUC created",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: userSchema,
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
	}, userHandler.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "GetUser",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "user by id",
		Description: "Get a user by id.",
		Tags:        []string{"Users"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IUserUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: userWithFormsSchema,
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
	}, userHandler.GetUserWithForms)

	huma.Register(api, huma.Operation{
		OperationID: "UpdateUser",
		Method:      http.MethodPut,
		Path:        "/users/{id}",
		Summary:     "update user",
		Description: "Update an existing user by ID.",
		Tags:        []string{"Users"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IUserUC updated",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: userSchema,
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
	}, userHandler.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID:   "DeleteUserForm",
		Method:        http.MethodDelete,
		Path:          "/users/{userId}/{spaceId}",
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
	}, userHandler.DeleteUser)

	huma.Register(api, huma.Operation{
		OperationID:   "GetUserForm",
		Method:        http.MethodGet,
		Path:          "/users/{userId}/{spaceId}",
		Summary:       "get user form in space",
		Description:   "returns user form in provided space",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusNoContent,
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IUserUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: formSchema,
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
	}, userHandler.GetForm)

	huma.Register(api, huma.Operation{
		OperationID:   "UpdateUserForm",
		Method:        http.MethodPut,
		Path:          "/users/{userId}/{spaceId}",
		Summary:       "update user form in space",
		Description:   "returns updated user form in provided space",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusNoContent,
		Responses: map[string]*huma.Response{
			"200": {
				Description: "IUserUC response",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: formSchema,
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
	}, userHandler.UpdateForm)
}
