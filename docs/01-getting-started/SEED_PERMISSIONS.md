# 🌱 Seeding Permissions to Casbin

## ❌ Problem: Empty Permissions Array

You're seeing this response:

```json
{
  "permissions": []
}
```

**Why?** Because you haven't added any permissions to Casbin yet!

## ✅ Solution: Seed Permissions

You have **3 methods** to seed permissions. Choose the easiest one for you.

---

## 📝 Method 1: Direct SQL (Fastest!)

### Step 1: Connect to your database

```bash
# PostgreSQL
psql -U your_username -d your_database_name
```

### Step 2: Insert permissions directly

```sql
-- Add permissions for 'admin' role
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'admin', 'users', 'read'),
('p', 'admin', 'users', 'write'),
('p', 'admin', 'users', 'delete'),
('p', 'admin', 'novels', 'read'),
('p', 'admin', 'novels', 'write'),
('p', 'admin', 'novels', 'delete'),
('p', 'admin', 'chapters', 'read'),
('p', 'admin', 'chapters', 'write'),
('p', 'admin', 'chapters', 'delete'),
('p', 'admin', 'roles', 'read'),
('p', 'admin', 'roles', 'write');

-- Add permissions for 'user' role (limited)
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'user', 'novels', 'read'),
('p', 'user', 'chapters', 'read'),
('p', 'user', 'profile', 'read'),
('p', 'user', 'profile', 'write');

-- Add permissions for 'author' role (can write novels)
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'author', 'novels', 'read'),
('p', 'author', 'novels', 'write'),
('p', 'author', 'chapters', 'read'),
('p', 'author', 'chapters', 'write'),
('p', 'author', 'profile', 'read'),
('p', 'author', 'profile', 'write');
```

### Step 3: Verify

```sql
SELECT * FROM casbin_rule WHERE ptype = 'p';
```

You should see all your permissions!

### Step 4: Restart your server

```bash
# Stop the server (Ctrl+C)
# Start again
./tmp/main
```

### Step 5: Test `/me` endpoint again

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN" | jq
```

Now you should see permissions! 🎉

---

## 🚀 Method 2: Create a Seed Script (Recommended for Production)

### Step 1: Create seed script

Create file: `scripts/seed_casbin.go`

```go
package main

import (
	"fmt"
	"log"
	"path/filepath"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/config"
	"github.com/FeisalDy/nogo/internal/database"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Casbin
	modelPath := filepath.Join("config", "casbin", "model.conf")
	if _, err := casbinService.InitCasbin(database.DB, modelPath); err != nil {
		log.Fatalf("Failed to initialize Casbin: %v", err)
	}

	// Create Casbin service
	casbin := casbinService.NewCasbinService(database.DB)

	fmt.Println("🌱 Seeding Casbin permissions...")

	// Admin permissions (full access)
	adminPerms := []struct {
		resource string
		action   string
	}{
		{"users", "read"},
		{"users", "write"},
		{"users", "delete"},
		{"novels", "read"},
		{"novels", "write"},
		{"novels", "delete"},
		{"chapters", "read"},
		{"chapters", "write"},
		{"chapters", "delete"},
		{"genres", "read"},
		{"genres", "write"},
		{"tags", "read"},
		{"tags", "write"},
		{"roles", "read"},
		{"roles", "write"},
	}

	for _, perm := range adminPerms {
		if err := casbin.AddPermissionForRole("admin", perm.resource, perm.action); err != nil {
			log.Printf("Warning: Failed to add admin permission %s:%s - %v", perm.resource, perm.action, err)
		} else {
			fmt.Printf("✓ Added admin permission: %s:%s\n", perm.resource, perm.action)
		}
	}

	// Author permissions (can create/edit own content)
	authorPerms := []struct {
		resource string
		action   string
	}{
		{"novels", "read"},
		{"novels", "write"},
		{"chapters", "read"},
		{"chapters", "write"},
		{"genres", "read"},
		{"tags", "read"},
		{"profile", "read"},
		{"profile", "write"},
	}

	for _, perm := range authorPerms {
		if err := casbin.AddPermissionForRole("author", perm.resource, perm.action); err != nil {
			log.Printf("Warning: Failed to add author permission %s:%s - %v", perm.resource, perm.action, err)
		} else {
			fmt.Printf("✓ Added author permission: %s:%s\n", perm.resource, perm.action)
		}
	}

	// User permissions (read-only)
	userPerms := []struct {
		resource string
		action   string
	}{
		{"novels", "read"},
		{"chapters", "read"},
		{"genres", "read"},
		{"tags", "read"},
		{"profile", "read"},
		{"profile", "write"},
	}

	for _, perm := range userPerms {
		if err := casbin.AddPermissionForRole("user", perm.resource, perm.action); err != nil {
			log.Printf("Warning: Failed to add user permission %s:%s - %v", perm.resource, perm.action, err)
		} else {
			fmt.Printf("✓ Added user permission: %s:%s\n", perm.resource, perm.action)
		}
	}

	fmt.Println("\n✅ Casbin permissions seeded successfully!")
	fmt.Println("\n📊 Summary:")

	// Get all roles
	allPerms, _ := casbin.GetEnforcer().GetPolicy()
	fmt.Printf("Total permissions: %d\n", len(allPerms))

	adminPermsCount, _ := casbin.GetPermissionsForRole("admin")
	fmt.Printf("Admin permissions: %d\n", len(adminPermsCount))

	authorPermsCount, _ := casbin.GetPermissionsForRole("author")
	fmt.Printf("Author permissions: %d\n", len(authorPermsCount))

	userPermsCount, _ := casbin.GetPermissionsForRole("user")
	fmt.Printf("User permissions: %d\n", len(userPermsCount))
}
```

### Step 2: Run the seed script

```bash
cd /home/feisal/project/shilan/nogo
go run scripts/seed_casbin.go
```

You should see output like:

```
🌱 Seeding Casbin permissions...
✓ Added admin permission: users:read
✓ Added admin permission: users:write
✓ Added admin permission: users:delete
...
✅ Casbin permissions seeded successfully!

