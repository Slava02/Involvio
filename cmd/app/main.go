package main

import (
	"log"

	"github.com/Slava02/Involvio/config"
	"github.com/Slava02/Involvio/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
