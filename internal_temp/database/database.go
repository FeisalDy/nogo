package database

import (
	"fmt"
	"log"

	"github.com/FeisalDy/nogo/config"
	"github.com/FeisalDy/nogo/internal/database/migrations"
	"github.com/FeisalDy/nogo/internal/database/seeds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the database connection
var DB *gorm.DB

// Init initializes the database connection and runs migrations
func Init(cfg config.DBConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Database connected")

	// Run migrations
	log.Println("Running database migrations...")
	if err := migrations.RunMigrations(DB); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")
}

// RunSeeds runs all database seeders
// This should be called after Casbin is initialized
func RunSeeds() {
	log.Println("Running database seeders...")
	if err := seeds.RunAllSeeds(DB); err != nil {
		log.Printf("Warning: Failed to run seeders: %v", err)
	}
}
