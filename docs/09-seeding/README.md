# Database Seeding - Quick Start

## Overview

✅ **Automatic seeding system implemented successfully!**

The database seeding system automatically populates all tables with sample data when the server starts. It's idempotent - safe to run multiple times without duplicating data.

## What Was Created

### New Seeder Files

All files in `internal/database/seeds/`:

| File | Records Created |
|------|----------------|
| `seeds.go` | Main orchestrator |
| `seed_roles.go` | 3 roles |
| `seed_users.go` | 3 users with passwords |
| `seed_genres.go` | 12 genres |
| `seed_tags.go` | 20 tags |
| `seed_novels.go` | 3 novels with translations |
| `seed_chapters.go` | 15 chapters (5 per novel) |
| `seed_casbin_rule.go` | 35 permission policies |

### Sample Data Created

#### Users (password: `password123`)
- `admin@example.com` (admin role)
- `author1@example.com` (author role)  
- `john@example.com` (user role)

#### Novels
1. **The Legendary Cultivator** (zh-CN) - Chinese cultivation novel
2. **Reborn in Another World** (ja-JP) - Japanese isekai novel
3. **Shadow Monarch** (ko-KR) - Korean dungeon novel

Each novel has:
- Translations in English (en-US) and Indonesian (id-ID)
- 5 chapters with translated content
- Assigned genres and tags

#### Languages Used
- `zh-CN` (Chinese Simplified)
- `ja-JP` (Japanese)
- `ko-KR` (Korean)
- `en-US` (English)
- `id-ID` (Indonesian)

## How It Works

### Auto-run on Server Start

```go
// cmd/server/main.go
database.Init(cfg.DB)              // 1. Initialize database
casbinService.InitCasbin(...)      // 2. Initialize Casbin
database.RunSeeds()                // 3. Run all seeders ← NEW!
router.SetupRoutes(...)            // 4. Start server
```

### Idempotent Execution

On first run:
```
🌱 Seeding roles...
✅ Created role: admin
✅ Created role: author
✅ Created role: user
```

On subsequent runs:
```
🌱 Seeding roles...
⏭️  Role already exists: admin
⏭️  Role already exists: author
⏭️  Role already exists: user
```

## Changes Made

### 1. Removed Role Seeding from Migration

**Before:** Migration `008_seed_roles.go` created roles
**After:** Migration is now a no-op, roles seeded via `seed_roles.go`

### 2. Updated Database Initialization

**File:** `internal/database/database.go`

```go
// OLD
func SeedCasbin() {
    seeds.SeedCasbinPolicies(DB)
}

// NEW
func RunSeeds() {
    seeds.RunAllSeeds(DB)  // Runs all seeders including Casbin
}
```

### 3. Updated Server Startup

**File:** `cmd/server/main.go`

```go
// OLD
database.SeedCasbin()

// NEW
database.RunSeeds()  // Runs all seeders
```

## Testing

### Test First Run (Empty Database)

```bash
# Drop database (if needed)
# Then run server
go run cmd/server/main.go
```

Expected output:
```
🌱 Starting database seeding...
Running seeder: Roles
✅ Created role: admin
✅ Created role: author
✅ Created role: user
...
✅ All seeders completed successfully!
```

### Test Idempotency (Existing Data)

```bash
# Run server again
go run cmd/server/main.go
```

Expected output:
```
🌱 Starting database seeding...
Running seeder: Roles
⏭️  Role already exists: admin
⏭️  Role already exists: author
⏭️  Role already exists: user
...
✅ All seeders completed successfully!
```

### Verify Data

```sql
-- Check record counts
SELECT COUNT(*) FROM roles;           -- 3
SELECT COUNT(*) FROM users;           -- 3
SELECT COUNT(*) FROM genres;          -- 12
SELECT COUNT(*) FROM tags;            -- 20
SELECT COUNT(*) FROM novels;          -- 3
SELECT COUNT(*) FROM chapters;        -- 15
SELECT COUNT(*) FROM novel_translations;  -- 6
SELECT COUNT(*) FROM chapter_translations; -- 30
SELECT COUNT(*) FROM casbin_rule;     -- 35
```

## Key Features

✅ **Auto-run** - Executes on every server start
✅ **Idempotent** - Safe to run multiple times
✅ **Smart Detection** - Only creates missing records
✅ **Multi-language** - Uses proper locale codes (zh-CN, ja-JP, etc.)
✅ **Complete Coverage** - Seeds all tables with relationships
✅ **Ordered Execution** - Respects foreign key dependencies
✅ **Clear Logging** - Shows what's created vs skipped

## Usage

### For Development

Just run the server - seeds will automatically populate:

```bash
go run cmd/server/main.go
```

### For Testing

The seeded data provides:
- Test users with known passwords
- Sample content for UI testing
- Realistic multi-language data
- Complete relationship examples

### For Production

**⚠️ Security Warning:** Consider disabling auto-seeding in production or using environment variables to control it:

```go
// Option 1: Environment variable
if os.Getenv("AUTO_SEED") == "true" {
    database.RunSeeds()
}

// Option 2: Config flag
if cfg.App.AutoSeed {
    database.RunSeeds()
}
```

## Adding New Seeders

1. Create `internal/database/seeds/seed_<table>.go`
2. Implement `Seed<Table>(db *gorm.DB) error`
3. Add to `seeds.go` in correct dependency order
4. Use helper functions `strPtr()` and `intPtr()`

Example:

```go
package seeds

func SeedComments(db *gorm.DB) error {
    log.Println("🌱 Seeding comments...")
    
    // Check if exists
    var existing Comment
    result := db.Where("unique_field = ?", value).First(&existing)
    
    if result.Error == gorm.ErrRecordNotFound {
        // Create new record
        db.Create(&comment)
        log.Printf("✅ Created comment")
    } else {
        log.Printf("⏭️  Comment already exists")
    }
    
    return nil
}
```

## Documentation

Full documentation: [docs/09-seeding/SEEDING_SYSTEM.md](docs/09-seeding/SEEDING_SYSTEM.md)

## Summary

✅ Comprehensive seeding system implemented
✅ All tables have sample data
✅ Multi-language support (zh-CN, ja-JP, ko-KR, en-US, id-ID)
✅ Auto-runs on server start
✅ Idempotent - safe to run multiple times
✅ Role seeding moved from migration to seeder
✅ Well documented with examples

**You can now start the server and have a fully populated database ready for testing!**
