# Complete DDD Compliance: Removing Authorization from Domain Layer

**Date:** October 20, 2025  
**Status:** ✅ Complete

## Problem: Authorization in Domain Layer

### Issue Identified

The **User domain** was breaking DDD principles by:

1. ❌ **Depending on Casbin** (authorization infrastructure)
2. ❌ **Handling cross-domain logic** (User + Role + Permissions)
3. ❌ **Mixing domain logic with infrastructure concerns**

### Problematic Code

```go
// ❌ BAD: User domain service depends on Casbin
type UserService struct {
    userRepo      *repository.UserRepository
    casbinService *casbinService.CasbinService  // ❌ Infrastructure dependency
}

// ❌ BAD: Cross-domain operation in domain service
func (s *UserService) GetUserWithPermissions(userID uint) (*dto.UserWithPermissionsDTO, error) {
    user, _ := s.userRepo.GetUserByID(userID)           // User domain
    roleNames, _ := s.casbinService.GetRolesForUser()   // ❌ Authorization infrastructure
    permissions, _ := s.casbinService.GetPermissions()  // ❌ Authorization infrastructure
    // ... returns combined data
}
```

## Why This Breaks DDD

### 1. **Domain Service Should Be Pure Business Logic**
- User domain should only handle user-related business operations
- Authorization is infrastructure/cross-cutting concern
- Mixing them creates tight coupling

### 2. **Violates Single Responsibility Principle**
```
User Domain Responsibilities:
✅ User CRUD operations
✅ User authentication (login)
✅ User profile management
❌ Authorization (roles/permissions) ← Should be elsewhere
```

### 3. **Creates Hidden Cross-Domain Dependencies**
```
User Service
    ├─> User Repository ✅ (Same domain)
    ├─> Casbin Service ❌ (Infrastructure)
    └─> Implicitly depends on Role domain ❌
```

### 4. **Prevents Domain Independence**
- Can't test User domain without Casbin
- Can't replace authorization system without touching User domain
- Can't extract User domain to microservice easily

## Solution: Move to Application Layer

### Architecture Before

```
┌────────────────────────────────────┐
│      User Domain (Handler)         │
│              ↓                      │
│      User Domain (Service)         │
│              ↓                      │
│         ┌────┴────┐                │
│         ↓         ↓                 │
│    User Repo   Casbin ❌           │  Violates domain boundaries
│                                     │
└────────────────────────────────────┘
```

### Architecture After

```
┌─────────────────────────────────────────────┐
│       Application Layer                     │
│                                             │
│   UserProfileService (Coordinator)         │
│              ↓                              │
│      ┌───────┼────────┐                    │
│      ↓       ↓        ↓                     │
│   User    Role    Casbin                   │  ✅ Coordinates all
│   Repo    Repo    Service                  │
│                                             │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│       User Domain (Pure)                    │
│                                             │
│   UserService                               │
│       ↓                                     │
│   User Repo  ✅ Only user domain           │
│                                             │
└─────────────────────────────────────────────┘
```

## Changes Made

### 1. Cleaned User Domain Service

**Before:**
```go
type UserService struct {
    userRepo      *repository.UserRepository
    casbinService *casbinService.CasbinService  // ❌
}

func NewUserService(userRepository, casbin) *UserService { ... }
func (s *UserService) GetUserWithPermissions() { ... }  // ❌
```

**After:**
```go
type UserService struct {
    userRepo *repository.UserRepository  // ✅ Pure domain
}

func NewUserService(userRepository) *UserService { ... }
// GetUserWithPermissions removed - moved to application layer
```

### 2. Created Application Layer Service

**New: `UserProfileService` (Application Layer)**
```go
type UserProfileService struct {
    userRepo      *userRepo.UserRepository      // User domain
    roleRepo      *roleRepo.RoleRepository      // Role domain
    casbinService *casbinService.CasbinService  // Auth infrastructure
    db            *gorm.DB
}

func (s *UserProfileService) GetUserWithPermissions(userID uint) (*dto.UserWithPermissionsDTO, error) {
    // 1. Get user from User domain
    user := s.userRepo.GetUserByID(userID)
    
    // 2. Get role IDs from user_roles
    roleIDs := s.userRepo.GetUserRoleIDs(userID)
    
    // 3. Get role details from Role domain
    for _, roleID := range roleIDs {
        role := s.roleRepo.GetByID(roleID)
        // ...
    }
    
    // 4. Get permissions from Casbin
    permissions := s.casbinService.GetPermissionsForRole(roleName)
    
    // 5. Combine everything and return
    return response
}
```

