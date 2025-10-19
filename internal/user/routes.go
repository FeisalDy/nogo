package user

import (
	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/middleware"
	roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
	"github.com/FeisalDy/nogo/internal/user/handler"
	"github.com/FeisalDy/nogo/internal/user/repository"
	"github.com/FeisalDy/nogo/internal/user/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(db *gorm.DB, router *gin.RouterGroup) {
	// Initialize user domain dependencies
	userRepository := repository.NewUserRepository(db)
	roleRepository := roleRepo.NewRoleRepository(db)
	casbinSvc := casbinService.NewCasbinService(db)
	userService := service.NewUserService(userRepository, roleRepository, casbinSvc)
	userHandler := handler.NewUserHandler(userService)

	// Public routes (no authentication required)
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.GET("/:email", userHandler.GetUserByEmail)

	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", userHandler.GetMe)
		// Add more protected routes here:
		// protected.PUT("/:id", userHandler.UpdateUser)
		// protected.DELETE("/:id", userHandler.DeleteUser)
		// protected.POST("/:id/avatar", userHandler.UploadAvatar)
	}
}
