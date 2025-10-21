# Application Layer Refactoring

## Overview

This document explains the **Application Layer** architecture that was introduced to handle cross-domain operations in the application. The refactoring moved cross-domain logic from individual domain services to a dedicated application layer.

## Architecture Changes

### Before: Domain-Level Cross-Domain Logic

Previously, cross-domain operations were handled directly in domain services:

```
Role Domain (RoleService)
    ├── Has dependency on UserRepository
    ├── AssignRoleToUser() - validates both user and role
    └── RemoveRoleFromUser() - validates both user and role

User Domain (UserService)
    ├── Has dependency on RoleRepository
    └── Register() - creates user AND assigns default role
```

**Problems with this approach:**
- ❌ Domain services depended on other domains' repositories
- ❌ Tight coupling between domains
- ❌ Violation of Single Responsibility Principle
- ❌ Difficult to test and maintain

### After: Application Layer Architecture

Now, cross-domain operations are coordinated by the application layer:

```
Application Layer
    ├── UserRoleService (coordinates User + Role domains)
    │   ├── AssignRoleToUser()
    │   ├── RemoveRoleFromUser()
    │   └── GetUserRoles()
    │
    └── AuthService (coordinates User + Role domains)
        └── Register() - creates user + assigns default role

Domain Layer
    ├── UserService (pure user operations)
    │   ├── CreateUser()
    │   ├── Login()
    │   ├── GetUserByID()
    │   └── GetUserWithPermissions()
    │
    └── RoleService (pure role operations)
        ├── CreateRole()
        ├── GetRoleByID()
        ├── UpdateRole()
        └── DeleteRole()
```

**Benefits:**
- ✅ Domains are independent and loosely coupled
- ✅ Clear separation of concerns
- ✅ Cross-domain logic is explicit and centralized
- ✅ Easier to test and maintain
- ✅ Better scalability for complex workflows

## Project Structure

```
internal/
├── application/                    # Application Layer (NEW)
│   ├── dto/
│   │   └── user_role_dto.go       # DTOs for cross-domain operations
│   ├── handler/
│   │   ├── auth_handler.go        # Handles registration (user + role)
│   │   └── user_role_handler.go   # Handles user-role assignments
│   ├── service/
│   │   ├── auth_service.go        # Coordinates user creation + role assignment
│   │   └── user_role_service.go   # Coordinates user-role operations
│   └── routes.go                  # Application layer routes
│
├── user/                           # User Domain (REFACTORED)
│   ├── service/
│   │   └── user_service.go        # NO role dependencies
│   └── ...
│
├── role/                           # Role Domain (REFACTORED)
│   ├── service/
│   │   └── role_service.go        # NO user dependencies
│   └── ...
│
└── router/
    └── router.go                   # Updated to register application routes
```

## API Endpoints

### Application Layer Endpoints

**Authentication (Cross-Domain)**
- `POST /api/v1/auth/register` - Register user with default role
  ```json
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secret123",
    "confirm_password": "secret123"
  }
  ```

**User-Role Management (Cross-Domain)**
- `POST /api/v1/user-roles/assign` - Assign role to user (requires auth)
  ```json
  {
    "user_id": 1,
    "role_id": 2
  }
  ```

- `POST /api/v1/user-roles/remove` - Remove role from user (requires auth)
  ```json
  {
    "user_id": 1,
    "role_id": 2
  }
  ```

- `GET /api/v1/user-roles/users/:user_id/roles` - Get all roles for a user (requires auth)

### Domain-Specific Endpoints

**User Domain**
- `POST /api/v1/users/login` - User login
- `GET /api/v1/users/me` - Get current user with permissions
- `GET /api/v1/users/:email` - Get user by email

**Role Domain**
- `POST /api/v1/roles` - Create role
- `GET /api/v1/roles` - List all roles
- `GET /api/v1/roles/:id` - Get role by ID
- `PUT /api/v1/roles/:id` - Update role
- `DELETE /api/v1/roles/:id` - Delete role

## Key Refactoring Changes

### 1. RoleService (Removed Cross-Domain Logic)

**Before:**
```go
type RoleService struct {
    roleRepo *repository.RoleRepository
    userRepo *userRepo.UserRepository  // ❌ Dependency on user domain
}

func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
    // Cross-domain validation and assignment
}
```

**After:**
```go
type RoleService struct {
    roleRepo *repository.RoleRepository  // ✅ Only role domain dependency
}

// AssignRoleToUser moved to application layer
```

### 2. UserService (Removed Cross-Domain Logic)

