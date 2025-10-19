# Cross-Domain Relationships Guide

This guide explains how to handle relationships between domains in a multi-domain architecture, specifically focusing on many-to-many relationships like User-Role.

## Problem

In a multi-domain architecture:

- User domain is in `internal/user/`
- Role domain is in `internal/role/`
- They need a many-to-many relationship
- We want to maintain domain separation
- We want to avoid circular dependencies

## Three Approaches

### Approach 1: Junction Table in Common (RECOMMENDED) ⭐

Create the junction/pivot table in a shared location since it relates to both domains.

**Pros:**

- ✅ No circular dependencies
- ✅ Clear separation of concerns
- ✅ Easy to query from either domain
- ✅ Follows DDD principles

**Cons:**

- ⚠️ Requires a common/shared models location

### Approach 2: Embed in One Domain

Put the relationship in the "owner" domain (usually the one you query from most).

**Pros:**

- ✅ Simpler structure
- ✅ Clear ownership

**Cons:**

- ❌ One domain depends on another
- ❌ Harder to query from the other direction
- ❌ Can cause circular dependencies

### Approach 3: Separate Relationship Domain

Create a dedicated domain for the relationship.

**Pros:**

- ✅ Maximum separation
- ✅ Scalable for complex relationships
- ✅ Clear bounded context

**Cons:**

- ❌ More complex
- ❌ More boilerplate
- ❌ Overkill for simple relationships

## Recommended Implementation (Approach 1)

### 1. Directory Structure

```
internal/
├── common/
│   └── model/
│       └── user_role.go          # Junction table
├── user/
│   ├── model/
│   │   └── user.go               # User model
│   ├── repository/
│   │   └── user_repository.go
│   └── service/
│       └── user_service.go
└── role/
    ├── model/
    │   └── role.go               # Role model
    ├── repository/
    │   └── role_repository.go
    └── service/
        └── role_service.go
```

### 2. Models

#### User Model (`internal/user/model/user.go`)

```go
package model

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model

    Username  *string `json:"username"`
    Email     string  `json:"email" gorm:"unique;not null"`
    Password  *string `json:"-"`
    AvatarURL *string `json:"avatar_url"`
    Bio       *string `json:"bio" gorm:"type:text"`
    Status    string  `json:"status" gorm:"default:'active';index"`

    // Don't define Roles here to avoid importing role package
}
```

#### Role Model (`internal/role/model/role.go`)

```go
package model

import "gorm.io/gorm"

type Role struct {
    gorm.Model

    Name        string  `json:"name" gorm:"unique;not null"`
    Description *string `json:"description" gorm:"type:text"`

    // Don't define Users here to avoid circular dependency
}
```

#### Junction Table (`internal/common/model/user_role.go`)

```go
package model

import "time"

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
    UserID    uint      `gorm:"primaryKey;index:idx_user_role"`
    RoleID    uint      `gorm:"primaryKey;index:idx_user_role"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// TableName specifies the table name
func (UserRole) TableName() string {
    return "user_roles"
}
```

### 3. Migrations

```go
// internal/database/migrations/008_create_user_roles.go
package migrations

import (
    commonModel "github.com/FeisalDy/nogo/internal/common/model"
    "gorm.io/gorm"
)

func init() {
    migrationsList = append(migrationsList, Migration{
        ID: "008_create_user_roles",
        Migrate: func(db *gorm.DB) error {
            return db.AutoMigrate(&commonModel.UserRole{})
        },
        Rollback: func(db *gorm.DB) error {
            return db.Migrator().DropTable("user_roles")
        },
    })
}
```

### 4. Repository Methods

#### User Repository - Role Operations

```go
// internal/user/repository/user_repository.go
package repository

import (
    "github.com/FeisalDy/nogo/internal/common/model"
    "github.com/FeisalDy/nogo/internal/database"
    userModel "github.com/FeisalDy/nogo/internal/user/model"
    roleModel "github.com/FeisalDy/nogo/internal/role/model"
)

type UserRepository struct{}

// GetUserWithRoles gets a user and their roles
func (r *UserRepository) GetUserWithRoles(userID uint) (*userModel.User, []roleModel.Role, error) {
    var user userModel.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        return nil, nil, err
    }

    var roles []roleModel.Role
    err := database.DB.
        Table("roles").
        Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
        Where("user_roles.user_id = ?", userID).
        Find(&roles).Error

    return &user, roles, err
}

// AssignRoleToUser assigns a role to a user
func (r *UserRepository) AssignRoleToUser(userID, roleID uint) error {
    userRole := model.UserRole{
        UserID: userID,
        RoleID: roleID,
    }
    return database.DB.Create(&userRole).Error
}

// RemoveRoleFromUser removes a role from a user
func (r *UserRepository) RemoveRoleFromUser(userID, roleID uint) error {
    return database.DB.
        Where("user_id = ? AND role_id = ?", userID, roleID).
        Delete(&model.UserRole{}).Error
}