### 3. Created Application Layer Handler

**New: `UserProfileHandler`**
```go
type UserProfileHandler struct {
    userProfileService *service.UserProfileService
}

func (h *UserProfileHandler) GetMe(c *gin.Context) {
    userID := c.Get("user_id")
    profile, _ := h.userProfileService.GetUserWithPermissions(userID)
    utils.RespondSuccess(c, http.StatusOK, profile)
}
```

### 4. Updated Routes

**Before:**
```go
// ❌ Cross-domain endpoint in User domain
GET /api/v1/users/me  → UserHandler.GetMe → UserService.GetUserWithPermissions
```

**After:**
```go
// ✅ Cross-domain endpoint in Application layer
GET /api/v1/profile/me  → UserProfileHandler.GetMe → UserProfileService.GetUserWithPermissions

// ✅ Pure domain endpoint remains in User domain
GET /api/v1/users/:email  → UserHandler.GetUserByEmail → UserService.GetUserByEmail
```

## API Changes

### Breaking Change

| Before | After | Status |
|--------|-------|--------|
| `GET /api/v1/users/me` | `GET /api/v1/profile/me` | ⚠️ **Moved** |

### Migration Guide for Clients

**Old endpoint:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer TOKEN"
```

**New endpoint:**
```bash
curl -X GET http://localhost:8080/api/v1/profile/me \
  -H "Authorization: Bearer TOKEN"
```

**Response format:** Unchanged ✅
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "john_doe",
    "roles": [
      {"id": 1, "name": "user"},
      {"id": 2, "name": "admin"}
    ],
    "permissions": [
      {"resource": "users", "action": "read"},
      {"resource": "users", "action": "write"}
    ]
  }
}
```

## Files Changed

### Created
1. ✅ `internal/application/service/user_profile_service.go`
2. ✅ `internal/application/handler/user_profile_handler.go`

### Modified
1. ✅ `internal/user/service/user_service.go`
   - Removed `casbinService` dependency
   - Removed `GetUserWithPermissions()` method
   - Now pure User domain service

2. ✅ `internal/user/handler/user_handler.go`
   - Removed `GetMe()` handler
   - Added comment explaining migration

3. ✅ `internal/user/routes.go`
   - Removed `casbinService` import
   - Removed `/me` route
   - Updated comments

4. ✅ `internal/application/routes.go`
   - Added `UserProfileService` initialization
   - Added `UserProfileHandler` initialization
   - Added `/profile/me` route

## Domain Boundaries Now Properly Defined

### User Domain (Pure)
```
Responsibilities:
✅ User CRUD operations
✅ User authentication (login)
✅ Password management
✅ User profile data

Dependencies:
✅ User Repository ONLY
❌ No Casbin
❌ No Role Repository
```

### Role Domain (Pure)
```
Responsibilities:
✅ Role CRUD operations
✅ Role management

Dependencies:
✅ Role Repository ONLY
❌ No User Repository
❌ No Casbin
```

### Application Layer (Coordinator)
```
Responsibilities:
✅ Cross-domain orchestration
✅ User + Role + Permission coordination
✅ Complex workflows
✅ Transaction management

Dependencies:
✅ User Repository
✅ Role Repository
✅ Casbin Service
✅ All domains it coordinates
```

## Benefits Achieved

### 1. **True Domain Independence**
```go
// ✅ Can test User domain without Casbin
func TestUserService_GetUserByID(t *testing.T) {
    mockUserRepo := &MockUserRepository{}
    // NO need to mock Casbin anymore!
    
    userService := service.NewUserService(mockUserRepo)
    user, err := userService.GetUserByID(1)
    
    assert.NoError(t, err)
}
```

