# 📚 NoGo - Novel Platform Documentation

Welcome to the NoGo documentation! This guide will help you understand, develop, and maintain the Novel Platform application.

## 📖 Table of Contents

### [01. Getting Started](#01-getting-started-1)

Quick setup and introduction to the project

### [02. Architecture](#02-architecture-1)

System design and structure

### [03. Authentication](#03-authentication-1)

JWT-based authentication system

### [04. Authorization](#04-authorization-1)

Casbin ABAC/RBAC permission system

### [05. Database](#05-database-1)

Database schema and migrations

### [06. Error Handling](#06-error-handling-1)

Comprehensive error management

### [07. API](#07-api-1)

REST API documentation

### [08. Development](#08-development-1)

Development guides and roadmap

---

## 01. Getting Started

Start here if you're new to the project!

### Quick Links

- **[⚡ Seed Permissions](01-getting-started/SEED_PERMISSIONS.md)** - **START HERE** if you see empty permissions!
- **[Casbin Quick Start](01-getting-started/CASBIN_QUICK_START.md)** - Get Casbin running in 5 minutes
- **[/me Endpoint Guide](01-getting-started/ME_ENDPOINT_SUMMARY.md)** - User permissions endpoint
- **[Casbin Implementation Summary](01-getting-started/CASBIN_IMPLEMENTATION_SUMMARY.md)** - Complete Casbin overview
- **[Cleanup Complete](01-getting-started/CLEANUP_COMPLETE.md)** - Recent permission table cleanup

### ⚠️ Common Issue: Empty Permissions?

If `/api/v1/users/me` returns empty permissions array, you need to seed permissions first!

**Quick Fix:**

```bash
go run scripts/seed_casbin.go
```

See **[SEED_PERMISSIONS.md](01-getting-started/SEED_PERMISSIONS.md)** for details.

### What is NoGo?

A novel/web fiction platform built with:

- **Go** (Gin framework)
- **PostgreSQL** (GORM ORM)
- **JWT** authentication
- **Casbin** authorization
- **RESTful API**

### First Steps

1. Clone the repository
2. Install dependencies: `go mod download`
3. Setup database (PostgreSQL)
4. Run migrations: `go run cmd/migrate/main.go`
5. **Seed Casbin permissions**: `go run scripts/seed_casbin.go` ← **IMPORTANT!**
6. Start server: `go run cmd/server/main.go`

---

## 02. Architecture

Understand the system design and structure.

### Documents

- **[Architecture Overview](02-architecture/ARCHITECTURE.md)** - System architecture and design patterns
- **[Cross-Domain Relationships](02-architecture/CROSS_DOMAIN_RELATIONSHIPS.md)** - How domains interact

### Key Concepts

- **Domain-Driven Design** - Code organized by business domains
- **Layered Architecture** - Handler → Service → Repository → Database
- **Dependency Injection** - Loose coupling between layers
- **Clean Code** - Each domain is self-contained

### Project Structure

```
nogo/
├── cmd/
│   ├── server/     # Main application
│   └── migrate/    # Database migrations
├── config/         # Configuration
│   └── casbin/     # Casbin model & policies
├── internal/
│   ├── common/     # Shared code (middleware, errors, casbin)
│   ├── database/   # Database & migrations
│   ├── user/       # User domain
│   ├── role/       # Role domain
│   ├── novel/      # Novel domain
│   └── router/     # Route setup
├── docs/           # Documentation (you are here)
└── scripts/        # Utility scripts
```

---

## 03. Authentication

JWT-based authentication system.

### Documents

- **[Authentication Guide](03-authentication/AUTHENTICATION.md)** - Complete authentication guide
- **[Auth Implementation Summary](03-authentication/AUTH_IMPLEMENTATION_SUMMARY.md)** - Implementation details
- **[Auth Quick Reference](03-authentication/AUTH_QUICK_REFERENCE.md)** - Quick lookup guide
- **[Auth Testing](03-authentication/AUTH_TESTING.md)** - Test authentication flows

### How It Works

1. **User Registration** → Hash password, create user
2. **User Login** → Validate credentials, generate JWT token
3. **Protected Routes** → Validate JWT via middleware
4. **Token Refresh** → Generate new token when expired

### Key Files

- `internal/common/utils/jwt.go` - JWT generation & validation
- `internal/common/middleware/auth.go` - Authentication middleware
- `internal/common/utils/password.go` - Password hashing

### Example

```go
// Login
POST /api/v1/users/login
{
  "email": "user@example.com",
  "password": "password123"
}

// Response
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}

// Use token in requests
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## 04. Authorization

Casbin-based ABAC/RBAC permission system.

### Documents

- **[Casbin ABAC Guide](04-authorization/CASBIN_ABAC_GUIDE.md)** 📘 - Complete Casbin guide
- **[Casbin Database Schema](04-authorization/CASBIN_DATABASE_SCHEMA.md)** - Database structure
- **[Casbin Route Examples](04-authorization/CASBIN_ROUTE_EXAMPLES.md)** - Protect your routes
- **[RBAC Implementation](04-authorization/RBAC_IMPLEMENTATION.md)** - Role-based access control

### How It Works

```
User → Has Roles → Roles Have Permissions → Access Granted/Denied
```

### Key Features

- ✅ Dynamic role creation
- ✅ Runtime permission changes
- ✅ Fine-grained access control
- ✅ Database-backed policies
- ✅ In-memory caching

### Quick Example

```go
// Protect a route
router.POST("/users",
    middleware.AuthMiddleware(),
    middleware.CasbinMiddleware("users", "write"),
    handler.CreateUser,
)

