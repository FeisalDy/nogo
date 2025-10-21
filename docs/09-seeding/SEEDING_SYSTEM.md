# Database Seeding System

## Overview

The seeding system automatically populates the database with sample data on server startup. It intelligently checks if data already exists and only creates new records when needed, making it safe to run on every server start.

## Features

- ✅ **Auto-run on Server Start**: Seeds are automatically executed when the server starts
- ✅ **Idempotent**: Safe to run multiple times - only creates data if it doesn't exist
- ✅ **Ordered Execution**: Seeders run in dependency order to maintain referential integrity
- ✅ **Multi-language Support**: Uses locale codes (zh-CN, ja-JP, ko-KR, en-US, id-ID, etc.)
- ✅ **Comprehensive Coverage**: Seeds all tables including relationships

## Seeder Files

All seeders are located in `internal/database/seeds/`:

| File | Description |
|------|-------------|
| `seeds.go` | Main orchestrator that runs all seeders in order |
| `seed_roles.go` | Seeds user roles (admin, author, user) |
| `seed_users.go` | Seeds sample users with role assignments |
| `seed_genres.go` | Seeds novel genres (Fantasy, Romance, Action, etc.) |
| `seed_tags.go` | Seeds novel tags (Reincarnation, System, Cultivation, etc.) |
| `seed_novels.go` | Seeds sample novels with translations and relationships |
| `seed_chapters.go` | Seeds sample chapters for each novel |
| `seed_casbin_rule.go` | Seeds Casbin RBAC permissions |

## Execution Order

The seeders run in the following order to respect database relationships:

1. **Roles** - User roles (admin, author, user)
2. **Users** - Sample users with role assignments
3. **Genres** - Novel genre categories
4. **Tags** - Novel tags/labels
5. **Novels** - Sample novels with translations
6. **Chapters** - Sample chapters for novels
7. **Casbin Policies** - Permission rules

## Seeded Data

### Roles (3)
- **admin**: Full access to all resources
- **author**: Can create and manage novels/chapters
- **user**: Read-only access

### Users (3)
- **admin@example.com** (admin role)
  - Username: admin
  - Password: password123
  
- **author1@example.com** (author role)
  - Username: author1
  - Password: password123
  
- **john@example.com** (user role)
  - Username: john_doe
  - Password: password123

### Genres (12)
Fantasy, Romance, Action, Mystery, Science Fiction, Horror, Comedy, Drama, Slice of Life, Historical, Martial Arts, Psychological

### Tags (20)
Reincarnation, Overpowered MC, System, Magic, Cultivation, Transmigration, Harem, Revenge, Weak to Strong, Academy, Dungeon, Isekai, Virtual Reality, Monster, Adventure, Kingdom Building, Time Travel, Anti-Hero, Female Lead, Male Lead

### Novels (3)

#### 1. The Legendary Cultivator
- **Original Language**: zh-CN (Chinese)
- **Original Author**: 李明
- **Status**: ongoing
- **Genres**: Fantasy, Action, Martial Arts
- **Tags**: Cultivation, Weak to Strong, Overpowered MC, Magic
- **Translations**: English (en-US), Indonesian (id-ID)
- **Chapters**: 5

#### 2. Reborn in Another World
- **Original Language**: ja-JP (Japanese)
- **Original Author**: 田中太郎
- **Status**: ongoing
- **Genres**: Fantasy, Comedy, Slice of Life
- **Tags**: Isekai, Reincarnation, System, Adventure
- **Translations**: English (en-US), Indonesian (id-ID)
- **Chapters**: 5

#### 3. Shadow Monarch
- **Original Language**: ko-KR (Korean)
- **Original Author**: 김철수
- **Status**: completed
- **Genres**: Action, Fantasy, Horror
- **Tags**: Weak to Strong, System, Dungeon, Monster, Overpowered MC
- **Translations**: English (en-US), Indonesian (id-ID)
- **Chapters**: 5

### Casbin Policies (35)
Permissions for all roles across all resources (users, novels, chapters, genres, tags, roles, media)

## Language Codes Used

The system uses standard locale codes:

- `zh-CN` - Chinese (Simplified)
- `ja-JP` - Japanese
- `ko-KR` - Korean
- `en-US` - English (US)
- `id-ID` - Indonesian

## How It Works

### 1. Server Startup Flow