**Before:**
```go
type UserService struct {
    userRepo      *repository.UserRepository
    roleRepo      *roleRepo.RoleRepository  // ❌ Dependency on role domain
    casbinService *casbinService.CasbinService
}

func (s *UserService) Register(dto *dto.RegisterUserDTO) (*model.User, error) {
    // Creates user AND assigns default role
}
```

**After:**
```go
type UserService struct {
    userRepo      *repository.UserRepository  // ✅ Only user domain dependency
    casbinService *casbinService.CasbinService
}

func (s *UserService) CreateUser(dto *dto.RegisterUserDTO) (*model.User, error) {
    // Only creates user, no role logic
}
```

### 3. Application Layer (NEW)

**AuthService** - Handles user registration with role assignment:
```go
type AuthService struct {
    userRepo      *userRepo.UserRepository
    roleRepo      *roleRepo.RoleRepository
    casbinService *casbinService.CasbinService
    db            *gorm.DB
}

func (s *AuthService) Register(dto *userDto.RegisterUserDTO) (*userModel.User, error) {
    // Transaction ensures atomicity:
    // 1. Create user
    // 2. Get default "user" role
    // 3. Assign role in database
    // 4. Sync with Casbin
}
```

**UserRoleService** - Handles user-role relationships:
```go
type UserRoleService struct {
    userRepo      *userRepo.UserRepository
    roleRepo      *roleRepo.RoleRepository
    casbinService *casbinService.CasbinService
    db            *gorm.DB
}

func (s *UserRoleService) AssignRoleToUser(userID, roleID uint) error {
    // 1. Validate user exists
    // 2. Validate role exists
    // 3. Check if already assigned
    // 4. Assign in database
    // 5. Sync with Casbin
}
```

## Transaction Handling

The application layer uses database transactions to ensure atomicity of cross-domain operations:

```go
err := database.DB.Transaction(func(tx *gorm.DB) error {
    // All operations within this block are atomic
    
    // 1. Domain operation 1
    user, err := s.userRepo.WithTx(tx).CreateUser(user)
    
    // 2. Domain operation 2
    role, err := s.roleRepo.WithTx(tx).GetByName("user")
    
    // 3. Cross-domain operation
    err := s.userRepo.WithTx(tx).AssignRoleToUser(userID, roleID)
    
    // 4. Authorization sync
    err := s.casbinService.AssignRoleToUser(userID, roleName)
    
    return nil // Commit transaction
})
```

## Benefits of This Architecture

### 1. **Clear Domain Boundaries**
- Each domain focuses only on its own concerns
- No circular dependencies between domains
- Easier to understand and navigate codebase

### 2. **Testability**
- Domain services can be tested independently
- Application layer can mock domain services
- Easier to write unit tests

### 3. **Maintainability**
- Changes in one domain don't affect others
- Cross-domain logic is centralized
- Easier to add new features

### 4. **Scalability**
- Can add new domains without affecting existing ones
- Complex workflows can be orchestrated at application layer
- Supports future microservices migration

### 5. **Consistency**
- Transaction management at application layer
- Ensures data consistency across domains
- Single source of truth for cross-domain operations

## Migration Guide

If you need to add new cross-domain operations:

1. **Identify if it's truly cross-domain**
   - Does it require data/logic from multiple domains?
   - If yes → Application Layer
   - If no → Stay in domain layer

2. **Create Application Service**
   ```go
   // internal/application/service/your_service.go
   type YourService struct {
       domain1Repo *Domain1Repository
       domain2Repo *Domain2Repository
       db          *gorm.DB
   }
   ```

3. **Create Handler**
   ```go
   // internal/application/handler/your_handler.go
   type YourHandler struct {
       yourService *service.YourService
   }
   ```

4. **Register Routes**
   ```go
   // internal/application/routes.go
   yourService := service.NewYourService(...)
   yourHandler := handler.NewYourHandler(yourService)
   router.POST("/your-endpoint", yourHandler.YourMethod)
   ```

## Related Documentation

- [Architecture Overview](../02-architecture/ARCHITECTURE.md)
- [Cross-Domain Relationships](../02-architecture/CROSS_DOMAIN_RELATIONSHIPS.md)
- [RBAC Implementation](../04-authorization/RBAC_IMPLEMENTATION.md)
- [Authentication](../03-authentication/AUTHENTICATION.md)

## Summary

The application layer refactoring successfully:

- ✅ Removed cross-domain dependencies from domain services
- ✅ Centralized cross-domain logic in application layer
- ✅ Improved separation of concerns
- ✅ Made the codebase more maintainable and testable
- ✅ Prepared the architecture for future scaling

**Result:** A cleaner, more maintainable architecture that follows Domain-Driven Design (DDD) principles and supports both current and future requirements.
