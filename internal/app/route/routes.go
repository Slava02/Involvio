package route

import (
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router *fiber.App, pg *database.Postgres) {
	openapiConfig := huma.DefaultConfig("Involvio", "1.0.0")
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
	//setupSpaceRoutes(api, pg)
}
