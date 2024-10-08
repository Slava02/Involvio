// Package Involvio configures and runs application.
package app

import (
	"fmt"
	"github.com/Slava02/Involvio/config"
	"github.com/Slava02/Involvio/internal/app/route"
	"github.com/Slava02/Involvio/pkg/database"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
)

// Run creates objects via constructors.
// TODO: добавить внедрение зависимостей
type App struct {
	Server *fiber.App
}

func NewApp() *App {
	return &App{
		Server: fiber.New(),
	}
}

func Run(router *fiber.App, cfg *config.Config) {
	// fiber middlewares
	router.Use(logger.New())

	// open telemetry
	router.Use(otelfiber.Middleware())
	router.Use(healthcheck.New())
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	if os.Getenv("ENV_NAME") == "dev" {
		router.Use(pprof.New())
		router.Get("/monitor", monitor.New())
	}

	// Connect to Database
	pg, err := database.New(cfg, database.MaxPoolSize(cfg.DB.PoolMax), database.Isolation(pgx.ReadCommitted))
	if err != nil {
		slog.Error("postgres connection failed", slog.String("error", err.Error()))
		return
	}
	defer pg.Close()

	err = applyMigrations(cfg.DB)
	if err != nil {
		slog.Error("apply migrations failed", slog.String("error", err.Error()))
		return
	}

	// Setup routes
	route.SetupRoutes(router, pg)

	PrintSystemData()
	PrintMemoryInfo()

	// Start server
	slog.Info("Starting server on port: " + cfg.HTTP.Port)
	if err := router.Listen(":" + cfg.HTTP.Port); err != nil {
		slog.Error(fmt.Sprintf("server starting error: %v", err))
	}
}
