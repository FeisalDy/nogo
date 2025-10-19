# Role-Based Access Control (RBAC) Implementation Guide

## Overview

This guide documents the complete implementation of Role-Based Access Control (RBAC) in the multi-domain Go/Gin architecture, demonstrating how to properly handle many-to-many relationships between domains while maintaining clean architecture principles.

## Architecture Pattern

We used **Approach 1** from the cross-domain relationships guide:

- **Junction Table in Common**: `user_roles` table in `internal/common/model`
- **Separate Domain Models**: `User` and `Role` models remain in their respective domains
- **Repository Layer Relationships**: Cross-domain queries are handled in repositories using SQL joins
- **No Circular Dependencies**: Clean separation maintained throughout

## File Structure

```
internal/
├── common/
│   └── model/
│       └── user_role.go          # Junction table model
├── role/
│   ├── dto/
│   │   └── role_dto.go           # Data transfer objects
│   ├── handler/
│   │   └── role_handler.go       # HTTP handlers
│   ├── model/
│   │   └── role.go               # Role model
│   ├── repository/
│   │   └── role_repository.go    # Data access layer
│   ├── service/
│   │   └── role_service.go       # Business logic
│   └── routes.go                 # Route registration
├── user/
│   └── repository/
│       └── user_repository.go    # Extended with role methods
└── database/
    └── migrations/
        └── 008_create_roles_and_user_roles.go  # Database migration
```

## Database Schema

### Roles Table

```sql
CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);
```

### User_Roles Junction Table

```sql
CREATE TABLE user_roles (
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL,
  PRIMARY KEY (user_id, role_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
```

## Implementation Details

### 1. Models

#### Role Model (`internal/role/model/role.go`)

```go
type Role struct {
    gorm.Model
    Name        string  `gorm:"unique;not null;size:50"`
    Description *string `gorm:"type:text"`
}
```

#### UserRole Junction Table (`internal/common/model/user_role.go`)

```go
type UserRole struct {
    UserID    uint `gorm:"primaryKey"`
    RoleID    uint `gorm:"primaryKey"`
    CreatedAt int64
    UpdatedAt int64
}
```

### 2. DTOs

#### Role DTOs (`internal/role/dto/role_dto.go`)

- `RoleDTO`: Response DTO for role data
- `CreateRoleDTO`: Request DTO for creating a role (with validation)
- `UpdateRoleDTO`: Request DTO for updating a role (all fields optional)
- `RoleWithUsersDTO`: Response DTO with role and assigned users
- `UserDTO`: Simplified user representation for role responses

### 3. Repository Layer

#### Role Repository (`internal/role/repository/role_repository.go`)

**CRUD Operations:**

- `Create(role)`: Create a new role
- `GetByID(id)`: Get role by ID
- `GetByName(name)`: Get role by name
- `GetAll()`: Get all roles
- `Update(role)`: Update a role
- `Delete(id)`: Delete a role

**Relationship Operations:**

- `GetRoleWithUsers(roleID)`: Get a role with all assigned users
- `GetUserCountByRole(roleID)`: Count users with a specific role
- `Exists(id)`: Check if role exists by ID
- `ExistsByName(name)`: Check if role exists by name

#### User Repository Extended (`internal/user/repository/user_repository.go`)

**Role Relationship Methods:**

- `GetUserWithRoles(userID)`: Get user with all their roles
- `AssignRoleToUser(userID, roleID)`: Assign a role to a user
- `RemoveRoleFromUser(userID, roleID)`: Remove a role from a user
- `GetUserRoleIDs(userID)`: Get all role IDs for a user
- `HasRole(userID, roleName)`: Check if user has a specific role (by name)
- `HasRoleByID(userID, roleID)`: Check if user has a specific role (by ID)
- `HasAnyRole(userID, roleNames)`: Check if user has any of the specified roles

### 4. Service Layer

#### Role Service (`internal/role/service/role_service.go`)

The service layer handles business logic and coordinates between repositories:

**Constructor:**

```go
NewRoleService(roleRepo, userRepo) *RoleService
```

**Methods:**

