package main

import (
	"log"

	"boiler/config"
	"boiler/internal/database"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database and run migrations
	log.Println("Starting migration system...")
	database.Init(cfg.DB)
	log.Println("Migration system completed successfully!")
}
