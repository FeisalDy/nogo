# DDD Compliance Fix: Repository Layer

## Problem Identified

The `UserRepository` was **breaking Domain-Driven Design (DDD) principles** by:

1. ❌ Importing `roleModel` from the Role domain
2. ❌ Directly querying the `roles` table
3. ❌ Returning `Role` entities from another domain
4. ❌ Creating tight coupling between User and Role domains at the data layer

### Problematic Code (Before)

```go
// ❌ BAD: Importing role domain model
import (
    roleModel "github.com/FeisalDy/nogo/internal/role/model"
)

// ❌ BAD: Repository method returns entities from another domain
func (r *UserRepository) GetUserWithRoles(userID uint) (*model.User, []roleModel.Role, error) {
    var roles []roleModel.Role
    err := r.db.
        Table("roles").  // ❌ Directly querying role table
        Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
        Where("user_roles.user_id = ?", userID).
        Find(&roles).Error
    
    return &user, roles, err  // ❌ Returning role entities
}

// ❌ BAD: Joining with roles table to check by name
func (r *UserRepository) HasRole(userID uint, roleName string) (bool, error) {
    var count int64
    err := r.db.
        Table("user_roles").
        Joins("INNER JOIN roles ON roles.id = user_roles.role_id").  // ❌ Cross-domain join
        Where("user_roles.user_id = ? AND roles.name = ?", userID, roleName).
        Count(&count).Error
    return count > 0, err
}
```

## DDD Principles Violated

### 1. **Bounded Context Violation**
- User repository was crossing into Role domain's bounded context
- Direct access to Role table from User repository

### 2. **Dependency Direction**
- Repository layer was creating dependencies between domains
- Violated the principle that domains should be independent

### 3. **Encapsulation Breach**
- User domain was aware of Role domain's internal structure
- Role table schema changes would impact User repository

## Solution Implemented

### Fixed Code (After)

```go
// ✅ GOOD: No role domain imports
import (
    commonModel "github.com/FeisalDy/nogo/internal/common/model"
    "github.com/FeisalDy/nogo/internal/user/model"
)

// ✅ GOOD: Only returns IDs, not entities from another domain
func (r *UserRepository) GetUserRoleIDs(userID uint) ([]uint, error) {
    var roleIDs []uint
    err := r.db.
        Model(&commonModel.UserRole{}).  // ✅ Only touches junction table
        Where("user_id = ?", userID).
        Pluck("role_id", &roleIDs).Error
    return roleIDs, err
}

// ✅ GOOD: Only checks junction table, no cross-domain joins
func (r *UserRepository) HasRoleByID(userID, roleID uint) (bool, error) {
    var count int64
    err := r.db.
        Table("user_roles").  // ✅ Only junction table
        Where("user_id = ? AND role_id = ?", userID, roleID).
        Count(&count).Error
    return count > 0, err
}
```

## Key Changes

### 1. Removed Cross-Domain Methods
- ❌ Removed: `GetUserWithRoles()` - was returning Role entities
- ❌ Removed: `HasRole(userID, roleName)` - was joining with roles table
- ❌ Removed: `HasAnyRole(userID, roleNames)` - was joining with roles table

### 2. Repository Now Only Handles User-Roles Junction
- ✅ `GetUserRoleIDs()` - Returns IDs only, not entities
- ✅ `HasRoleByID()` - Uses IDs, not names (no join needed)
- ✅ `AssignRoleToUser()` - Only touches `user_roles` table
- ✅ `RemoveRoleFromUser()` - Only touches `user_roles` table

### 3. Application Layer Coordinates Domains

**Before (Breaking DDD):**
```go
// ❌ User repository doing cross-domain work
user, roles, err := userRepo.GetUserWithRoles(userID)
```

**After (DDD Compliant):**
```go
// ✅ Application layer coordinates both domains
roleIDs, err := userRepo.GetUserRoleIDs(userID)  // User domain
for _, roleID := range roleIDs {
    role, err := roleRepo.GetByID(roleID)  // Role domain
    // Use role...
}
```

## Benefits of This Fix

### 1. **True Domain Independence**
```
Before:
User Domain ──depends on──> Role Domain  ❌

After:
User Domain                Role Domain
     │                          │
     └────> Application <───────┘         ✅
            (Coordinator)
```

### 2. **Flexibility**
- Can change Role table schema without affecting User repository
- Can replace Role implementation without touching User code
- Easier to test each domain in isolation

### 3. **Maintainability**
- Clear boundaries between domains
- Each repository only knows about its own tables
- Cross-domain logic is explicit in application layer

### 4. **Scalability**
- Domains can be extracted to microservices
- No hidden dependencies between domains
- Clear integration points

## Updated Architecture

### Repository Layer Boundaries