- `CreateRole(req)`: Creates a role with uniqueness validation
- `GetRoleByID(id)`: Retrieves a role by ID
- `GetRoleByName(name)`: Retrieves a role by name
- `GetAllRoles()`: Retrieves all roles
- `UpdateRole(id, req)`: Updates a role with conflict checking
- `DeleteRole(id)`: Deletes a role (prevents deletion if users assigned)
- `GetRoleWithUsers(id)`: Gets a role with all assigned users
- `AssignRoleToUser(userID, roleID)`: Assigns a role to a user
- `RemoveRoleFromUser(userID, roleID)`: Removes a role from a user

**Error Codes Used:**

- `ROLE001`: Role with this name already exists
- `ROLE002`: Role not found
- `ROLE003`: Cannot delete role with assigned users
- `ROLE004`: User already has this role
- `ROLE005`: User does not have this role
- `USER003`: User not found

### 5. Handler Layer

#### Role Handler (`internal/role/handler/role_handler.go`)

HTTP endpoints with proper error handling and validation:

**Endpoints:**

- `POST /roles`: Create a new role
- `GET /roles`: Get all roles
- `GET /roles/:id`: Get a role by ID
- `PUT /roles/:id`: Update a role
- `DELETE /roles/:id`: Delete a role
- `GET /roles/:id/users`: Get a role with all assigned users
- `POST /roles/:role_id/users/:user_id`: Assign a role to a user
- `DELETE /roles/:role_id/users/:user_id`: Remove a role from a user

**Features:**

- Input validation using validator v10
- Standardized error responses using AppError system
- Consistent response format using RespondSuccess/RespondWithAppError
- Swagger documentation comments

### 6. Routes

#### Route Registration (`internal/role/routes.go`)

```go
func RegisterRoutes(router *gin.RouterGroup) {
    roleRepository := repository.NewRoleRepository()
    userRepository := userRepo.NewUserRepository()
    roleService := service.NewRoleService(roleRepository, userRepository)
    roleHandler := handler.NewRoleHandler(roleService)

    roleRoutes := router.Group("/roles")
    roleRoutes.Use(middleware.AuthMiddleware())
    {
        // CRUD operations
        roleRoutes.POST("", roleHandler.CreateRole)
        roleRoutes.GET("", roleHandler.GetAllRoles)
        roleRoutes.GET("/:id", roleHandler.GetRole)
        roleRoutes.PUT("/:id", roleHandler.UpdateRole)
        roleRoutes.DELETE("/:id", roleHandler.DeleteRole)

        // Role-user relationships
        roleRoutes.GET("/:id/users", roleHandler.GetRoleWithUsers)
        roleRoutes.POST("/:role_id/users/:user_id", roleHandler.AssignUserToRole)
        roleRoutes.DELETE("/:role_id/users/:user_id", roleHandler.RemoveUserFromRole)
    }
}
```

Routes are registered in `internal/router/router.go`:

```go
roleRoutes := v1.Group("/roles")
role.RegisterRoutes(roleRoutes)
```

### 7. Migration

#### Migration 008 (`internal/database/migrations/008_create_roles_and_user_roles.go`)

**Features:**

- Creates `roles` and `user_roles` tables
- Adds foreign key constraints with CASCADE delete
- Inserts default roles: `admin`, `editor`, `user`
- Idempotent (safe to run multiple times)
- Includes rollback (Down) functionality

**Default Roles:**

- **admin**: Administrator with full access
- **editor**: Can create and edit content
- **user**: Regular user with basic access

## API Usage Examples

### 1. Create a Role

```bash
POST /api/v1/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "moderator",
  "description": "Can moderate content"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "id": 4,
    "name": "moderator",
    "description": "Can moderate content"
  },
  "message": "Role created successfully"
}
```

### 2. Get All Roles

```bash
GET /api/v1/roles
Authorization: Bearer <token>
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator with full access"
    },
    {
      "id": 2,
      "name": "editor",
      "description": "Can create and edit content"
    },
    {
      "id": 3,
      "name": "user",
      "description": "Regular user with basic access"
    }
  ],
  "message": "Roles retrieved successfully"
}
```

### 3. Assign Role to User

```bash
POST /api/v1/roles/1/users/5
Authorization: Bearer <token>
```

Response:

```json
{
  "success": true,
  "data": {
    "message": "Role assigned to user successfully"
  },
  "message": "Role assigned to user successfully"
}
```

