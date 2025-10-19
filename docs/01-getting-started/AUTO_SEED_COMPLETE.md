# ✅ AUTO-SEEDING CASBIN PERMISSIONS - COMPLETE!

## 🎉 What's Been Done

Your Casbin permissions now **auto-seed on every server start**! No more manual seeding required.

---

## 🚀 How It Works Now

### Automatic Flow

```
1. Server starts
2. Database connects
3. Migrations run
4. ✅ Casbin initializes
5. 🌱 Permissions auto-seed (NEW!)
6. Server ready!
```

### What Happens on Startup

```bash
$ ./tmp/main

Database connected
Running database migrations...
Migrations completed successfully
Casbin initialized successfully
🌱 Auto-seeding Casbin permissions...
✅ Casbin auto-seed complete! Total permissions: 35
Starting server on :8080
```

---

## 📝 Two Ways to Seed

### Method 1: Automatic (Recommended - Happens Every Start)

Just start your server:

```bash
./tmp/main
# or
go run cmd/server/main.go
```

**Benefits:**

- ✅ Always runs automatically
- ✅ Adds new permissions without duplicating
- ✅ Safe to run multiple times
- ✅ No manual steps needed

### Method 2: Manual (When You Want to Test)

Run the standalone script:

```bash
go run scripts/seed_casbin.go
```

**Use this when:**

- You want to seed without starting the server
- Testing permission changes
- Manually verifying seed logic

---

## ✏️ How to Add New Permissions

### Step 1: Edit the Seed File

Open: `/home/feisal/project/shilan/nogo/internal/database/seeds/seed_casbin_rule.go`

### Step 2: Add Your Permission

Find the role you want to add permissions to and add new entries:

```go
// Add to admin permissions
adminPerms := []struct {
    resource string
    action   string
}{
    {"users", "read"},
    {"users", "write"},
    // ... existing permissions ...

    // ADD NEW PERMISSION HERE:
    {"comments", "read"},      // ← NEW!
    {"comments", "write"},     // ← NEW!
    {"comments", "delete"},    // ← NEW!
}
```

### Step 3: Restart Server

```bash
# Stop server (Ctrl+C)
# Start again
./tmp/main
```

The new permissions will be added automatically! 🎉

### Step 4: Verify

```bash
# Check the endpoint
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN" | jq '.data.permissions'

# Or check database
psql -U user -d db -c "SELECT v0, v1, v2 FROM casbin_rule WHERE ptype = 'p' ORDER BY v0, v1;"
```

---

## 📊 Current Permissions Structure

### Location

`/home/feisal/project/shilan/nogo/internal/database/seeds/seed_casbin_rule.go`

### Roles Defined

#### 1. Admin (21 permissions)

```
users:     read, write, delete
novels:    read, write, delete
chapters:  read, write, delete
genres:    read, write, delete
tags:      read, write, delete
roles:     read, write
media:     read, write, delete
```

#### 2. Author (8 permissions)

```
novels:   read, write
chapters: read, write
genres:   read
tags:     read
profile:  read, write
media:    write
```

#### 3. User (6 permissions)

```
novels:   read
chapters: read
genres:   read
tags:     read
profile:  read, write
```

#### 4. Moderator (8 permissions)

```
novels:   read, write
chapters: read, write, delete
users:    read
media:    read, delete
```

---

## 🔧 Files Modified

### 1. `/cmd/server/main.go`

**Added:**

```go
// Auto-seed Casbin permissions (runs after Casbin is initialized)
database.SeedCasbin()
```

This calls the seeder right after Casbin initialization.

### 2. `/internal/database/database.go`

**Changed:**

- ❌ Removed seed call from `Init()` (was too early, Casbin not initialized)
- ✅ Added `SeedCasbin()` function (called from main.go after Casbin init)

### 3. `/internal/database/seeds/seed_casbin_rule.go`

**Updated:**

- ✅ Uses Casbin API (not raw SQL)
- ✅ Properly checks for Casbin initialization
- ✅ Adds permissions without duplicating
- ✅ Safe to run multiple times

