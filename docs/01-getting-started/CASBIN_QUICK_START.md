# Casbin ABAC Quick Start

## ⚠️ IMPORTANT: Empty Permissions Array?

If your `/me` endpoint returns empty permissions like this:

```json
{
  "permissions": []
}
```

**You need to seed permissions first!** See: **[📄 SEED_PERMISSIONS.md](./SEED_PERMISSIONS.md)** for detailed guide.

### Quick Fix (30 seconds):

```bash
# Run the seed script
cd /home/feisal/project/shilan/nogo
go run scripts/seed_casbin.go

# Restart your server
# Test again: GET /api/v1/users/me
```

---

## Summary

I've implemented Casbin for dynamic ABAC (Attribute-Based Access Control) in your application. Here's what was added:

## Files Created

### Core Implementation

1. `/internal/common/casbin/casbin_service.go` - Casbin service with policy management
2. `/internal/common/middleware/casbin.go` - Middleware for route protection
3. `/internal/common/model/permission.go` - Permission models
4. `/internal/common/repository/permission_repository.go` - Permission repository

### Configuration

5. `/config/casbin/model.conf` - Casbin RBAC model configuration
6. `/config/casbin/policy.csv` - Policy storage (managed by GORM adapter)

### Documentation

7. `/docs/CASBIN_ABAC_GUIDE.md` - Complete guide
8. `/docs/CASBIN_ROUTE_EXAMPLES.md` - Route examples

### Updated Files

9. `/cmd/server/main.go` - Added Casbin initialization

## Quick Start (5 Minutes)

### Step 1: Build and Run

```bash
cd /home/feisal/project/shilan/nogo
go mod tidy
go build -o server ./cmd/server
./server
```

### Step 2: Seed Permissions (Choose one method)

#### Method A: Using Casbin Service Directly

Create a test file `test_casbin.go`:

```go
package main

import (
	"log"
	"path/filepath"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/config"
	"github.com/FeisalDy/nogo/internal/database"
)

func main() {
	cfg := config.LoadConfig()
	database.Init(cfg.DB)

	modelPath := filepath.Join("config", "casbin", "model.conf")
	casbinService.InitCasbin(modelPath)

	svc := casbinService.NewCasbinService()

	// Add permissions for admin role
	svc.AddPermissionForRole("admin", "users", "read")
	svc.AddPermissionForRole("admin", "users", "write")
	svc.AddPermissionForRole("admin", "users", "delete")
	svc.AddPermissionForRole("admin", "roles", "read")
	svc.AddPermissionForRole("admin", "roles", "write")
	svc.AddPermissionForRole("admin", "roles", "delete")

	// Add permissions for editor role
	svc.AddPermissionForRole("editor", "novels", "read")
	svc.AddPermissionForRole("editor", "novels", "write")

	// Assign admin role to user ID 1
	svc.AssignRoleToUser(1, "admin")

	log.Println("Permissions seeded!")
}
```

Run:

```bash
go run test_casbin.go
```

#### Method B: Via API (after creating endpoints)

### Step 3: Protect Your Routes

Update your routes file (e.g., `internal/role/routes.go`):

```go
import "github.com/FeisalDy/nogo/internal/common/middleware"

func RegisterRoutes(router *gin.RouterGroup) {
	// ... handler init ...

	roleRoutes := router.Group("/roles")
	roleRoutes.Use(middleware.AuthMiddleware())
	{
		// Anyone authenticated can read
		roleRoutes.GET("", roleHandler.GetAllRoles)

		// Only users with "roles:write" permission can create
		roleRoutes.POST("",
			middleware.CasbinMiddleware("roles", "write"),
			roleHandler.CreateRole,
		)

		// Only users with "roles:delete" permission can delete
		roleRoutes.DELETE("/:id",
			middleware.CasbinMiddleware("roles", "delete"),
			roleHandler.DeleteRole,
		)
	}
}
```

### Step 4: Test

```bash
# 1. Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"yourpassword"}'

# 2. Try to create a role (if user has permission)
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"moderator","description":"Content moderator"}'
```

## Available Middleware Functions

```go
// 1. Check specific permission
middleware.CasbinMiddleware("resource", "action")

// 2. Require any role
middleware.RequireAnyRole("admin", "moderator")

// 3. Require all roles
middleware.RequireAllRoles("admin", "billing")

// 4. Dynamic permission (auto-detect from route)
middleware.DynamicCasbinMiddleware()
```

## Common Operations

### Add Permission to Role

```go
svc := casbinService.NewCasbinService()
err := svc.AddPermissionForRole("editor", "novels", "publish")
```

### Assign Role to User

```go
err := svc.AssignRoleToUser(userID, "editor")
```

### Check Permission

```go
allowed, err := svc.Enforce(userID, "novels", "write")
if allowed {
    // User has permission
}
```

### Get User Roles

```go
roles, err := svc.GetRolesForUser(userID)
// Returns: ["editor", "reader"]
```

### Get Role Permissions

```go
permissions, err := svc.GetPermissionsForRole("editor")
// Returns: [["editor", "novels", "write"], ["editor", "novels", "read"]]
```

## Integration with Your Existing Role System

When you assign a role in your database, also assign it in Casbin:

```go
func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
	// Your existing code to assign in database
	// ...

	// Get role name
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	// Assign in Casbin
	casbinSvc := casbinService.NewCasbinService()
	return casbinSvc.AssignRoleToUser(userID, role.Name)
}
```

## Resource Naming Convention

Use these standard names:

- `users` - User management
- `roles` - Role management
- `novels` - Novel resources
- `chapters` - Chapter resources
- `role_assignments` - Role assignment operations

## Action Types

- `read` - View/List operations (GET)
- `write` - Create operations (POST)
- `update` - Update operations (PUT/PATCH)
- `delete` - Delete operations (DELETE)

## Troubleshooting

### Casbin not initializing

```bash
# Check if model file exists
ls -la config/casbin/model.conf

# Check database connection
# Casbin creates 'casbin_rule' table automatically
```

### Permission denied

```bash
# Check if user has role
svc.GetRolesForUser(userID)

# Check if role has permission
svc.GetPermissionsForRole("roleName")

# Reload policies
svc.ReloadPolicies()
```

### Clear all policies (Dev only)

```go
svc.ClearAllPolicies()
```

## Next Steps

1. ✅ Create permission seeding script
2. ✅ Protect your existing routes with Casbin middleware
3. Create admin API endpoints for permission management
4. Build UI for role/permission management
5. Add permission checking in business logic
6. Implement audit logging

## Support

- Full guide: `/docs/CASBIN_ABAC_GUIDE.md`
- Route examples: `/docs/CASBIN_ROUTE_EXAMPLES.md`
- Casbin docs: https://casbin.org/docs/overview
- GORM adapter: https://github.com/casbin/gorm-adapter

## Key Features

✅ **Dynamic Roles** - Add/modify roles at runtime
✅ **Database Storage** - Policies stored in PostgreSQL
✅ **No Restart Required** - Changes take effect immediately
✅ **Flexible** - Support RBAC, ABAC, and hybrid models
✅ **Performant** - In-memory caching with database persistence
✅ **Scalable** - Works with your existing role system
