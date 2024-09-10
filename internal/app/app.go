// Package app configures and runs application.
package app

import (
	"fmt"
	postgres2 "github.com/Slava02/Involvio/internal/repo/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/Slava02/Involvio/config"
	v1 "github.com/Slava02/Involvio/internal/controller/http/v1"
	"github.com/Slava02/Involvio/internal/usecase"
	"github.com/Slava02/Involvio/pkg/httpserver"
	"github.com/Slava02/Involvio/pkg/logger"
	"github.com/Slava02/Involvio/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// INIT LOGGER
	l := logger.New(cfg.Log.Level)

	//  INIT REPOS
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	//  INIT DEPENDENCIES

	//  INIT USECASES
	useCase := usecase.New(&usecase.Deps{
		postgres2.New(pg),
	})

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, useCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