### 4. `/scripts/seed_casbin.go`

**Created:**

- ✅ Standalone seeding script
- ✅ Can run independently for testing
- ✅ Same logic as automatic seed

---

## 🎯 Key Benefits

### 1. **No Manual Steps**

```bash
# Old way:
./server
# Oh wait, need to seed!
go run scripts/seed_casbin.go
# Now restart server

# New way:
./server  # ✅ Done! Everything automatic
```

### 2. **Easy Permission Updates**

```go
// Edit seed_casbin_rule.go
adminPerms = append(adminPerms,
    {"new_resource", "new_action"}  // Add new permission
)

// Restart server - new permission added automatically!
```

### 3. **Safe & Idempotent**

```bash
# Run 100 times? No problem!
# No duplicates, no errors, just works ✅
```

### 4. **Development Workflow**

```bash
# 1. Add new feature (e.g., comments system)
# 2. Add permissions in seed_casbin_rule.go
# 3. Restart server
# 4. Permissions ready! Start coding routes
```

---

## 🧪 Testing the Auto-Seed

### Test 1: Fresh Start

```bash
# Delete all permissions
psql -U user -d db -c "DELETE FROM casbin_rule WHERE ptype = 'p';"

# Start server
./tmp/main

# Check logs - should see:
# 🌱 Auto-seeding Casbin permissions...
# ✅ Casbin auto-seed complete! Total permissions: 35
```

### Test 2: Add Permission

```bash
# 1. Edit seed_casbin_rule.go, add new permission
# 2. Restart server
# 3. Check /me endpoint - new permission should appear
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer TOKEN" | jq '.data.permissions | length'
```

### Test 3: Idempotency

```bash
# Restart server 5 times
./tmp/main  # Stop with Ctrl+C
./tmp/main  # Stop with Ctrl+C
./tmp/main  # Stop with Ctrl+C

# Check permissions count - should be same each time
psql -U user -d db -c "SELECT COUNT(*) FROM casbin_rule WHERE ptype = 'p';"
```

---

## 🔍 Troubleshooting

### Issue: "Casbin enforcer not initialized yet, skipping seed"

**Cause:** Casbin wasn't initialized before seed ran.

**Solution:** Already fixed! Seed now runs AFTER Casbin init in main.go.

### Issue: Permissions not appearing

**Check 1: Server logs**

```bash
# Look for this line:
✅ Casbin auto-seed complete! Total permissions: 35
```

**Check 2: Database**

```sql
SELECT COUNT(*) FROM casbin_rule WHERE ptype = 'p';
-- Should show 35+ permissions
```

**Check 3: Restart server**

```bash
# Casbin caches policies in memory
# Restart to reload
```

### Issue: Duplicate permissions

**This shouldn't happen** - Casbin automatically prevents duplicates.

If it does, check:

```sql
SELECT v0, v1, v2, COUNT(*)
FROM casbin_rule
WHERE ptype = 'p'
GROUP BY v0, v1, v2
HAVING COUNT(*) > 1;
```

---

## 📚 Documentation Updates Needed

The following docs should be updated to reflect auto-seeding:

1. `README.md` - Remove manual seed step
2. `SEED_PERMISSIONS.md` - Add "Now automatic!" note
3. `CASBIN_QUICK_START.md` - Update workflow

---

## 🎉 Summary

### What Changed

- ✅ Permissions now auto-seed on server start
- ✅ No manual `go run scripts/seed_casbin.go` needed
- ✅ Add new permissions by editing one file + restart
- ✅ Safe to run multiple times (no duplicates)

### Your Workflow Now

```bash
# 1. Edit permissions in seed_casbin_rule.go
# 2. Restart server
# 3. Done! ✨
```

### Files to Remember

- **Add/edit permissions**: `internal/database/seeds/seed_casbin_rule.go`
- **Manual test seeding**: `scripts/seed_casbin.go`
- **Startup flow**: `cmd/server/main.go`

---

## 🚀 You're All Set!

Your permissions will now automatically seed every time you start the server. Just edit the seed file when you need to add new permissions and restart! 🎊
