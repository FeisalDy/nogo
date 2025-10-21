# Application Layer Architecture Diagram

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          HTTP Client                                 │
│                    (Browser, Mobile App, etc.)                       │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             │ HTTP Requests
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Router Layer                                 │
│                    (internal/router/router.go)                       │
│                                                                      │
│  Routes:                                                            │
│  • /api/v1/auth/*          → Application Layer                     │
│  • /api/v1/user-roles/*    → Application Layer                     │
│  • /api/v1/users/*         → User Domain                           │
│  • /api/v1/roles/*         → Role Domain                           │
└────────────────────────────┬────────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
┌──────────────────┐  ┌─────────────────────┐  ┌──────────────────┐
│  User Domain     │  │ Application Layer   │  │  Role Domain     │
│    Handler       │  │     Handler         │  │    Handler       │
└────────┬─────────┘  └──────────┬──────────┘  └────────┬─────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌──────────────────┐  ┌─────────────────────┐  ┌──────────────────┐
│  User Domain     │  │ Application Layer   │  │  Role Domain     │
│    Service       │  │     Service         │  │    Service       │
│                  │  │                     │  │                  │
│ • CreateUser()   │  │ • Register()        │  │ • CreateRole()   │
│ • Login()        │  │ • AssignRole()      │  │ • GetRoleByID()  │
│ • GetUserByID()  │  │ • RemoveRole()      │  │ • UpdateRole()   │
│                  │  │ • GetUserRoles()    │  │ • DeleteRole()   │
└────────┬─────────┘  └──────────┬──────────┘  └────────┬─────────┘
         │                       │                       │
         │                       │                       │
         │            ┌──────────┴──────────┐            │
         │            │                     │            │
         ▼            ▼                     ▼            ▼
┌──────────────────────────────────────────────────────────────┐
│                    Repository Layer                           │
│                                                               │
│  ┌─────────────────┐         ┌─────────────────┐           │
│  │ User Repository │         │ Role Repository │           │
│  │                 │         │                 │           │
│  │ • GetUserByID() │         │ • GetByID()     │           │
│  │ • CreateUser()  │         │ • Create()      │           │
│  │ • GetByEmail()  │         │ • Update()      │           │
│  │ • AssignRole()  │         │ • GetByName()   │           │
│  └─────────────────┘         └─────────────────┘           │
└────────────────────────┬──────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Database Layer                               │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │  users       │  │  roles       │  │  user_roles  │            │
│  │  table       │  │  table       │  │  table       │            │
│  └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                      │
│  ┌──────────────────────────────────────────────┐                  │
│  │  casbin_rule (authorization policies)        │                  │
│  └──────────────────────────────────────────────┘                  │
└─────────────────────────────────────────────────────────────────────┘
```

## Request Flow Examples

### Example 1: User Registration (Cross-Domain)

```
1. Client                    POST /api/v1/auth/register
   │
   ▼
2. Router                    → application.RegisterRoutes()
   │
   ▼
3. AuthHandler              → Register(ctx)
   │
   ▼
4. AuthService              → Register(dto)
   │
   ├─────────────────────────────────────────┐
   │                                         │
   ▼                                         ▼
5. UserRepository                      RoleRepository
   │ CreateUser()                        │ GetByName("user")
   │                                     │
   │                                     │
6. ├─── Database ──────────────────────┤
   │     (Transaction)                  │
   │                                    │
   ▼                                    │
7. UserRepository                       │
   │ AssignRoleToUser()                 │
   │                                    │
8. └────────────────────────────────────┘
   │
   ▼
9. CasbinService
   │ AssignRoleToUser()
   │
   ▼
10. Response with JWT Token
```

### Example 2: Assign Role to User (Cross-Domain)

```
1. Client                    POST /api/v1/user-roles/assign
   │                         {user_id: 1, role_id: 2}
   ▼
2. Router                    → application.RegisterRoutes()
   │
   ▼
3. UserRoleHandler          → AssignRoleToUser(ctx)
   │
   ▼
4. UserRoleService          → AssignRoleToUser(userID, roleID)
   │
   ├─────────────────────────────────────────┐
   │                                         │
   ▼                                         ▼
5. UserRepository                      RoleRepository
   │ GetUserByID()                      │ GetByID()
   │ (validate exists)                  │ (validate exists)
   │                                    │
   │                                    │
   ▼                                    │
6. UserRepository                       │
   │ HasRoleByID()                      │
   │ (check not already assigned)       │
   │                                    │
   ▼                                    │
7. UserRepository                       │
   │ AssignRoleToUser()                 │
   │                                    │
8. └────────────────────────────────────┘
   │
   ▼
9. CasbinService
   │ AssignRoleToUser()
   │
   ▼
10. Success Response
```

### Example 3: Get User by ID (Single Domain)

```
1. Client                    GET /api/v1/users/:id
   │
   ▼
2. Router                    → user.RegisterRoutes()
   │
   ▼
3. UserHandler              → GetUserByID(ctx)
   │
   ▼
4. UserService              → GetUserByID(id)
   │
   ▼
5. UserRepository           → GetUserByID(id)
   │
   ▼
6. Database                  SELECT * FROM users WHERE id = ?
   │
   ▼
7. User Model
   │
   ▼
8. Response with User Data
```

## Layer Responsibilities

### 1. Handler Layer (HTTP)
**Responsibility:** Handle HTTP requests and responses

```
┌─────────────────────────────────────┐
│          Handler Layer              │
│                                     │
│  • Parse request body               │
│  • Validate input                   │
│  • Call service layer               │
│  • Format response                  │
│  • Handle HTTP errors               │
│                                     │
└─────────────────────────────────────┘
```

### 2. Service Layer
**Responsibility:** Business logic

#### Domain Services (Single Domain)
```
┌─────────────────────────────────────┐
│       Domain Service Layer          │
│                                     │
│  • Single domain operations         │
│  • Domain-specific validation       │
│  • Business rules for domain        │
│  • No cross-domain dependencies     │
│                                     │
└─────────────────────────────────────┘
```

#### Application Services (Cross-Domain)
```
┌─────────────────────────────────────┐
│     Application Service Layer       │
│                                     │
│  • Coordinate multiple domains      │
│  • Transaction management           │
│  • Cross-domain validation          │
│  • Orchestrate complex workflows    │
│  • Maintain consistency             │
│                                     │
└─────────────────────────────────────┘
```

### 3. Repository Layer
**Responsibility:** Data access

```
┌─────────────────────────────────────┐
│        Repository Layer             │
│                                     │
│  • Database queries                 │
│  • CRUD operations                  │
│  • Transaction support              │
│  • Data mapping                     │
│                                     │
└─────────────────────────────────────┘
```

## Dependency Flow

```
Outer Layer → Inner Layer (Dependencies point inward)

┌────────────────────────────────────────────────┐
│                  Handlers                      │  ← HTTP Layer
│  (Depends on: Services)                        │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│              Application Services               │  ← Orchestration Layer
│  (Depends on: Domain Services, Repositories)   │
└──────────────────┬─────────────────────────────┘
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
┌──────────────────┐  ┌──────────────────┐
│ Domain Services  │  │   Repositories   │        ← Domain Layer
│                  │  │                  │
└──────────────────┘  └────────┬─────────┘
                               │
                               ▼
                    ┌──────────────────┐
                    │    Database      │           ← Data Layer
                    └──────────────────┘
```

## Key Principles

### 1. Separation of Concerns
- Each layer has a specific responsibility
- Layers don't skip levels (no handler → repository directly)

### 2. Dependency Rule
- Dependencies point inward
- Inner layers don't know about outer layers
- Application layer coordinates domains

### 3. Transaction Boundaries
- Application layer manages transactions
- Ensures atomicity across domains

### 4. Single Responsibility
- Domain services: single domain operations
- Application services: cross-domain coordination

## Benefits Visualization

### Before: Tangled Dependencies

```
┌─────────────┐     ┌─────────────┐
│    User     │────▶│    Role     │
│   Service   │     │   Service   │
│             │◀────│             │
└─────────────┘     └─────────────┘
      ▲                   ▲
      │                   │
      └───────┬───────────┘
              │
        Circular Dependencies
        Hard to Maintain
```

### After: Clean Architecture

```
        ┌───────────────────┐
        │   Application     │
        │     Layer         │
        │   (Coordinator)   │
        └─────────┬─────────┘
                  │
         ┌────────┴────────┐
         ▼                 ▼
┌─────────────┐   ┌─────────────┐
│    User     │   │    Role     │
│   Service   │   │   Service   │
│             │   │             │
└─────────────┘   └─────────────┘

     Independent Domains
     Easy to Maintain
```

---

**Summary:** The application layer sits between handlers and domain services, coordinating operations that span multiple domains while keeping the domains themselves independent and focused.