### 2. **Clear Separation of Concerns**
- User domain: Business logic for users
- Role domain: Business logic for roles
- Application layer: Coordination and workflows
- Infrastructure (Casbin): Authorization policies

### 3. **Flexibility**
- Can replace Casbin with another authorization system
- Changes only affect Application layer
- Domains remain untouched

### 4. **Testability**
- Each domain can be unit tested independently
- Application layer tests coordination logic
- Clear boundaries for integration tests

### 5. **Maintainability**
- Easy to understand what each layer does
- Clear file organization
- Explicit dependencies

## Testing Strategy

### Unit Tests

**User Domain (Pure)**
```go
func TestUserService_Login(t *testing.T) {
    mockUserRepo := &MockUserRepository{}
    userService := service.NewUserService(mockUserRepo)
    
    // Test only user domain logic
    user, err := userService.Login(loginDTO)
    assert.NoError(t, err)
}
```

**Application Layer (Coordination)**
```go
func TestUserProfileService_GetUserWithPermissions(t *testing.T) {
    mockUserRepo := &MockUserRepository{}
    mockRoleRepo := &MockRoleRepository{}
    mockCasbin := &MockCasbinService{}
    
    service := NewUserProfileService(mockUserRepo, mockRoleRepo, mockCasbin, db)
    
    // Test coordination logic
    profile, err := service.GetUserWithPermissions(1)
    
    assert.NoError(t, err)
    assert.NotNil(t, profile.Roles)
    assert.NotNil(t, profile.Permissions)
}
```

### Integration Tests
```bash
# Test the full flow
curl -X POST http://localhost:8080/api/v1/users/login \
  -d '{"email":"user@example.com","password":"pass123"}'

TOKEN=$(extract_token_from_response)

curl -X GET http://localhost:8080/api/v1/profile/me \
  -H "Authorization: Bearer $TOKEN"
```

## DDD Principles Now Followed

### ✅ 1. Domain Independence
- Each domain only knows about itself
- No cross-domain dependencies in domain services

### ✅ 2. Bounded Contexts
- User domain: User aggregate and operations
- Role domain: Role aggregate and operations
- Application layer: Cross-domain workflows

### ✅ 3. Layered Architecture
```
Presentation Layer (Handlers)
         ↓
Application Layer (Coordinators)
         ↓
Domain Layer (Business Logic)
         ↓
Infrastructure Layer (Repositories, Casbin)
```

### ✅ 4. Dependency Direction
```
Outer layers depend on inner layers
Never the reverse

Handlers → Application Services → Domain Services → Repositories
```

### ✅ 5. Single Responsibility
- Domain services: Domain-specific business logic
- Application services: Cross-domain coordination
- Repositories: Data access
- Handlers: HTTP concerns

## Comparison Summary

| Aspect | Before | After |
|--------|--------|-------|
| **User Domain Dependencies** | User Repo + Casbin ❌ | User Repo only ✅ |
| **Cross-Domain Logic Location** | User Service ❌ | Application Layer ✅ |
| **Authorization Coupling** | Tight ❌ | Loose ✅ |
| **Testability** | Difficult ❌ | Easy ✅ |
| **DDD Compliance** | No ❌ | Yes ✅ |
| **API Endpoint** | `/users/me` | `/profile/me` |
| **Maintainability** | Complex ❌ | Clear ✅ |

## Conclusion

✅ **Full DDD Compliance Achieved!**

The application now properly separates concerns:

1. **Domain Layer** = Pure business logic for specific domains
2. **Application Layer** = Cross-domain coordination
3. **Infrastructure** = Technical concerns (DB, auth, etc.)

All cross-domain operations have been moved to the application layer, making the codebase:
- More maintainable
- Easier to test
- Better aligned with DDD principles
- Ready for future scaling (microservices, etc.)

---

**Related Documentation:**
- [Application Layer Architecture](APPLICATION_LAYER.md)
- [DDD Repository Compliance](DDD_COMPLIANCE_FIX.md)
- [Refactoring Summary](REFACTORING_SUMMARY.md)
