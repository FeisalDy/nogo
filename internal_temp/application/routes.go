package application

import (
	"github.com/FeisalDy/nogo/internal/application/handler"
	"github.com/FeisalDy/nogo/internal/application/service"
	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/middleware"
	roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.RouterGroup) {
	userRepository := userRepo.NewUserRepository(db)
	roleRepository := roleRepo.NewRoleRepository(db)
	casbinSvc := casbinService.NewCasbinService(db)

	userRoleService := service.NewUserRoleService(userRepository, roleRepository, casbinSvc)
	authService := service.NewAuthService(userRepository, roleRepository, casbinSvc)
	userProfileService := service.NewUserProfileService(userRepository, roleRepository, casbinSvc)

	userRoleHandler := handler.NewUserRoleHandler(userRoleService)
	authHandler := handler.NewAuthHandler(authService)
	userProfileHandler := handler.NewUserProfileHandler(userProfileService)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
	}

	profileRoutes := router.Group("/profile")
	profileRoutes.Use(middleware.AuthMiddleware())
	{
		profileRoutes.GET("/me", userProfileHandler.GetMe)
	}

	userRoleRoutes := router.Group("/user-roles")
	userRoleRoutes.Use(middleware.AuthMiddleware())
	{
		userRoleRoutes.POST("/assign", userRoleHandler.AssignRoleToUser)
		userRoleRoutes.POST("/remove", userRoleHandler.RemoveRoleFromUser)

		userRoleRoutes.GET("/users/:user_id/roles", userRoleHandler.GetUserRoles)
	}
}
