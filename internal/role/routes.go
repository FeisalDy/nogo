package role

import (
	"github.com/FeisalDy/nogo/internal/common/middleware"
	"github.com/FeisalDy/nogo/internal/role/handler"
	"github.com/FeisalDy/nogo/internal/role/repository"
	"github.com/FeisalDy/nogo/internal/role/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.RouterGroup) {
	roleRepository := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepository)
	roleHandler := handler.NewRoleHandler(roleService)

	roleRoutes := router.Group("/")
	roleRoutes.Use(middleware.AuthMiddleware())
	{
		// CRUD operations
		roleRoutes.POST("", roleHandler.CreateRole)       // Create a new role
		roleRoutes.GET("", roleHandler.GetAllRoles)       // Get all roles
		roleRoutes.GET("/:id", roleHandler.GetRole)       // Get a role by ID
		roleRoutes.PUT("/:id", roleHandler.UpdateRole)    // Update a role
		roleRoutes.DELETE("/:id", roleHandler.DeleteRole) // Delete a role
	}
}
