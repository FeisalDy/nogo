# Database Migration System

## Overview

This project uses a **phased migration approach** to handle database evolution and complex relationships between models. This approach allows us to:

- Build features incrementally without breaking existing functionality
- Handle complex relationships (like Novel ↔ Genre) without circular dependencies
- Test each phase independently
- Easily roll back changes if needed
- Avoid database constraint errors during development

## How It Works

### 1. Migration Framework

The migration system is built around these core components:

```go
// Migration represents a single database change
type Migration struct {
    ID          string                    // Unique identifier (e.g., "001_create_users")
    Description string                    // Human-readable description
    Up          func(*gorm.DB) error     // Apply the migration
    Down        func(*gorm.DB) error     // Rollback the migration
}

// MigrationHistory tracks which migrations have been applied
type MigrationHistory struct {
    ID          uint   `gorm:"primaryKey"`
    MigrationID string `gorm:"unique;not null"`
    AppliedAt   int64  `gorm:"not null"`
}
```

### 2. Migration Execution Flow

```
1. App starts → database.Init() is called
2. RunMigrations() checks MigrationHistory table
3. For each migration in GetAllMigrations():
   - Check if already applied (skip if yes)
   - Run migration.Up() function
   - Record success in MigrationHistory
4. Continue with app startup
```

### 3. Phased Development Strategy

Instead of creating all tables at once, we use **phases**:

#### **Phase 1: Foundation (✅ ACTIVE)**
```go
Migration001CreateUsers()  // Users table with authentication
```

#### **Phase 2: Content Management (🔄 READY)**
```go
// Migration002CreateNovels()     // Novels table (Novel → User relationship)
// Migration003CreateChapters()   // Chapters table (Chapter → Novel relationship)
```

#### **Phase 3: Genre System (⏳ PREPARED)**
```go
// Migration004CreateGenres()     // Genres table
// Migration005AddNovelGenres()   // Novel ↔ Genre many-to-many relationship
```

## File Structure

```
internal/database/migrations/
├── migrations.go              # Main migration framework
├── 001_create_users.go       # User table migration
├── 002_create_novels.go      # Novel table migration (ready)
├── 003_create_chapters.go    # Chapter table migration (ready)
├── 004_create_genres.go      # Genre table migration (ready)
└── 005_add_novel_genres.go   # Novel-Genre junction table (ready)

cmd/migrate/
└── main.go                   # Migration runner command
```

## Usage Examples

### Running Migrations

Migrations run automatically when the app starts:
```go
// In database.Init()
log.Println("Running database migrations...")
if err := migrations.RunMigrations(DB); err != nil {
    log.Fatalf("failed to run migrations: %v", err)
}
```

Or run manually:
```bash
go run cmd/migrate/main.go
```

### Activating Next Phase

To move to **Phase 2** (Content Management), simply uncomment in `migrations.go`:

```go
func GetAllMigrations() []Migration {
    return []Migration{
        Migration001CreateUsers(),
        
        // Phase 1: Uncomment when ready for Content Management
        Migration002CreateNovels(),     // ← Uncomment this
        Migration003CreateChapters(),   // ← Uncomment this
        
        // Phase 2: Uncomment when ready for Genre System  
        // Migration004CreateGenres(),     
        // Migration005AddNovelGenres(),   
    }
}
```

Next app restart will automatically apply the new migrations.

## Migration Examples

### Simple Table Migration
```go
// 001_create_users.go
func Migration001CreateUsers() Migration {
    return Migration{
        ID:          "001_create_users",
        Description: "Create users table with basic authentication fields",
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&User{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&User{})
        },
    }
}
```

### Relationship Migration
```go
// 002_create_novels.go
type Novel struct {
    ID          uint   `gorm:"primaryKey"`
    Title       string `gorm:"not null"`
    AuthorID    uint   `gorm:"not null"`        // Required relationship
    Author      User   `gorm:"foreignKey:AuthorID"`
    // Note: Genres relationship added later in Phase 3
}
```

### Many-to-Many Relationship Migration
```go
// 005_add_novel_genres.go
type NovelGenre struct {
    NovelID uint  `gorm:"primaryKey"`
    GenreID uint  `gorm:"primaryKey"`
    Novel   Novel `gorm:"foreignKey:NovelID"`
    Genre   Genre `gorm:"foreignKey:GenreID"`
}
```

## Advantages of This Approach

### 🎯 **Incremental Development**
- Build and test one feature at a time
- No complex dependency management
- Each phase is independently functional

### 🔒 **Safe Relationships**
- No circular dependencies during development
- Foreign key constraints work correctly
- Optional relationships added when both tables exist

### 📊 **Clear Development Path**
```
Users → Novels → Chapters → Genres → Novel-Genres → Advanced Features
```

### 🐛 **Easy Debugging**
- Each migration is isolated
- Clear migration history
- Simple rollback process

### 🧪 **Testable**
- Test each phase independently
- Seed data for specific phases
- Integration tests per phase

## Best Practices

### 1. **Never Modify Existing Migrations**
Once a migration is applied, create a new migration for changes:
```go
// ❌ Don't modify 001_create_users.go
// ✅ Create 006_add_user_fields.go
```

### 2. **Keep Migrations Simple**
Each migration should do one thing:
```go
// ✅ Good: Create one table
Migration001CreateUsers()

// ❌ Avoid: Create multiple unrelated tables
Migration001CreateUsersAndNovelAndGenres()
```

### 3. **Test Migrations**
```go
// Test that migration works
func TestMigration001(t *testing.T) {
    db := setupTestDB()
    migration := Migration001CreateUsers()
    
    // Test Up
    err := migration.Up(db)
    assert.NoError(t, err)
    
    // Test Down
    err = migration.Down(db)
    assert.NoError(t, err)
}
```

### 4. **Use Descriptive Names**
```go
// ✅ Clear purpose
"001_create_users"
"002_create_novels" 
"005_add_novel_genres"

// ❌ Unclear
"001_initial"
"002_stuff"
```

## Troubleshooting

### Migration Failed
```bash
# Check migration history
SELECT * FROM migration_histories;

# Manual rollback (if needed)
DROP TABLE novels;
DELETE FROM migration_histories WHERE migration_id = '002_create_novels';
```

### Relationship Issues
```go
// ❌ Problem: Trying to reference non-existent table
type Novel struct {
    GenreID uint  // Genre table doesn't exist yet!
}

// ✅ Solution: Add relationship in later migration
type Novel struct {
    // Add GenreID in Migration005 when Genre table exists
}
```

## Current Status

- ✅ **Phase 1**: User system ready
- 🔄 **Phase 2**: Novel/Chapter migrations prepared (uncomment to activate)
- ⏳ **Phase 3**: Genre system migrations prepared (uncomment to activate)

## Next Steps

1. **Activate Phase 2**: Uncomment Novel/Chapter migrations
2. **Build Novel CRUD**: Create services, handlers, routes
3. **Test Phase 2**: Ensure everything works before Phase 3
4. **Activate Phase 3**: Add Genre system when ready