### 4. Get Role with Users

```bash
GET /api/v1/roles/1/users
Authorization: Bearer <token>
```

Response:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "admin",
    "description": "Administrator with full access",
    "users": [
      {
        "id": 5,
        "username": "john_doe",
        "email": "john@example.com",
        "status": "active"
      }
    ],
    "user_count": 1
  },
  "message": "Role with users retrieved successfully"
}
```

### 5. Remove Role from User

```bash
DELETE /api/v1/roles/1/users/5
Authorization: Bearer <token>
```

## Key Design Decisions

### 1. Junction Table in Common Domain

**Why:** Prevents circular dependencies between User and Role domains. The `user_roles` table is a pure relationship without business logic, making it suitable for the common layer.

### 2. Repository Layer for Cross-Domain Queries

**Why:** Keeps domain services focused on their own domain logic. Repositories handle database-level relationships using SQL joins.

### 3. Service Layer Coordinates Multiple Repositories

**Why:** Business logic that involves multiple domains (like assigning roles) belongs in the service layer, which orchestrates repository calls.

### 4. Composite Primary Key for Junction Table

**Why:** Ensures uniqueness of user-role combinations at the database level and improves query performance.

### 5. No Circular Dependencies

**Why:** Clean architecture principles require that domains don't directly depend on each other. We achieve this through:

- Junction table in common
- Import aliasing (`userModel`, `roleModel`)
- Repository layer handling cross-domain queries

## Testing the Implementation

### 1. Run Migrations

```bash
cd /home/feisal/project/shilan/nogo
go run cmd/migrate/main.go
```

This will:

- Create the `roles` table
- Create the `user_roles` junction table
- Add foreign key constraints
- Insert default roles (admin, editor, user)

### 2. Start the Server

```bash
go run cmd/server/main.go
```

### 3. Test the Endpoints

First, register a user and get a JWT token:

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123!",
    "confirm_password": "Password123!"
  }'
```

Then use the token to access role endpoints:

```bash
# Get all roles
curl -X GET http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer <your-token>"

# Assign admin role (ID: 1) to user (ID: 1)
curl -X POST http://localhost:8080/api/v1/roles/1/users/1 \
  -H "Authorization: Bearer <your-token>"
```

## Common Pitfalls to Avoid

1. **Don't import user domain in role service directly**

   - ✅ Import `userRepo "github.com/FeisalDy/nogo/internal/user/repository"`
   - ❌ Import `"github.com/FeisalDy/nogo/internal/user/service"`

2. **Don't put business logic in repository**

   - ✅ Repository only does database operations
   - ❌ Repository doesn't check business rules

3. **Don't skip validation in service layer**

   - ✅ Check if role exists before assigning
   - ❌ Assume data is valid just because it passed handler validation

4. **Don't forget cascade deletes**

   - ✅ Foreign keys with ON DELETE CASCADE
   - ❌ Orphaned records in junction table

5. **Don't expose internal IDs without checking permissions**
   - ✅ All role endpoints require authentication
   - ❌ Public endpoints that expose user-role relationships

## Next Steps

Now that you have RBAC implemented, you can:

1. **Add Role-Based Authorization Middleware**

   ```go
   func RequireRole(role string) gin.HandlerFunc {
       // Check if authenticated user has the required role
   }
   ```

2. **Extend User DTOs to Include Roles**

   ```go
   type UserWithRolesDTO struct {
       UserDTO
       Roles []RoleDTO `json:"roles"`
   }
   ```

3. **Add Permission-Based Access Control**

   - Create `permissions` table
   - Create `role_permissions` junction table
   - Implement fine-grained access control

4. **Add Role Hierarchy**
   - Super admin > Admin > Editor > User
   - Implement inheritance of permissions

## Conclusion

This implementation demonstrates clean architecture principles for handling many-to-many relationships in a multi-domain Go application:

- **Separation of Concerns**: Each layer has a clear responsibility
- **No Circular Dependencies**: Achieved through junction table in common and repository-level queries
- **Maintainability**: Easy to extend with new roles or permissions
- **Testability**: Each component can be tested independently
- **Performance**: Efficient queries using SQL joins instead of N+1 queries

The pattern used here (junction table in common, repository layer handling relationships) can be applied to other many-to-many relationships in your application.
