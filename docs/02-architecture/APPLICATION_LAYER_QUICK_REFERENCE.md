# Application Layer Quick Reference

## When to Use Application Layer vs Domain Layer

### Use Domain Layer When:
- ✅ Operation involves only ONE domain
- ✅ No coordination between multiple domains needed
- ✅ Simple CRUD operations
- ✅ Business logic specific to one domain

**Examples:**
- Create a role
- Get user by ID
- Update user profile
- Delete a role

### Use Application Layer When:
- ✅ Operation involves MULTIPLE domains
- ✅ Coordination between domains is required
- ✅ Complex workflows that span domains
- ✅ Business processes (not just data operations)

**Examples:**
- User registration (User + Role + Casbin)
- Assign role to user (User + Role + Casbin)
- Order checkout (Order + Payment + Inventory + Notification)

## Quick Architecture Comparison

```
┌─────────────────────────────────────────────────┐
│           Application Layer                      │
│  Coordinates cross-domain operations            │
│                                                  │
│  ┌──────────────┐      ┌──────────────┐        │
│  │ AuthService  │      │UserRoleService│        │
│  │              │      │              │         │
│  │ Register()   │      │ AssignRole() │         │
│  └──────┬───────┘      └──────┬───────┘         │
└─────────┼────────────────────┼──────────────────┘
          │                    │
          ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│  User Domain    │  │  Role Domain    │
│                 │  │                 │
│ ┌─────────────┐ │  │ ┌─────────────┐ │
│ │UserService  │ │  │ │RoleService  │ │
│ │             │ │  │ │             │ │
│ │CreateUser() │ │  │ │CreateRole() │ │
│ │GetUserByID()│ │  │ │GetRoleByID()│ │
│ └─────────────┘ │  │ └─────────────┘ │
└─────────────────┘  └─────────────────┘
```

## Code Templates

### Creating an Application Service

```go
// internal/application/service/your_service.go
package service

import (
    "github.com/FeisalDy/nogo/internal/database"
    domain1Repo "github.com/FeisalDy/nogo/internal/domain1/repository"
    domain2Repo "github.com/FeisalDy/nogo/internal/domain2/repository"
    "gorm.io/gorm"
)

type YourService struct {
    domain1Repo *domain1Repo.Repository
    domain2Repo *domain2Repo.Repository
    db          *gorm.DB
}

func NewYourService(
    d1Repo *domain1Repo.Repository,
    d2Repo *domain2Repo.Repository,
    db *gorm.DB,
) *YourService {
    return &YourService{
        domain1Repo: d1Repo,
        domain2Repo: d2Repo,
        db:          db,
    }
}

func (s *YourService) YourCrossDomainOperation() error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        // 1. Domain 1 operation
        result1, err := s.domain1Repo.WithTx(tx).SomeMethod()
        if err != nil {
            return err
        }

        // 2. Domain 2 operation
        result2, err := s.domain2Repo.WithTx(tx).AnotherMethod()
        if err != nil {
            return err
        }

        // 3. Cross-domain coordination
        // ...

        return nil // Commit
    })
}
```

### Creating an Application Handler

```go
// internal/application/handler/your_handler.go
package handler

import (
    "net/http"
    "github.com/FeisalDy/nogo/internal/application/dto"
    "github.com/FeisalDy/nogo/internal/application/service"
    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type YourHandler struct {
    yourService *service.YourService
    validator   *validator.Validate
}

func NewYourHandler(yourService *service.YourService) *YourHandler {
    return &YourHandler{
        yourService: yourService,
        validator:   validator.New(),
    }
}

func (h *YourHandler) HandleOperation(c *gin.Context) {
    var req dto.YourRequestDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
        return
    }

    if err := h.validator.Struct(req); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
        return
    }

    result, err := h.yourService.YourCrossDomainOperation()
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, result, "Operation successful")
}
```

### Registering Routes

```go
// internal/application/routes.go
package application

import (
    "github.com/FeisalDy/nogo/internal/application/handler"
    "github.com/FeisalDy/nogo/internal/application/service"
    "github.com/FeisalDy/nogo/internal/common/middleware"
    domain1Repo "github.com/FeisalDy/nogo/internal/domain1/repository"
    domain2Repo "github.com/FeisalDy/nogo/internal/domain2/repository"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.RouterGroup) {
    // Initialize repositories
    d1Repo := domain1Repo.NewRepository(db)
    d2Repo := domain2Repo.NewRepository(db)

    // Initialize service
    yourService := service.NewYourService(d1Repo, d2Repo, db)
    
    // Initialize handler
    yourHandler := handler.NewYourHandler(yourService)

    // Register routes
    yourRoutes := router.Group("/your-resource")
    yourRoutes.Use(middleware.AuthMiddleware())
    {
        yourRoutes.POST("/action", yourHandler.HandleOperation)
    }
}
```

