package user

import (
	"boiler/internal/user/handler"
	"boiler/internal/user/repository"
	"boiler/internal/user/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router *gin.RouterGroup) {
	// Initialize user domain dependencies
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	// Register user routes
	router.POST("/", userHandler.CreateUser)
	router.GET("/:id", userHandler.GetUser)

	// Add more user routes here as needed:
	// router.PUT("/:id", userHandler.UpdateUser)
	// router.DELETE("/:id", userHandler.DeleteUser)
	// router.GET("/", userHandler.GetAllUsers)
	// router.POST("/:id/avatar", userHandler.UploadAvatar)
}
