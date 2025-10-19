# 01. Getting Started

Quick start guides to get you up and running.

## 📄 Documents

### ⚡ [Casbin Quick Start](CASBIN_QUICK_START.md)

**Start here!** Get Casbin running in 5 minutes.

- Setup steps
- Seed permissions script
- Protect routes
- Test permissions

### 📋 [Casbin Implementation Summary](CASBIN_IMPLEMENTATION_SUMMARY.md)

Complete overview of the Casbin implementation.

- What was added
- How it works
- Key features
- Usage examples
- Next steps

### ✅ [Cleanup Complete](CLEANUP_COMPLETE.md)

Recent permission table cleanup documentation.

- What changed
- Database schema
- Migration updates
- Next steps

## Quick Setup

```bash
# 1. Install dependencies
go mod download

# 2. Setup database
createdb nogo_dev

# 3. Run migrations
go run cmd/migrate/main.go

# 4. Seed Casbin permissions
go run scripts/seed_permissions.go

# 5. Start server
go run cmd/server/main.go
```

## First API Request

```bash
# Register
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

[← Back to Main Documentation](../README.md)
