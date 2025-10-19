# Example: Updating Routes with Casbin Middleware

## Before (Auth Only)

```go
package role

import (
	"github.com/FeisalDy/nogo/internal/common/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	roleRoutes := router.Group("/roles")
	roleRoutes.Use(middleware.AuthMiddleware())
	{
		roleRoutes.POST("", roleHandler.CreateRole)
		roleRoutes.GET("", roleHandler.GetAllRoles)
		roleRoutes.GET("/:id", roleHandler.GetRole)
		roleRoutes.PUT("/:id", roleHandler.UpdateRole)
		roleRoutes.DELETE("/:id", roleHandler.DeleteRole)
	}
}
```

## After (Auth + Casbin)

```go
package role

import (
	"github.com/FeisalDy/nogo/internal/common/middleware"
	"github.com/FeisalDy/nogo/internal/database"
	"github.com/FeisalDy/nogo/internal/role/handler"
	"github.com/FeisalDy/nogo/internal/role/repository"
	"github.com/FeisalDy/nogo/internal/role/service"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	roleRepository := repository.NewRoleRepository(database.DB)
	userRepository := userRepo.NewUserRepository(database.DB)
	roleService := service.NewRoleService(roleRepository, userRepository)
	roleHandler := handler.NewRoleHandler(roleService)

	roleRoutes := router.Group("/roles")
	roleRoutes.Use(middleware.AuthMiddleware()) // Must be authenticated
	{
		// Public read access (authenticated users can list/view roles)
		roleRoutes.GET("", roleHandler.GetAllRoles)
		roleRoutes.GET("/:id", roleHandler.GetRole)
		
		// Write operations - require "roles:write" permission
		roleRoutes.POST("",
			middleware.CasbinMiddleware("roles", "write"),
			roleHandler.CreateRole,
		)
		
		roleRoutes.PUT("/:id",
			middleware.CasbinMiddleware("roles", "write"),
			roleHandler.UpdateRole,
		)
		
		// Delete operations - require "roles:delete" permission
		roleRoutes.DELETE("/:id",
			middleware.CasbinMiddleware("roles", "delete"),
			roleHandler.DeleteRole,
		)
		
		// Role assignment operations
		roleRoutes.POST("/:id/users/:user_id",
			middleware.CasbinMiddleware("role_assignments", "write"),
			roleHandler.AssignRoleToUser,
		)
		
		roleRoutes.DELETE("/:id/users/:user_id",
			middleware.CasbinMiddleware("role_assignments", "delete"),
			roleHandler.RemoveRoleFromUser,
		)
	}
}
```

## Alternative: Role-Based Middleware

```go
func RegisterRoutes(router *gin.RouterGroup) {
	// ...handler initialization...

	roleRoutes := router.Group("/roles")
	roleRoutes.Use(middleware.AuthMiddleware())
	{
		// Anyone can read
		roleRoutes.GET("", roleHandler.GetAllRoles)
		roleRoutes.GET("/:id", roleHandler.GetRole)
		
		// Only admins can modify
		adminOnly := roleRoutes.Group("")
		adminOnly.Use(middleware.RequireAnyRole("admin", "super_admin"))
		{
			adminOnly.POST("", roleHandler.CreateRole)
			adminOnly.PUT("/:id", roleHandler.UpdateRole)
			adminOnly.DELETE("/:id", roleHandler.DeleteRole)
			adminOnly.POST("/:id/users/:user_id", roleHandler.AssignRoleToUser)
			adminOnly.DELETE("/:id/users/:user_id", roleHandler.RemoveRoleFromUser)
		}
	}
}
```

## Setup Permissions (One-time setup)

Create a script or migration to set up initial permissions:

```go
// scripts/seed_permissions.go
package main

import (
	"log"
	
	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/config"
	"github.com/FeisalDy/nogo/internal/database"
)

func main() {
	// Initialize database
	cfg := config.LoadConfig()
	database.Init(cfg.DB)
	
	// Initialize Casbin
	modelPath := "config/casbin/model.conf"
	casbinService.InitCasbin(modelPath)
	
	svc := casbinService.NewCasbinService()
	
	// Define permissions for admin role
	adminPerms := [][]string{
		{"roles", "read"},
		{"roles", "write"},
		{"roles", "delete"},
		{"role_assignments", "write"},
		{"role_assignments", "delete"},
		{"users", "read"},
		{"users", "write"},
		{"users", "delete"},
		{"novels", "read"},
		{"novels", "write"},
		{"novels", "delete"},
		{"chapters", "read"},
		{"chapters", "write"},
		{"chapters", "delete"},
	}
	
	// Add permissions to admin role
	err := svc.AddPermissionsForRole("admin", adminPerms)
	if err != nil {
		log.Fatalf("Failed to add admin permissions: %v", err)
	}
	
	// Define permissions for editor role
	editorPerms := [][]string{
		{"roles", "read"},
		{"users", "read"},
		{"novels", "read"},
		{"novels", "write"},
		{"chapters", "read"},
		{"chapters", "write"},
	}
	
	err = svc.AddPermissionsForRole("editor", editorPerms)
	if err != nil {
		log.Fatalf("Failed to add editor permissions: %v", err)
	}
	
	// Define permissions for reader role
	readerPerms := [][]string{
		{"novels", "read"},
		{"chapters", "read"},
	}
	
	err = svc.AddPermissionsForRole("reader", readerPerms)
	if err != nil {
		log.Fatalf("Failed to add reader permissions: %v", err)
	}
	
	log.Println("Permissions seeded successfully!")
	
	// Optionally assign admin role to first user
	// err = svc.AssignRoleToUser(1, "admin")
}
```

Run the script:
```bash
go run scripts/seed_permissions.go
```

## Testing the Protected Routes

```bash
# Login to get token
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}'

# Use token to create role (requires admin permission)
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"name":"moderator","description":"Content moderator"}'

# If user doesn't have permission, will get 401 Unauthorized
```
