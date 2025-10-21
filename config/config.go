package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type AppConfig struct {
	Environment    string        // development, staging, production
	Port           string        // server port
	BaseURL        string        // base URL for the application
	Timezone       string        // timezone (e.g., "UTC", "Asia/Jakarta")
	LogLevel       string        // log level (debug, info, warn, error)
	Debug          bool          // debug mode
	RequestTimeout time.Duration // request timeout
}

// Config holds all configuration
type Config struct {
	DB  DBConfig
	App AppConfig
}

// LoadConfig loads all application configuration from environment variables
func LoadConfig() Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		App: LoadAppConfig(),
		DB:  LoadDBConfig(),
	}
}

// LoadAppConfig loads the application configuration
func LoadAppConfig() AppConfig {
	timeout, err := strconv.Atoi(getEnv("REQUEST_TIMEOUT_SECONDS", "30"))
	if err != nil {
		timeout = 30
	}

	debug, err := strconv.ParseBool(getEnv("DEBUG", "false"))
	if err != nil {
		debug = false
	}

	return AppConfig{
		Environment:    getEnv("ENVIRONMENT", "development"),
		Port:           getEnv("PORT", "8080"),
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
		Timezone:       getEnv("TIMEZONE", "UTC"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		Debug:          debug,
		RequestTimeout: time.Duration(timeout) * time.Second,
	}
}

// LoadDBConfig loads the database configuration from environment variables
func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "postgres"),
	}
}

// InitializeApp initializes global application settings
func InitializeApp(config AppConfig) error {
	// Set timezone
	if err := setTimezone(config.Timezone); err != nil {
		log.Printf("Warning: Failed to set timezone to %s: %v", config.Timezone, err)
	} else {
		log.Printf("Timezone set to: %s", config.Timezone)
	}

	// Set Gin mode based on environment
	setGinMode(config.Environment)

	// Log configuration info
	log.Printf("Application initialized:")
	log.Printf("  Environment: %s", config.Environment)
	log.Printf("  Port: %s", config.Port)
	log.Printf("  BaseURL: %s", config.BaseURL)
	log.Printf("  Debug: %t", config.Debug)
	log.Printf("  LogLevel: %s", config.LogLevel)

	return nil
}

// setTimezone sets the application timezone
func setTimezone(timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	// Set the local timezone for the application
	time.Local = loc

	// Also set the TZ environment variable
	os.Setenv("TZ", timezone)

	return nil
}

// setGinMode sets the Gin framework mode based on environment
func setGinMode(environment string) {
	switch environment {
	case "production":
		os.Setenv("GIN_MODE", "release")
		log.Println("Gin mode set to: release")
	case "staging":
		os.Setenv("GIN_MODE", "test")
		log.Println("Gin mode set to: test")
	default:
		os.Setenv("GIN_MODE", "debug")
		log.Println("Gin mode set to: debug")
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
