package router

import (
	"github.com/FeisalDy/nogo/config"
	"github.com/FeisalDy/nogo/internal/common"
	"github.com/FeisalDy/nogo/internal/role"
	"github.com/FeisalDy/nogo/internal/user"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *gorm.DB, cfg config.AppConfig) *gin.Engine {
	r := gin.Default()

	uploadService := common.GetUploadService(cfg.BaseURL)
	r.Static("/uploads", uploadService.GetUploadDir())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
				"status":  "healthy",
			})
		})
		userRoutes := v1.Group("/users")
		user.RegisterRoutes(db, userRoutes)

		// Role routes
		roleRoutes := v1.Group("/roles")
		role.RegisterRoutes(db, roleRoutes)

		// Upload routes
		// uploadRoutes := v1.Group("/upload")
		// common.RegisterUploadRoutes(uploadRoutes, cfg.BaseURL)

		// Novel routes
		// novelRoutes := v1.Group("/novels")
		// novel.RegisterRoutes(novelRoutes)

		// Add more domain routes here:
		// authRoutes := v1.Group("/auth")
		// auth.RegisterRoutes(authRoutes)

		// authRoutes := v1.Group("/auth")
		// auth.RegisterRoutes(authRoutes)
	}

	return r
}

// SetupRoutesWithMiddleware sets up routes with additional middleware
func SetupRoutesWithMiddleware(db *gorm.DB, cfg config.AppConfig, middlewares ...gin.HandlerFunc) *gin.Engine {
	r := SetupRoutes(db, cfg)

	// Apply global middleware
	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	return r
}