// GetUserRoleIDs gets all role IDs for a user
func (r *UserRepository) GetUserRoleIDs(userID uint) ([]uint, error) {
    var roleIDs []uint
    err := database.DB.
        Model(&model.UserRole{}).
        Where("user_id = ?", userID).
        Pluck("role_id", &roleIDs).Error
    return roleIDs, err
}

// HasRole checks if a user has a specific role
func (r *UserRepository) HasRole(userID uint, roleName string) (bool, error) {
    var count int64
    err := database.DB.
        Table("user_roles").
        Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
        Where("user_roles.user_id = ? AND roles.name = ?", userID, roleName).
        Count(&count).Error
    return count > 0, err
}
```

#### Role Repository - User Operations

```go
// internal/role/repository/role_repository.go
package repository

import (
    "github.com/FeisalDy/nogo/internal/common/model"
    "github.com/FeisalDy/nogo/internal/database"
    roleModel "github.com/FeisalDy/nogo/internal/role/model"
    userModel "github.com/FeisalDy/nogo/internal/user/model"
)

type RoleRepository struct{}

// GetRoleWithUsers gets a role and all users with that role
func (r *RoleRepository) GetRoleWithUsers(roleID uint) (*roleModel.Role, []userModel.User, error) {
    var role roleModel.Role
    if err := database.DB.First(&role, roleID).Error; err != nil {
        return nil, nil, err
    }

    var users []userModel.User
    err := database.DB.
        Table("users").
        Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
        Where("user_roles.role_id = ?", roleID).
        Find(&users).Error

    return &role, users, err
}

// GetRoleByName gets a role by name
func (r *RoleRepository) GetRoleByName(name string) (*roleModel.Role, error) {
    var role roleModel.Role
    err := database.DB.Where("name = ?", name).First(&role).Error
    return &role, err
}

// GetUserCountByRole gets the count of users with a specific role
func (r *RoleRepository) GetUserCountByRole(roleID uint) (int64, error) {
    var count int64
    err := database.DB.
        Model(&model.UserRole{}).
        Where("role_id = ?", roleID).
        Count(&count).Error
    return count, err
}
```

### 5. DTOs for Cross-Domain Data

```go
// internal/user/dto/user_dto.go
package dto

// UserWithRolesDTO represents a user with their roles
type UserWithRolesDTO struct {
    ID        uint        `json:"id"`
    Username  *string     `json:"username"`
    Email     string      `json:"email"`
    Status    string      `json:"status"`
    Roles     []RoleDTO   `json:"roles"`
}

// RoleDTO represents a role (from role domain)
type RoleDTO struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Description *string `json:"description,omitempty"`
}

// AssignRoleDTO for assigning roles to users
type AssignRoleDTO struct {
    RoleID uint `json:"role_id" validate:"required"`
}
```

```go
// internal/role/dto/role_dto.go
package dto

// RoleWithUsersDTO represents a role with users
type RoleWithUsersDTO struct {
    ID          uint       `json:"id"`
    Name        string     `json:"name"`
    Description *string    `json:"description,omitempty"`
    Users       []UserDTO  `json:"users"`
    UserCount   int64      `json:"user_count"`
}

// UserDTO represents a user (from user domain)
type UserDTO struct {
    ID       uint    `json:"id"`
    Username *string `json:"username"`
    Email    string  `json:"email"`
    Status   string  `json:"status"`
}
```

### 6. Service Layer

```go
// internal/user/service/user_service.go
package service

import (
    "github.com/FeisalDy/nogo/internal/user/dto"
    "github.com/FeisalDy/nogo/internal/user/repository"
    roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
)

type UserService struct {
    UserRepository *repository.UserRepository
    RoleRepository *roleRepo.RoleRepository  // Inject role repository
}

// GetUserWithRoles gets user with their roles
func (s *UserService) GetUserWithRoles(userID uint) (*dto.UserWithRolesDTO, error) {
    user, roles, err := s.UserRepository.GetUserWithRoles(userID)
    if err != nil {
        return nil, err
    }

    // Convert to DTO
    roleDTOs := make([]dto.RoleDTO, len(roles))
    for i, role := range roles {
        roleDTOs[i] = dto.RoleDTO{
            ID:          role.ID,
            Name:        role.Name,
            Description: role.Description,
        }
    }

    return &dto.UserWithRolesDTO{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        Status:   user.Status,
        Roles:    roleDTOs,
    }, nil
}

// AssignRole assigns a role to a user
func (s *UserService) AssignRole(userID, roleID uint) error {
    // Verify role exists
    if _, err := s.RoleRepository.GetRoleByID(roleID); err != nil {
        return err
    }

    return s.UserRepository.AssignRoleToUser(userID, roleID)
}