// Manage permissions
svc := casbinService.NewCasbinService()
svc.AddPermissionForRole("editor", "novels", "write")
svc.AssignRoleToUser(userID, "editor")
```

### Common Permissions

- `resource:read` - View
- `resource:write` - Create/Edit
- `resource:delete` - Delete
- `resource:publish` - Publish (custom)

---

## 05. Database

PostgreSQL database with GORM migrations.

### Documents

- **[Migration System](05-database/MIGRATION_SYSTEM.md)** - How migrations work
- **[Migration Quick Reference](05-database/MIGRATION_QUICK_REFERENCE.md)** - Quick commands

### Database Schema

```
users
  ├─ user_roles ──┐
  │               │
roles ────────────┘

casbin_rule (Casbin policies)

novels
  └─ chapters

genres
tags
novel_genres
novel_tags
```

### Running Migrations

```bash
# Run all migrations
go run cmd/migrate/main.go

# Check migration status
psql -U user -d database -c "SELECT * FROM migrations;"
```

### Creating Migrations

See `internal/database/migrations/` for examples.

---

## 06. Error Handling

Standardized error handling system.

### Documents

- **[Error Handling Guide](06-error-handling/ERROR_HANDLING.md)** - Complete guide
- **[Before & After](06-error-handling/ERROR_HANDLING_BEFORE_AFTER.md)** - Examples
- **[Error Flow](06-error-handling/ERROR_HANDLING_FLOW.md)** - Error flow diagram
- **[Quick Reference](06-error-handling/ERROR_HANDLING_QUICK_REFERENCE.md)** - Quick lookup
- **[Summary](06-error-handling/ERROR_HANDLING_SUMMARY.md)** - Overview
- **[Testing](06-error-handling/ERROR_HANDLING_TESTS.md)** - Test errors

### Error Structure

```go
{
  "error": {
    "code": "USER001",
    "message": "User not found",
    "details": {}
  }
}
```

### Error Codes

- `USER001-099` - User domain
- `ROLE001-099` - Role domain
- `AUTH001-099` - Authentication
- `VAL001-099` - Validation
- `GEN001-099` - General errors

### Usage

```go
// In handler
if user == nil {
    utils.RespondWithAppError(c, errors.ErrUserNotFound)
    return
}

// In service
if exists {
    return nil, errors.ErrUserAlreadyExists
}
```

---

## 07. API

REST API documentation.

### Documents

- **[API Documentation](07-api/API.md)** - Complete API reference

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoints

#### Authentication

- `POST /users/register` - Register new user
- `POST /users/login` - Login user

#### Users

- `GET /users` - List users
- `GET /users/:id` - Get user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

#### Roles

- `GET /roles` - List roles
- `POST /roles` - Create role
- `PUT /roles/:id` - Update role
- `DELETE /roles/:id` - Delete role
- `POST /roles/:id/users/:user_id` - Assign role

#### Novels

- `GET /novels` - List novels
- `POST /novels` - Create novel
- `GET /novels/:id` - Get novel
- `PUT /novels/:id` - Update novel
- `DELETE /novels/:id` - Delete novel

### Authentication Header

```
Authorization: Bearer <jwt_token>
```

---

## 08. Development

Development guides and project roadmap.

### Documents

- **[Development Guide](08-development/DEVELOPMENT.md)** - Development workflow
- **[Development Roadmap](08-development/DEVELOPMENT_ROADMAP.md)** - Future plans

### Setup Development Environment

```bash
# Install dependencies
go mod download

# Setup database
createdb nogo_dev

# Run migrations
go run cmd/migrate/main.go

# Seed permissions
go run scripts/seed_permissions.go

# Run server with hot reload (install air first)
air

# Or run normally
go run cmd/server/main.go
```

### Code Style

- Follow Go conventions
- Use meaningful names
- Write tests
- Document public functions
- Keep handlers thin (business logic in services)

### Adding a New Domain

1. Create folder in `internal/`
2. Add model, repository, service, handler
3. Register routes
4. Add migrations
5. Define Casbin permissions
6. Write tests

---

## 🔍 Quick Reference

### Common Tasks

#### Protect a Route

```go
router.POST("/resource",
    middleware.AuthMiddleware(),
    middleware.CasbinMiddleware("resource", "action"),
    handler.Method,
)
```

#### Add Permission

```go
svc := casbinService.NewCasbinService()
svc.AddPermissionForRole("role", "resource", "action")
svc.AssignRoleToUser(userID, "role")
```

#### Handle Errors

```go
if err != nil {
    utils.RespondWithAppError(c, errors.ErrSomething)
    return
}
```

#### Create Migration

```go
func MigrationXXX() Migration {
    return Migration{
        ID: "XXX_description",
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&Model{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&Model{})
        },
    }
}
```

---

## 🆘 Troubleshooting

### Build Errors

```bash
go mod tidy
go build ./cmd/server
```

### Permission Denied

```go
// Debug permissions
svc := casbinService.NewCasbinService()
roles, _ := svc.GetRolesForUser(userID)
log.Println("User roles:", roles)
```

### Database Issues

```bash
# Check migrations
psql -U user -d database -c "SELECT * FROM migrations;"

# Reset database
dropdb nogo_dev && createdb nogo_dev
go run cmd/migrate/main.go
```

---

## 📞 Need Help?

1. Check the relevant section above
2. Read the detailed docs in each folder
3. Look at code examples in `internal/`
4. Check error codes in `internal/common/errors/`

---

## 🎯 Project Status

- ✅ Authentication (JWT)
- ✅ Authorization (Casbin ABAC/RBAC)
- ✅ Error Handling
- ✅ User Management
- ✅ Role Management
- ✅ Database Migrations
- 🚧 Novel Management (In Progress)
- 🚧 Chapter Management (In Progress)
- 📝 API Documentation (In Progress)

---

**Happy Coding! 🚀**
