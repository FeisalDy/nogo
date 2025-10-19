package main

import (
	"log"
	"path/filepath"

	"github.com/FeisalDy/nogo/config"
	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/database"
	"github.com/FeisalDy/nogo/internal/router"
)

func main() {
	cfg := config.LoadConfig()
	if err := config.InitializeApp(cfg.App); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	database.Init(cfg.DB)

	modelPath := filepath.Join("config", "casbin", "model.conf")
	_, err := casbinService.InitCasbin(database.DB, modelPath)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin: %v", err)
	}
	log.Println("Casbin initialized successfully")

	// Auto-seed Casbin permissions (runs after Casbin is initialized)
	database.SeedCasbin()

	r := router.SetupRoutes(database.DB, cfg.App)

	serverAddr := ":" + cfg.App.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