// RemoveRole removes a role from a user
func (s *UserService) RemoveRole(userID, roleID uint) error {
    return s.UserRepository.RemoveRoleFromUser(userID, roleID)
}

// HasRole checks if user has a specific role
func (s *UserService) HasRole(userID uint, roleName string) (bool, error) {
    return s.UserRepository.HasRole(userID, roleName)
}
```

### 7. Handler with Cross-Domain Operations

```go
// internal/user/handler/user_handler.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/FeisalDy/nogo/internal/user/dto"
    "github.com/FeisalDy/nogo/internal/user/service"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    UserService *service.UserService
}

// GetUserWithRoles gets a user with their roles
func (h *UserHandler) GetUserWithRoles(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrInvalidInput)
        return
    }

    userWithRoles, err := h.UserService.GetUserWithRoles(uint(id))
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrUserNotFound)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, userWithRoles)
}

// AssignRole assigns a role to a user
func (h *UserHandler) AssignRole(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrInvalidInput)
        return
    }

    var assignDTO dto.AssignRoleDTO
    if err := c.ShouldBindJSON(&assignDTO); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
        return
    }

    if err := h.UserService.AssignRole(uint(id), assignDTO.RoleID); err != nil {
        utils.RespondWithAppError(c, errors.ErrInternalServer)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, nil, "Role assigned successfully")
}

// RemoveRole removes a role from a user
func (h *UserHandler) RemoveRole(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrInvalidInput)
        return
    }

    roleID, err := strconv.ParseUint(c.Param("roleId"), 10, 32)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrInvalidInput)
        return
    }

    if err := h.UserService.RemoveRole(uint(id), uint(roleID)); err != nil {
        utils.RespondWithAppError(c, errors.ErrInternalServer)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, nil, "Role removed successfully")
}
```

### 8. Routes

```go
// internal/user/routes.go
package user

import (
    "github.com/FeisalDy/nogo/internal/common/middleware"
    "github.com/FeisalDy/nogo/internal/user/handler"
    "github.com/FeisalDy/nogo/internal/user/repository"
    "github.com/FeisalDy/nogo/internal/user/service"
    roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
    userRepository := repository.NewUserRepository()
    roleRepository := roleRepo.NewRoleRepository()
    userService := service.NewUserService(userRepository, roleRepository)
    userHandler := handler.NewUserHandler(userService)

    // Public routes
    router.POST("/register", userHandler.Register)
    router.POST("/login", userHandler.Login)

    // Protected routes
    protected := router.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/me", userHandler.GetMe)
        protected.GET("/:id", userHandler.GetUser)
        protected.GET("/:id/roles", userHandler.GetUserWithRoles)
        protected.POST("/:id/roles", userHandler.AssignRole)
        protected.DELETE("/:id/roles/:roleId", userHandler.RemoveRole)
    }
}
```

## Best Practices

### 1. Avoid Circular Dependencies

❌ **Don't:**

```go
// user/model/user.go
import "github.com/FeisalDy/nogo/internal/role/model"

type User struct {
    Roles []model.Role `gorm:"many2many:user_roles"`
}

// role/model/role.go
import "github.com/FeisalDy/nogo/internal/user/model"

type Role struct {
    Users []model.User `gorm:"many2many:user_roles"`
}
// This creates circular dependency!
```

✅ **Do:**

```go
// Keep models separate
// Use junction table in common
// Query relationships in repository layer
```

### 2. Use DTOs for Cross-Domain Data

✅ **Do:**

```go
// Return DTOs, not domain models
type UserWithRolesDTO struct {
    ID    uint      `json:"id"`
    Email string    `json:"email"`
    Roles []RoleDTO `json:"roles"` // Simple DTO, not full model
}
```

### 3. Inject Dependencies

✅ **Do:**

```go
type UserService struct {
    UserRepository *userRepo.UserRepository
    RoleRepository *roleRepo.RoleRepository // Inject needed repositories
}
```

### 4. Keep Business Logic in Service Layer

✅ **Do:**

```go
// Service handles cross-domain operations
func (s *UserService) AssignRole(userID, roleID uint) error {
    // Validate role exists
    role, err := s.RoleRepository.GetRoleByID(roleID)
    if err != nil {
        return errors.New("role not found")
    }

    // Check if already assigned
    hasRole, _ := s.UserRepository.HasRole(userID, role.Name)
    if hasRole {
        return errors.New("role already assigned")
    }

    // Assign
    return s.UserRepository.AssignRoleToUser(userID, roleID)
}
```

## Summary

For many-to-many relationships in multi-domain architecture:

1. **Create junction table in `common/model`**
2. **Keep domain models separate** (no cross-imports)
3. **Handle relationships in repository layer** with SQL joins
4. **Use DTOs** to transfer cross-domain data
5. **Inject dependencies** at service layer
6. **Keep business logic** in service layer

This approach maintains clean architecture while handling complex relationships efficiently!
