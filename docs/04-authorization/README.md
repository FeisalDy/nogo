# 04. Authorization (Casbin ABAC/RBAC)

Complete authorization system using Casbin.

## 📄 Documents

### 📘 [Casbin ABAC Guide](CASBIN_ABAC_GUIDE.md)

**Main Guide** - Complete Casbin implementation guide.

- Overview and architecture
- Model configuration
- Usage examples
- Middleware usage
- Managing permissions dynamically
- Integration with existing roles
- Testing and best practices

### 📊 [Casbin Database Schema](CASBIN_DATABASE_SCHEMA.md)

Database structure and schema details.

- Table structure
- How Casbin uses tables
- Migration steps
- Data flow
- Troubleshooting

### 📝 [Casbin Route Examples](CASBIN_ROUTE_EXAMPLES.md)

Practical examples for protecting routes.

- Before/after comparisons
- Different middleware approaches
- Permission seeding scripts
- Testing protected routes

### 🔐 [RBAC Implementation](RBAC_IMPLEMENTATION.md)

Role-Based Access Control implementation details.

- RBAC concepts
- Role hierarchy
- Permission models
- Legacy RBAC info

## Quick Reference

### Protect a Route

```go
router.POST("/users",
    middleware.AuthMiddleware(),
    middleware.CasbinMiddleware("users", "write"),
    handler.CreateUser,
)
```

### Manage Permissions

```go
svc := casbinService.NewCasbinService()

// Add permission to role
svc.AddPermissionForRole("editor", "novels", "write")

// Assign role to user
svc.AssignRoleToUser(userID, "editor")

// Check permission
allowed, _ := svc.Enforce(userID, "novels", "write")
```

### Middleware Options

```go
// Option 1: Specific permission
middleware.CasbinMiddleware("resource", "action")

// Option 2: Any role required
middleware.RequireAnyRole("admin", "editor")

// Option 3: All roles required
middleware.RequireAllRoles("admin", "manager")

// Option 4: Dynamic (auto-detect action)
middleware.DynamicCasbinMiddleware()
```

## Permission Patterns

### Standard Actions

- `read` - View/List operations
- `write` - Create/Edit operations
- `delete` - Delete operations
- `publish` - Custom workflow actions

### Resource Naming

- Use plural nouns: `users`, `novels`, `chapters`
- Match API routes: `/api/v1/users` → `users`
- Keep consistent across app

## Database Tables

### Roles

Store role definitions:

```sql
INSERT INTO roles (name, description)
VALUES ('admin', 'Administrator');
```

### User_Roles

Link users to roles:

```sql
INSERT INTO user_roles (user_id, role_id)
VALUES (1, 1);
```

### Casbin_Rule

Managed by Casbin (don't edit manually):

```sql
-- Policy: admin can write users
ptype='p', v0='admin', v1='users', v2='write'

-- Assignment: user 1 has admin role
ptype='g', v0='user:1', v1='admin'
```

## Common Tasks

### Seed Initial Permissions

```go
svc := casbinService.NewCasbinService()

// Admin
svc.AddPermissionsForRole("admin", [][]string{
    {"users", "read"},
    {"users", "write"},
    {"users", "delete"},
    {"roles", "write"},
})

// Editor
svc.AddPermissionsForRole("editor", [][]string{
    {"novels", "read"},
    {"novels", "write"},
})
```

### Check User Permissions

```go
// Get user roles
roles, _ := svc.GetRolesForUser(userID)

// Get role permissions
perms, _ := svc.GetPermissionsForRole("admin")

// Check specific permission
allowed, _ := svc.Enforce(userID, "users", "delete")
```

### Update Role Permissions

```go
// Add permission
svc.AddPermissionForRole("editor", "chapters", "write")

// Remove permission
svc.RemovePermissionForRole("editor", "chapters", "delete")

// Remove all permissions
svc.RemoveAllPermissionsForRole("editor")
```

[← Back to Main Documentation](../README.md)