## Current Application Layer Services

### 1. AuthService
**Purpose:** Handle user registration with role assignment

**Endpoint:** `POST /api/v1/auth/register`

**Use Case:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secret123",
    "confirm_password": "secret123"
  }'
```

### 2. UserRoleService
**Purpose:** Manage user-role relationships

**Endpoints:**
- `POST /api/v1/user-roles/assign` - Assign role to user
- `POST /api/v1/user-roles/remove` - Remove role from user
- `GET /api/v1/user-roles/users/:user_id/roles` - Get user's roles

**Use Cases:**
```bash
# Assign role
curl -X POST http://localhost:8080/api/v1/user-roles/assign \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "role_id": 2}'

# Remove role
curl -X POST http://localhost:8080/api/v1/user-roles/remove \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "role_id": 2}'

# Get user roles
curl -X GET http://localhost:8080/api/v1/user-roles/users/1/roles \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Testing Guidelines

### Testing Domain Services
```go
func TestUserService_CreateUser(t *testing.T) {
    // Mock only user repository
    mockUserRepo := &MockUserRepository{}
    
    userService := service.NewUserService(mockUserRepo, mockCasbin)
    
    // Test user creation (single domain)
    user, err := userService.CreateUser(dto)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Testing Application Services
```go
func TestAuthService_Register(t *testing.T) {
    // Mock both user and role repositories
    mockUserRepo := &MockUserRepository{}
    mockRoleRepo := &MockRoleRepository{}
    mockDB := setupTestDB()
    
    authService := service.NewAuthService(
        mockUserRepo, 
        mockRoleRepo, 
        mockCasbin,
        mockDB,
    )
    
    // Test cross-domain operation
    user, err := authService.Register(dto)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    // Verify role was assigned
    mockUserRepo.AssertCalled(t, "AssignRoleToUser")
}
```

## Common Patterns

### Pattern 1: Validation Before Action
```go
func (s *YourService) Operation(id1, id2 uint) error {
    // 1. Validate all entities exist
    entity1, err := s.repo1.GetByID(id1)
    if err != nil || entity1 == nil {
        return errors.ErrNotFound
    }
    
    entity2, err := s.repo2.GetByID(id2)
    if err != nil || entity2 == nil {
        return errors.ErrNotFound
    }
    
    // 2. Perform operation
    return s.performAction(entity1, entity2)
}
```

### Pattern 2: Transaction Wrapping
```go
func (s *YourService) Operation() error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        // All operations within transaction
        if err := s.step1(tx); err != nil {
            return err // Rollback
        }
        
        if err := s.step2(tx); err != nil {
            return err // Rollback
        }
        
        return nil // Commit
    })
}
```

### Pattern 3: Event Coordination
```go
func (s *YourService) Operation() error {
    // 1. Domain operation
    result, err := s.domainRepo.DoSomething()
    if err != nil {
        return err
    }
    
    // 2. Sync with external systems (Casbin, cache, etc.)
    if err := s.externalService.Sync(result); err != nil {
        return err
    }
    
    return nil
}
```

## Troubleshooting

### Issue: "Circular dependency detected"
**Solution:** Move the operation to application layer

### Issue: "Transaction not working across domains"
**Solution:** Use `WithTx(tx)` for all repository calls within transaction

### Issue: "Domain service has dependency on another domain"
**Solution:** Remove the dependency and create application service

## Best Practices

1. **Keep domains independent** - Don't import domain services into other domains
2. **Use transactions** - Wrap cross-domain operations in transactions
3. **Validate early** - Check all entities exist before performing operations
4. **Fail fast** - Return errors immediately, don't continue with invalid state
5. **Document cross-domain flows** - Comment why operations span domains
6. **Test with mocks** - Mock dependencies for unit testing
7. **Use DTOs** - Don't expose internal models across layer boundaries

## Related Documentation

- [Application Layer Architecture](APPLICATION_LAYER.md) - Detailed explanation
- [Architecture Overview](ARCHITECTURE.md) - Overall architecture
- [Cross-Domain Relationships](CROSS_DOMAIN_RELATIONSHIPS.md) - Domain interaction patterns
