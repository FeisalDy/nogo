# Database Schema Changes - Casbin Integration

## What Changed

### ✅ Removed Tables
- `permissions` - No longer needed (Casbin uses `casbin_rule`)
- `role_permissions` - No longer needed (Casbin uses `casbin_rule`)

### ✅ Kept Tables
- `roles` - Still used for role definitions (synced with Casbin)
- `user_roles` - Still used for user-role assignments (synced with Casbin)
- `users` - Unchanged

### ✅ New Table (Auto-created by Casbin)
- `casbin_rule` - Stores all policies and role assignments

## Database Schema

### Current Schema (After Cleanup)

```sql
-- Users table (existing)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    avatar_url VARCHAR(255),
    bio TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Roles table (existing, synced with Casbin)
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- User-Role assignments (existing, synced with Casbin)
CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- Casbin policies (auto-created by Casbin GORM adapter)
CREATE TABLE casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),  -- 'p' for policy, 'g' for grouping/role
    v0 VARCHAR(100),     -- subject (user/role)
    v1 VARCHAR(100),     -- object (resource)
    v2 VARCHAR(100),     -- action
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

## How Casbin Uses These Tables

### Roles Table
- **Purpose**: Store role definitions (name, description)
- **Casbin Usage**: Role names are referenced in `casbin_rule` table
- **Example**:
  ```sql
  INSERT INTO roles (name, description) 
  VALUES ('admin', 'Administrator with full access');
  ```

### User_Roles Table
- **Purpose**: Link users to roles
- **Casbin Sync**: When you assign a role, update both `user_roles` AND Casbin
- **Example**:
  ```sql
  INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);
  -- Also sync with Casbin:
  -- svc.AssignRoleToUser(1, "admin")
  ```

### Casbin_Rule Table
- **Purpose**: Store all permission policies and role assignments
- **Managed By**: Casbin (don't manually edit!)
- **Example Data**:

```sql
-- Policy: admin role can write to users resource
ptype='p', v0='admin', v1='users', v2='write'

-- Policy: editor role can write to novels resource
ptype='p', v0='editor', v1='novels', v2='write'

-- Grouping: user 1 has admin role
ptype='g', v0='user:1', v1='admin'

-- Grouping: user 2 has editor role
ptype='g', v0='user:2', v1='editor'
```

## Migration Steps

### Step 1: Drop Old Tables (If They Exist)

Run the SQL script:
```bash
psql -U your_user -d your_database -f scripts/drop_old_permission_tables.sql
```

Or manually in psql:
```sql
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
```

### Step 2: Verify Tables

```sql
\dt  -- List all tables

-- Should see:
-- users
-- roles
-- user_roles
-- casbin_rule (will be created when you start the server)
-- migrations (migration tracking)
```

### Step 3: Start Server

Casbin will automatically create the `casbin_rule` table on first run:

```bash
go run cmd/server/main.go
```

## Data Flow

### Creating a Role
```
1. Create in database (roles table)
   ↓
2. Define permissions in Casbin
   ↓
3. Casbin stores in casbin_rule table
```

### Assigning Role to User
```
1. Assign in database (user_roles table)
   ↓
2. Sync with Casbin
   ↓
3. Casbin stores in casbin_rule table (ptype='g')
```

### Checking Permission
```
1. User makes request with JWT
   ↓
2. Middleware extracts user_id
   ↓
3. Casbin checks:
   - Does user have role? (casbin_rule where ptype='g')
   - Does role have permission? (casbin_rule where ptype='p')
   ↓
4. Allow or Deny
```

## Example: Complete Workflow

```go
// 1. Create role in database
role := &model.Role{
    Name:        "editor",
    Description: "Content editor",
}
db.Create(&role)

// 2. Add permissions to role in Casbin
svc := casbinService.NewCasbinService()
svc.AddPermissionForRole("editor", "novels", "read")
svc.AddPermissionForRole("editor", "novels", "write")
svc.AddPermissionForRole("editor", "chapters", "write")

// 3. Assign role to user in database
userRole := &model.UserRole{
    UserID: 123,
    RoleID: role.ID,
}
db.Create(&userRole)

// 4. Sync with Casbin
svc.AssignRoleToUser(123, "editor")

// 5. Check permission (done automatically by middleware)
allowed, _ := svc.Enforce(123, "novels", "write")
// Returns: true
```

## Benefits of This Approach

✅ **Simpler Schema** - Only 3 tables instead of 5
✅ **Flexible Permissions** - Add/remove/modify without schema changes
✅ **Better Performance** - Casbin caches policies in memory
✅ **Industry Standard** - Casbin is battle-tested and widely used
✅ **Dynamic** - Change permissions at runtime without restarts

## Important Notes

⚠️ **Don't Edit casbin_rule Manually**
- Let Casbin manage this table
- Use CasbinService methods instead

⚠️ **Keep Sync Between Tables**
- When adding role in `roles` table, add permissions in Casbin
- When assigning in `user_roles` table, sync with Casbin

⚠️ **Existing Data**
- If you have existing roles in database, you need to define their permissions in Casbin
- If you have existing user-role assignments, sync them to Casbin

## Troubleshooting

### Permission checks not working?
```go
// Check if user has role in Casbin
svc := casbinService.NewCasbinService()
roles, _ := svc.GetRolesForUser(userID)
fmt.Println(roles) // Should show user's roles

// Check if role has permission
perms, _ := svc.GetPermissionsForRole("admin")
fmt.Println(perms) // Should show role's permissions
```

### Tables not synced?
```go
// Reload policies from database
svc.ReloadPolicies()
```

### Need to start fresh?
```sql
TRUNCATE TABLE casbin_rule;
```
Then re-seed permissions.

## Next Steps

1. ✅ Migration updated
2. ✅ Old tables removed
3. ✅ Code cleaned up
4. 📝 Next: Seed initial permissions (see CASBIN_QUICK_START.md)
5. 📝 Next: Update routes with middleware (see CASBIN_ROUTE_EXAMPLES.md)
