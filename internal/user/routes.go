package user

import (
	"github.com/FeisalDy/nogo/internal/common/middleware"
	"github.com/FeisalDy/nogo/internal/user/handler"
	"github.com/FeisalDy/nogo/internal/user/repository"
	"github.com/FeisalDy/nogo/internal/user/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(db *gorm.DB, router *gin.RouterGroup) {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router.POST("/login", userHandler.Login)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/:email", userHandler.GetUserByEmail)
	}
}