📊 Summary:
Total permissions: 27
Admin permissions: 15
Author permissions: 8
User permissions: 6
```

### Step 3: Test again!

---

## 🔧 Method 3: Using API Endpoints (Create Admin Endpoint)

### Step 1: Create admin permission management endpoint

This is optional but useful for production. Create file: `internal/role/handler/permission_handler.go`

```go
package handler

import (
	"net/http"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionHandler struct {
	casbinService *casbinService.CasbinService
}

func NewPermissionHandler(db *gorm.DB) *PermissionHandler {
	return &PermissionHandler{
		casbinService: casbinService.NewCasbinService(db),
	}
}

type AddPermissionRequest struct {
	RoleName string `json:"role_name" binding:"required"`
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

func (h *PermissionHandler) AddPermission(c *gin.Context) {
	var req AddPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeValidation)
		return
	}

	if err := h.casbinService.AddPermissionForRole(req.RoleName, req.Resource, req.Action); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, nil, "Permission added successfully")
}

func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleName := c.Param("role")

	permissions, err := h.casbinService.GetPermissionsForRole(roleName)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Convert to readable format
	type PermissionDTO struct {
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}

	perms := make([]PermissionDTO, 0, len(permissions))
	for _, p := range permissions {
		if len(p) >= 3 {
			perms = append(perms, PermissionDTO{
				Resource: p[1],
				Action:   p[2],
			})
		}
	}

	utils.RespondSuccess(c, http.StatusOK, perms)
}
```

### Step 2: Add routes (in `internal/role/routes.go`)

```go
// Add to RegisterRoutes function
permHandler := handler.NewPermissionHandler(db)

// Admin only routes
admin := router.Group("/")
admin.Use(middleware.AuthMiddleware())
admin.Use(middleware.RequireAnyRole("admin"))
{
	admin.POST("/permissions", permHandler.AddPermission)
	admin.GET("/permissions/:role", permHandler.GetRolePermissions)
}
```

### Step 3: Use the API to add permissions

```bash
# Login as admin
TOKEN="your_admin_token"

# Add permissions
curl -X POST http://localhost:8080/api/v1/roles/permissions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_name": "admin",
    "resource": "users",
    "action": "read"
  }'
```

---

## 🎯 Recommended Permission Structure

Here's a complete permission matrix for your app:

### Admin Role (Super User)

```
users:read, users:write, users:delete
novels:read, novels:write, novels:delete
chapters:read, chapters:write, chapters:delete
genres:read, genres:write, genres:delete
tags:read, tags:write, tags:delete
roles:read, roles:write
media:read, media:write, media:delete
```

### Author Role (Content Creator)

```
novels:read, novels:write
chapters:read, chapters:write
genres:read
tags:read
profile:read, profile:write
media:write (for uploads)
```

### User Role (Regular User)

```
novels:read
chapters:read
genres:read
tags:read
profile:read, profile:write
```

### Moderator Role (Optional)

```
novels:read, novels:write
chapters:read, chapters:write, chapters:delete
users:read
media:read, media:delete
```

---

## 🧪 Testing Your Permissions

### 1. Check what permissions exist

```bash
# Connect to database
psql -U your_username -d your_database_name

# Query all permissions
SELECT
    v0 as role,
    v1 as resource,
    v2 as action
FROM casbin_rule
WHERE ptype = 'p'
ORDER BY v0, v1, v2;
```

### 2. Test with `/me` endpoint

```bash
# Get your token after login
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN" | jq '.data.permissions'
```

### 3. Test permission checks

Try accessing protected routes:

```bash
# Should work if you have permission
curl -X GET http://localhost:8080/api/v1/novels \
  -H "Authorization: Bearer YOUR_TOKEN"

# Should fail if you don't have permission
curl -X DELETE http://localhost:8080/api/v1/novels/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🔍 Troubleshooting

### Issue: Still seeing empty permissions after seeding

**Check 1: Verify data in database**

```sql
SELECT COUNT(*) FROM casbin_rule WHERE ptype = 'p';
```

If count is 0, seeding failed.

**Check 2: Verify user has the role**

```sql
SELECT * FROM user_roles WHERE user_id = YOUR_USER_ID;
```

**Check 3: Check role name matches**

```sql
SELECT name FROM roles;
```

Make sure the role name in `roles` table matches the role name in `casbin_rule` table.

**Check 4: Restart your server**
Casbin caches policies in memory. Restart to reload.

### Issue: "role not found" error

Make sure the role exists in your `roles` table:

```sql
INSERT INTO roles (name, description, created_at, updated_at)
VALUES ('admin', 'Administrator role', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;
```

---

## 📚 Next Steps

After seeding permissions:

1. ✅ Test `/me` endpoint - should now show permissions
2. ✅ Protect your routes with `CasbinMiddleware`
3. ✅ Build permission-aware frontend
4. ✅ Create admin panel for permission management

---

## 💡 Pro Tips

1. **Use consistent naming**: Always lowercase for resources (users, novels, chapters)
2. **Standard actions**: Use `read`, `write`, `delete` (not create/update/edit)
3. **Seed during deployment**: Add seed script to your CI/CD
4. **Regular backups**: Backup `casbin_rule` table
5. **Audit changes**: Log who adds/removes permissions

---

## 🎉 You're Ready!

Choose one of the methods above and seed your permissions. Your `/me` endpoint will then return the full permissions array! 🚀
