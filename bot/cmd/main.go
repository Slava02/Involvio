package main

import (
	"fmt"
	"github.com/Slava02/Involvio/bot/config"
	"github.com/Slava02/Involvio/bot/internal/app"
	"github.com/Slava02/Involvio/bot/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	err = Run(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func Run(cfg *config.Config) error {
	// Initialize logger
	logger.SetupLogger(cfg)

	// Run the application
	bot := app.New(cfg)

	go func() {
		if err := bot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	return nil
}
