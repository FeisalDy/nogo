package router

import (
	"github.com/FeisalDy/nogo/config"
	"github.com/FeisalDy/nogo/internal/application"
	"github.com/FeisalDy/nogo/internal/novel"
	"github.com/FeisalDy/nogo/internal/role"
	"github.com/FeisalDy/nogo/internal/user"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *gorm.DB, cfg config.AppConfig) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
				"status":  "healthy",
			})
		})

		application.RegisterRoutes(db, v1)

		userRoutes := v1.Group("/users")
		user.RegisterRoutes(db, userRoutes)

		roleRoutes := v1.Group("/roles")
		role.RegisterRoutes(db, roleRoutes)

		novelRoutes := v1.Group("/novels")
		novel.RegisterRoutes(db, novelRoutes)
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
