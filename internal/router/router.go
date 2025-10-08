package router

import (
	"boiler/config"
	"boiler/internal/common"
	"boiler/internal/novel"
	"boiler/internal/user"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all application routes
func SetupRoutes(cfg config.AppConfig) *gin.Engine {
	r := gin.Default()

	// Serve uploaded files statically
	uploadService := common.GetUploadService(cfg.BaseURL)
	r.Static("/uploads", uploadService.GetUploadDir())

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "healthy",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// User routes
		userRoutes := v1.Group("/users")
		user.RegisterRoutes(userRoutes)

		// Upload routes
		uploadRoutes := v1.Group("/upload")
		common.RegisterUploadRoutes(uploadRoutes, cfg.BaseURL)

		// Novel routes
		novelRoutes := v1.Group("/novels")
		novel.RegisterRoutes(novelRoutes)

		// Add more domain routes here:
		// authRoutes := v1.Group("/auth")
		// auth.RegisterRoutes(authRoutes)

		// authRoutes := v1.Group("/auth")
		// auth.RegisterRoutes(authRoutes)
	}

	return r
}

// SetupRoutesWithMiddleware sets up routes with additional middleware
func SetupRoutesWithMiddleware(cfg config.AppConfig, middlewares ...gin.HandlerFunc) *gin.Engine {
	r := SetupRoutes(cfg)

	// Apply global middleware
	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	return r
}