```
┌─────────────────────────────────────────────────────────┐
│              Application Layer                          │
│  (Coordinates multiple domains)                         │
│                                                          │
│  roleIDs := userRepo.GetUserRoleIDs(userID)            │
│  for _, roleID := range roleIDs {                       │
│      role := roleRepo.GetByID(roleID)  ← Coordination │
│  }                                                       │
└──────────────┬─────────────────────┬────────────────────┘
               │                     │
               ▼                     ▼
┌───────────────────────┐  ┌───────────────────────┐
│  User Repository      │  │  Role Repository      │
│                       │  │                       │
│  ✅ users table       │  │  ✅ roles table       │
│  ✅ user_roles table  │  │  ✅ user_roles table  │
│     (IDs only)        │  │     (IDs only)        │
│                       │  │                       │
│  ❌ NO roles table    │  │  ❌ NO users table    │
└───────────────────────┘  └───────────────────────┘
```

### Data Flow for Getting User Roles

```
1. Client Request
   │
   ▼
2. Application Layer (UserRoleService.GetUserRoles)
   │
   ├──> 3a. userRepo.GetUserByID(userID)      [User Domain]
   │    └─> Returns: User entity
   │
   └──> 3b. userRepo.GetUserRoleIDs(userID)   [User Domain]
        └─> Returns: [1, 2, 3]  (just IDs)
        │
        └──> 4. For each roleID:
             ├─> roleRepo.GetByID(1)           [Role Domain]
             ├─> roleRepo.GetByID(2)           [Role Domain]
             └─> roleRepo.GetByID(3)           [Role Domain]
             └─> Returns: Role entities
   │
   ▼
5. Combine results in Application Layer
   │
   ▼
6. Return to Client
```

## Migration Impact

### Code Changed
- ✅ `internal/user/repository/user_repository.go` - Removed cross-domain methods
- ✅ `internal/application/service/user_role_service.go` - Updated to coordinate domains
- ✅ `internal/user/service/user_service.go` - Updated to use Casbin for role info

### Breaking Changes
**None!** All changes are internal refactoring. External API remains the same.

### Performance Impact
**Minimal.** 
- Before: 1 query with JOIN
- After: 1 query for IDs + N queries for roles (where N = number of roles per user)
- Typical users have 1-3 roles, so 2-4 queries total
- Can be optimized with batch fetching if needed

## Testing Recommendations

### Unit Tests
```go
func TestUserRepository_GetUserRoleIDs(t *testing.T) {
    // Test that it only returns IDs, not full role objects
    roleIDs, err := userRepo.GetUserRoleIDs(userID)
    
    assert.NoError(t, err)
    assert.IsType(t, []uint{}, roleIDs)  // Just IDs
}

func TestUserRoleService_GetUserRoles(t *testing.T) {
    // Test that application layer properly coordinates
    mockUserRepo := &MockUserRepository{}
    mockRoleRepo := &MockRoleRepository{}
    
    mockUserRepo.On("GetUserRoleIDs", userID).Return([]uint{1, 2}, nil)
    mockRoleRepo.On("GetByID", uint(1)).Return(role1, nil)
    mockRoleRepo.On("GetByID", uint(2)).Return(role2, nil)
    
    roles, err := service.GetUserRoles(userID)
    
    assert.NoError(t, err)
    assert.Len(t, roles, 2)
}
```

### Integration Tests
```bash
# Test getting user roles
curl -X GET http://localhost:8080/api/v1/user-roles/users/1/roles \
  -H "Authorization: Bearer TOKEN"

# Should return full role details even though repository only returns IDs
```

## Key Takeaways

### ✅ DDD Rules to Follow

1. **Repositories only query their own domain's tables**
   - User repository → users, user_roles (junction)
   - Role repository → roles, user_roles (junction)

2. **Cross-domain queries go through application layer**
   - Get IDs from one domain
   - Fetch entities from other domain
   - Coordinate in application layer

3. **Junction tables belong to common domain**
   - Both repositories can access junction table
   - But only work with IDs, not entities

4. **Never import domain models across domains**
   - No `roleModel` in User repository
   - No `userModel` in Role repository

### 🚫 Anti-Patterns to Avoid

1. **Direct cross-domain table access**
   ```go
   // ❌ BAD
   db.Table("roles").Joins("user_roles").Where(...)
   ```

2. **Returning entities from other domains**
   ```go
   // ❌ BAD
   func GetUserWithRoles() (*User, []Role, error)
   ```

3. **Cross-domain joins in repositories**
   ```go
   // ❌ BAD
   Joins("INNER JOIN roles ON roles.id = ...")
   ```

## Conclusion

✅ **DDD Compliance Achieved!**

The repository layer now properly respects domain boundaries:
- Each repository only accesses its own domain's tables
- Cross-domain operations are coordinated at the application layer
- Domains are truly independent and loosely coupled

This makes the codebase more maintainable, testable, and ready for future scaling needs like microservices.

---

**Related Documentation:**
- [Application Layer Architecture](APPLICATION_LAYER.md)
- [Cross-Domain Relationships](CROSS_DOMAIN_RELATIONSHIPS.md)
- [Architecture Diagram](ARCHITECTURE_DIAGRAM.md)
