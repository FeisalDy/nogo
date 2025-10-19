# Casbin ABAC Implementation Guide

## Overview

This implementation uses [Casbin](https://casbin.org/) for Attribute-Based Access Control (ABAC) with dynamic role management. Casbin policies are stored in the database using GORM adapter, allowing runtime modifications without service restarts.

## Architecture

### Components

1. **Casbin Enforcer** - Core authorization engine
2. **GORM Adapter** - Stores policies in PostgreSQL database
3. **Casbin Service** - Business logic for policy management
4. **Casbin Middleware** - HTTP middleware for route protection
5. **Permission Repository** - Database operations for permissions

### Model Configuration

Located at: `config/casbin/model.conf`

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

**Explanation:**
- `sub` (subject) - User identifier (e.g., "user:123")
- `obj` (object) - Resource name (e.g., "users", "novels")
- `act` (action) - Permission type (e.g., "read", "write", "delete")
- `g` - Role inheritance/assignment

## Database Schema

### Casbin Rules Table
Auto-created by GORM adapter as `casbin_rule`:

```sql
CREATE TABLE casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),  -- 'p' for policy, 'g' for grouping
    v0 VARCHAR(100),     -- subject/user
    v1 VARCHAR(100),     -- object/resource
    v2 VARCHAR(100),     -- action
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

### Permissions Table
Located in: `internal/common/model/permission.go`

```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Usage Examples

### 1. Initialize Casbin (Already in main.go)

```go
import casbinService "github.com/FeisalDy/nogo/internal/common/casbin"

modelPath := filepath.Join("config", "casbin", "model.conf")
_, err := casbinService.InitCasbin(modelPath)
if err != nil {
    log.Fatalf("Failed to initialize Casbin: %v", err)
}
```

### 2. Using Casbin Service

```go
package handler

import (
    casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
)

func ManagePermissions() {
    svc := casbinService.NewCasbinService()
    
    // Add permission to role
    err := svc.AddPermissionForRole("admin", "users", "write")
    err = svc.AddPermissionForRole("editor", "novels", "write")
    err = svc.AddPermissionForRole("reader", "novels", "read")
    
    // Assign role to user
    err = svc.AssignRoleToUser(userID, "admin")
    
    // Check permission
    allowed, err := svc.Enforce(userID, "users", "write")
    if allowed {
        // User has permission
    }
    
    // Get user roles
    roles, err := svc.GetRolesForUser(userID)
    
    // Get role permissions
    permissions, err := svc.GetPermissionsForRole("admin")
}
```

### 3. Using Middleware

#### Method 1: Static Resource/Action

```go
import "github.com/FeisalDy/nogo/internal/common/middleware"

// Protect specific routes
router.GET("/users",
    middleware.AuthMiddleware(),
    middleware.CasbinMiddleware("users", "read"),
    handler.ListUsers,
)

router.POST("/users",
    middleware.AuthMiddleware(),
    middleware.CasbinMiddleware("users", "write"),
    handler.CreateUser,
)

router.DELETE("/users/:id",
    middleware.AuthMiddleware(),
    middleware.CasbinMiddleware("users", "delete"),
    handler.DeleteUser,
)
```

#### Method 2: Dynamic Resource/Action

```go
// Automatically determine action from HTTP method
router.Use(middleware.AuthMiddleware())
router.Use(middleware.DynamicCasbinMiddleware())

// Or set custom resource/action
router.GET("/novels/:id",
    middleware.SetCasbinParams("novels", "read"),
    middleware.DynamicCasbinMiddleware(),
    handler.GetNovel,
)
```

#### Method 3: Role-Based

```go
// Require specific roles
router.GET("/admin/dashboard",
    middleware.AuthMiddleware(),
    middleware.RequireAnyRole("admin", "super_admin"),
    handler.AdminDashboard,
)

router.POST("/admin/settings",
    middleware.AuthMiddleware(),
    middleware.RequireAllRoles("admin", "settings_manager"),
    handler.UpdateSettings,
)
```

### 4. Example Route Setup

```go
// internal/role/routes.go
package role

import (
    "github.com/FeisalDy/nogo/internal/common/middleware"
    "github.com/FeisalDy/nogo/internal/role/handler"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
    h := handler.NewRoleHandler()
    
    // Public routes (no auth required)
    rg.GET("/roles", h.GetAllRoles)
    
    // Protected routes
    authorized := rg.Group("")
    authorized.Use(middleware.AuthMiddleware())
    {
        // Read permissions
        authorized.GET("/roles/:id",
            middleware.CasbinMiddleware("roles", "read"),
            h.GetRoleByID,
        )
        
        // Write permissions
        authorized.POST("/roles",
            middleware.CasbinMiddleware("roles", "write"),
            h.CreateRole,
        )
        
        authorized.PUT("/roles/:id",
            middleware.CasbinMiddleware("roles", "write"),
            h.UpdateRole,
        )
        
        // Delete permissions
        authorized.DELETE("/roles/:id",
            middleware.CasbinMiddleware("roles", "delete"),
            h.DeleteRole,
        )
        
        // Role assignments (require special permission)
        authorized.POST("/roles/:id/assign",
            middleware.CasbinMiddleware("role_assignments", "write"),
            h.AssignRoleToUser,
        )
    }
}
```

## Managing Permissions Dynamically

### Create Initial Permissions and Roles

```go
package main

import casbinService "github.com/FeisalDy/nogo/internal/common/casbin"

func SeedPermissions() {
    svc := casbinService.NewCasbinService()
    
    // Define permissions for each role
    adminPermissions := [][]string{
        {"users", "read"},
        {"users", "write"},
        {"users", "delete"},
        {"roles", "read"},
        {"roles", "write"},
        {"roles", "delete"},
        {"novels", "read"},
        {"novels", "write"},
        {"novels", "delete"},
    }
    
    editorPermissions := [][]string{
        {"novels", "read"},
        {"novels", "write"},
        {"chapters", "read"},
        {"chapters", "write"},
    }
    
    readerPermissions := [][]string{
        {"novels", "read"},
        {"chapters", "read"},
    }
    
    // Add permissions to roles
    svc.AddPermissionsForRole("admin", adminPermissions)
    svc.AddPermissionsForRole("editor", editorPermissions)
    svc.AddPermissionsForRole("reader", readerPermissions)
}
```

### API Endpoints for Permission Management

You should create REST API endpoints to manage permissions:

```go
// POST /api/v1/permissions - Create permission
// GET /api/v1/permissions - List all permissions
// DELETE /api/v1/permissions/:id - Delete permission

// POST /api/v1/roles/:id/permissions - Assign permission to role
// DELETE /api/v1/roles/:id/permissions/:permissionId - Remove permission from role
// GET /api/v1/roles/:id/permissions - Get role permissions

// POST /api/v1/users/:id/roles - Assign role to user
// DELETE /api/v1/users/:id/roles/:roleId - Remove role from user
// GET /api/v1/users/:id/roles - Get user roles
```

## Common Permission Patterns

### Resource-Based Permissions

```go
// Basic CRUD
"resource:read"
"resource:write"
"resource:update"
"resource:delete"

// Specific resources
"users:read"
"users:write"
"novels:read"
"novels:publish"
"chapters:read"
"chapters:write"
```

### Action Types

- `read` / `get` / `list` - View operations
- `write` / `create` / `post` - Create operations
- `update` / `put` / `patch` - Update operations
- `delete` / `remove` - Delete operations
- `publish` / `approve` - Workflow actions
- `*` - All actions (wildcard, use carefully)

## Integration with Existing Role System

Your existing role system (in database) should sync with Casbin:

```go
// When creating a role in database
func (s *RoleService) CreateRole(req dto.CreateRoleDTO) (*model.Role, error) {
    // ... existing code to create role in database ...
    
    // Sync with Casbin (optional, if you want default permissions)
    casbinSvc := casbinService.NewCasbinService()
    // Add default permissions for the new role if needed
    
    return role, nil
}

// When assigning role to user
func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
    // ... existing code ...
    
    // Get role name
    role, err := s.roleRepo.GetByID(roleID)
    if err != nil {
        return err
    }
    
    // Assign in Casbin
    casbinSvc := casbinService.NewCasbinService()
    err = casbinSvc.AssignRoleToUser(userID, role.Name)
    
    return err
}
```

## Testing

### Test Permission Check

```go
func TestUserPermission(t *testing.T) {
    svc := casbinService.NewCasbinService()
    
    // Setup: Create role and assign permissions
    svc.AddPermissionForRole("editor", "novels", "write")
    svc.AssignRoleToUser(123, "editor")
    
    // Test: Check permission
    allowed, err := svc.Enforce(123, "novels", "write")
    assert.NoError(t, err)
    assert.True(t, allowed)
    
    // Test: Check denied permission
    allowed, err = svc.Enforce(123, "users", "delete")
    assert.NoError(t, err)
    assert.False(t, allowed)
}
```

## Best Practices

1. **Use Meaningful Resource Names**
   - Use plural nouns: `users`, `novels`, `chapters`
   - Keep names consistent with your API routes

2. **Standard Action Names**
   - Stick to: `read`, `write`, `update`, `delete`
   - Add custom actions only when needed: `publish`, `approve`

3. **Role Naming**
   - Use lowercase: `admin`, `editor`, `reader`
   - Be descriptive: `content_moderator`, `billing_admin`

4. **Permission Granularity**
   - Start with coarse permissions, refine as needed
   - Don't create too many fine-grained permissions initially

5. **Sync Database and Casbin**
   - When roles change in DB, update Casbin
   - When users get roles in DB, update Casbin
   - Consider event-driven sync or periodic sync

6. **Performance**
   - Casbin caches policies in memory
   - Database queries only on policy changes
   - Call `LoadPolicy()` if policies change externally

## Troubleshooting

### Permissions Not Working

1. Check if Casbin is initialized
2. Verify policies are loaded: `svc.GetAllRoles()`
3. Check user subject format: `user:123`
4. Verify role assignments: `svc.GetRolesForUser(userID)`

### Role Changes Not Reflecting

```go
// Reload policies from database
svc := casbinService.NewCasbinService()
err := svc.ReloadPolicies()
```

### Clear All Policies (Development Only)

```go
svc := casbinService.NewCasbinService()
err := svc.ClearAllPolicies()
```

## Next Steps

1. Create permission management API endpoints
2. Build admin UI for role/permission management
3. Add permission seeding in migrations
4. Implement audit logging for permission changes
5. Add permission checking in business logic (not just middleware)

## Resources

- [Casbin Documentation](https://casbin.org/docs/overview)
- [RBAC Model](https://casbin.org/docs/rbac)
- [GORM Adapter](https://github.com/casbin/gorm-adapter)