```go
// cmd/server/main.go
func main() {
    // 1. Initialize database
    database.Init(cfg.DB)
    
    // 2. Initialize Casbin
    casbinService.InitCasbin(database.DB, modelPath)
    
    // 3. Run all seeders
    database.RunSeeds()
    
    // 4. Start server
    router.SetupRoutes(database.DB, cfg.App)
}
```

### 2. Seed Execution

```go
// internal/database/seeds/seeds.go
func RunAllSeeds(db *gorm.DB) error {
    // Runs each seeder in dependency order
    // Only creates records if they don't exist
    // Returns error if any seeder fails
}
```

### 3. Idempotent Check

Each seeder checks if data exists before creating:

```go
var existing Record
result := db.Where("unique_field = ?", value).First(&existing)

if result.Error == gorm.ErrRecordNotFound {
    // Create the record
    db.Create(&record)
} else {
    // Skip - record already exists
}
```

## Adding New Seeders

To add a new seeder:

1. **Create the seeder file**:
   ```go
   // internal/database/seeds/seed_<table>.go
   package seeds
   
   func Seed<Table>(db *gorm.DB) error {
       log.Println("🌱 Seeding <table>...")
       
       // Check and create records
       
       log.Println("✅ <Table> seeding completed")
       return nil
   }
   ```

2. **Add to seeds.go**:
   ```go
   seeders := []struct {
       name string
       fn   func(*gorm.DB) error
   }{
       // ... existing seeders
       {"YourTable", SeedYourTable},
   }
   ```

3. **Place in correct order** based on dependencies

## Logs Output

When seeders run, you'll see output like:

```
🌱 Starting database seeding...
Running seeder: Roles
🌱 Seeding roles...
✅ Created role: admin
✅ Created role: author
✅ Created role: user
✅ Roles seeding completed

Running seeder: Users
🌱 Seeding users...
✅ Created user: admin@example.com with role: admin
✅ Users seeding completed

...

✅ All seeders completed successfully!
```

On subsequent runs with existing data:

```
🌱 Seeding roles...
⏭️  Role already exists: admin
⏭️  Role already exists: author
⏭️  Role already exists: user
✅ Roles seeding completed
```

## Testing

To test the seeders:

1. **Fresh database**:
   ```bash
   # Drop and recreate database
   # Then run server - all data will be seeded
   go run cmd/server/main.go
   ```

2. **Re-run with existing data**:
   ```bash
   # Run again - should skip existing records
   go run cmd/server/main.go
   ```

3. **Verify in database**:
   ```sql
   SELECT COUNT(*) FROM roles;      -- Should be 3
   SELECT COUNT(*) FROM users;      -- Should be 3
   SELECT COUNT(*) FROM genres;     -- Should be 12
   SELECT COUNT(*) FROM tags;       -- Should be 20
   SELECT COUNT(*) FROM novels;     -- Should be 3
   SELECT COUNT(*) FROM chapters;   -- Should be 15 (5 per novel)
   ```

## Migration vs Seeding

| Aspect | Migrations | Seeds |
|--------|-----------|-------|
| **Purpose** | Schema changes | Sample data |
| **When** | Once per environment | Every server start |
| **Tracked** | migration_history table | Not tracked |
| **Idempotent** | By ID | By checking existence |
| **Required** | Yes (production) | Optional (dev/testing) |

## Best Practices

1. ✅ **Use unique constraints** to identify existing records
2. ✅ **Handle errors gracefully** - log warnings but continue
3. ✅ **Maintain referential integrity** - seed in dependency order
4. ✅ **Use realistic data** - helps with testing and demos
5. ✅ **Document seed data** - especially passwords and credentials
6. ✅ **Keep seeds optional** - production may not need them

## Security Note

⚠️ **Important**: The seeded users have the password `password123`. This is only suitable for development/testing. Never use these credentials in production!

For production:
- Disable automatic seeding or use environment variables
- Create admin users manually with strong passwords
- Use proper secret management for credentials

## Troubleshooting

### Seeder fails with foreign key error
- Check seeder execution order in `seeds.go`
- Ensure parent records are created before child records

### Records created on every run
- Verify uniqueness check uses correct fields
- Check that unique constraints exist in database

### Seeder takes too long
- Consider reducing sample data size
- Add database indexes on frequently queried fields
- Check for N+1 query problems

## Related Documentation

- [Migration System](../05-database/MIGRATION_SYSTEM.md)
- [Database Schema](../05-database/CURSOR_PAGINATION_GUIDE.md)
- [RBAC Implementation](../04-authorization/RBAC_IMPLEMENTATION.md)
