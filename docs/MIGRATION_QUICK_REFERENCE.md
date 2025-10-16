# Migration Quick Reference

## 🚀 How to Activate Next Phase

### Currently Active: Phase 1 (Users)
```
✅ Users table created
✅ Authentication ready
```

### To Activate Phase 2 (Novels + Chapters):

1. **Edit**: `internal/database/migrations/migrations.go`
2. **Uncomment** these lines:
   ```go
   func GetAllMigrations() []Migration {
       return []Migration{
           Migration001CreateUsers(),
           
           // Phase 1: Uncomment when ready for Content Management
           Migration002CreateNovels(),     // ← UNCOMMENT THIS
           Migration003CreateChapters(),   // ← UNCOMMENT THIS
           
           // Phase 2: Keep commented for now
           // Migration004CreateGenres(),     
           // Migration005AddNovelGenres(),   
       }
   }
   ```
3. **Restart** the application or run:
   ```bash
   go run cmd/migrate/main.go
   ```

### To Activate Phase 3 (Genres):

1. **After Phase 2 is working**, uncomment:
   ```go
   Migration004CreateGenres(),     // ← UNCOMMENT THIS
   Migration005AddNovelGenres(),   // ← UNCOMMENT THIS
   ```

## 📋 Migration Commands

```bash
# Run migrations manually
go run cmd/migrate/main.go

# Build and run server (auto-runs migrations)
go run cmd/server/main.go

# Check if migrations work
go build ./cmd/migrate
```

## 🔍 Check Migration Status

```sql
-- Check which migrations have been applied
SELECT * FROM migration_histories ORDER BY applied_at;

-- Check if tables exist
\dt  -- In PostgreSQL
```

## 📊 Current Database Schema

### Phase 1 (Active)
```
users
├── id (primary key)
├── created_at
├── updated_at  
├── deleted_at
├── name
├── email (unique)
├── password
├── role (default: 'user')
└── is_active (default: true)

migration_histories
├── id (primary key)
├── migration_id (unique)
└── applied_at
```

### Phase 2 (Ready to activate)
```
novels
├── id (primary key)
├── created_at, updated_at, deleted_at
├── title
├── description
├── status (default: 'draft')
├── author_id → users.id
└── author (relationship)

chapters  
├── id (primary key)
├── created_at, updated_at, deleted_at
├── title
├── content
├── number
├── is_public (default: false)
├── novel_id → novels.id
└── novel (relationship)
```

### Phase 3 (Prepared)
```
genres
├── id (primary key)  
├── created_at, updated_at, deleted_at
├── name (unique)
├── description
└── slug (unique)

novel_genres (junction table)
├── novel_id → novels.id
├── genre_id → genres.id
├── novel (relationship)
└── genre (relationship)
```

## 🛠️ Development Workflow

1. **Work on Current Phase**: Build features for active tables
2. **Test Thoroughly**: Ensure current phase works perfectly  
3. **Activate Next Phase**: Uncomment next migrations
4. **Update Models**: Add new relationships to existing models
5. **Repeat**: Continue to next phase

## ⚠️ Important Notes

- **Never modify existing migrations** once applied
- **Always test** before activating next phase  
- **One phase at a time** - don't skip ahead
- **Check logs** for migration errors
- **Backup database** before major changes

## 🔄 Rollback (Emergency)

```sql
-- If you need to rollback a migration manually:
-- 1. Drop the table
DROP TABLE IF EXISTS novels;

-- 2. Remove from migration history  
DELETE FROM migration_histories WHERE migration_id = '002_create_novels';

-- 3. Restart application
```

## 📞 Need Help?

1. Check `docs/MIGRATION_SYSTEM.md` for detailed explanation
2. Look at existing migration files in `internal/database/migrations/`
3. Test migrations on development database first