package main

import (
	"log"

	"boiler/config"
	"boiler/internal/database"
	"boiler/internal/router"
)

func main() {
	cfg := config.LoadConfig()
	if err := config.InitializeApp(cfg.App); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	database.Init(cfg.DB)
	r := router.SetupRoutes(cfg.App)

	serverAddr := ":" + cfg.App.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
