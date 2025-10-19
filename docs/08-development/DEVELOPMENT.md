# Development Guide

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Setup Instructions](#setup-instructions)
3. [Migration-Based Development](#migration-based-development)
4. [Phase-by-Phase Development](#phase-by-phase-development)
5. [Code Standards](#code-standards)
6. [Expanding to New Domains](#expanding-to-new-domains)
7. [Testing Guidelines](#testing-guidelines)
8. [Database Management](#database-management)
9. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software
- **Go**: Version 1.21 or higher
  - [Installation Guide](https://golang.org/doc/install)
  - Verify: `go version`

- **PostgreSQL**: Version 12 or higher
  - [Installation Guide](https://www.postgresql.org/download/)
  - Verify: `psql --version`

- **Git**: For version control
  - [Installation Guide](https://git-scm.com/downloads)
  - Verify: `git --version`

### Recommended Tools
- **Air**: For hot reloading during development
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

- **golangci-lint**: For code linting
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- **Postman** or **Insomnia**: For API testing

## Setup Instructions

### 1. Clone Repository
```bash
git clone <repository-url>
cd boiler
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Database Setup

#### Create Database
```sql
-- Connect to PostgreSQL
psql -U postgres

-- Create database
CREATE DATABASE boiler_dev;

-- Create user (optional)
CREATE USER boiler_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE boiler_dev TO boiler_user;
```

#### Environment Configuration
Create a `.env` file in the project root:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=boiler_dev

# Application Configuration
PORT=8080
GIN_MODE=debug

# Future configurations
# JWT_SECRET=your_jwt_secret
# LOG_LEVEL=debug
```

### 4. Run Application

#### With Air (Hot Reloading)
```bash
air
```

#### With Go Run
```bash
go run cmd/server/main.go
```

#### Build and Run Binary
```bash
go build -o my-app cmd/server/main.go
./my-app
```

### 5. Run Initial Migration
The application uses a **phased migration system**. On first run, it will create the users table:

```bash
# Run the application (migrations run automatically)
go run cmd/server/main.go

# Or run migrations only
go run cmd/migrate/main.go
```

You should see output like:
```
Running database migrations...
Running migration 001_create_users: Create users table with basic authentication fields
Migration 001_create_users completed successfully
Migrations completed successfully
```

### 6. Verify Setup
Test the API:
```bash
curl http://localhost:8080/ping
# Expected response: {"message":"pong"}
```

Check the database:
```sql
-- Connect to your database
psql -d boiler_dev

-- Verify tables were created
\dt
-- Should show: migration_histories, users

-- Check migration status
SELECT * FROM migration_histories;
```

## Migration-Based Development

This project uses a **phased migration approach** to handle complex database relationships and avoid circular dependencies. Instead of building everything at once, we develop incrementally:

### Why Migration-Based Development?

**Traditional Approach Problems**:
```go
// ❌ Circular dependency issue
type Novel struct {
    GenreID uint   // Genre doesn't exist yet!
    Genres  []Genre // Can't reference non-existent model
}

type Genre struct {
    Novels []Novel // Circular reference!
}
```

**Our Solution**:
```go
// ✅ Phase 1: Users only
type User struct { ... }

// ✅ Phase 2: Add novels (User → Novel relationship)  
type Novel struct {
    AuthorID uint // References existing User table
    // Genres added later in Phase 3
}

// ✅ Phase 3: Add genres and relationships
type Genre struct { ... }
type NovelGenre struct { ... } // Junction table
```

### Migration System Overview

```
Phase 1 (✅ Active): Users
├── migration_histories table
└── users table

Phase 2 (🔄 Ready): Content Management  
├── novels table (Novel → User)
└── chapters table (Chapter → Novel)

Phase 3 (⏳ Prepared): Genre System
├── genres table  
└── novel_genres table (Novel ↔ Genre)
```

### Current Status: Phase 1 Complete
- ✅ User authentication system ready
- ✅ Migration framework implemented
- 🔄 Ready to activate Phase 2 (Novels + Chapters)

## Phase-by-Phase Development

### 🎯 Step 1: Understanding Current Phase

Check what's currently active:
```bash
# Check migration status
go run cmd/migrate/main.go

# Check database tables
psql -d boiler_dev -c "\dt"
```

### 🎯 Step 2: Activating Next Phase (Content Management)

When ready to add novels and chapters:

1. **Edit Migration File**:
   ```bash
   # Edit: internal/database/migrations/migrations.go
   ```

2. **Uncomment Phase 2 Migrations**:
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

3. **Restart Application**:
   ```bash
   go run cmd/server/main.go
   # New migrations will run automatically
   ```

4. **Verify New Tables**:
   ```sql
   \dt  -- Should show: users, novels, chapters, migration_histories
   ```

### 🎯 Step 3: Building Domain Logic

After activating Phase 2 migrations, build the domain logic:

#### A. Create Novel Domain Structure
```bash
mkdir -p internal/novel/{model,repository,service,handler,dto}
```

#### B. Implement Novel Model
```go
// internal/novel/model/novel.go
package model

import (
    "time"
    "gorm.io/gorm"
    userModel "github.com/FeisalDy/nogo/internal/user/model"
)

type Novel struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    Title       string `json:"title" gorm:"not null"`
    Description string `json:"description" gorm:"type:text"`
    Status      string `json:"status" gorm:"default:draft"`
    
    // Required relationship (exists in Phase 2)
    AuthorID    uint              `json:"author_id" gorm:"not null"`
    Author      userModel.User    `json:"author" gorm:"foreignKey:AuthorID"`
    
    // Optional relationships (add in Phase 3)
    // Genres   []genreModel.Genre `json:"genres" gorm:"many2many:novel_genres;"`
}
```

#### C. Implement Repository Pattern
```go
// internal/novel/repository/novel_repository.go
package repository

import (
    "github.com/FeisalDy/nogo/internal/database"
    "github.com/FeisalDy/nogo/internal/novel/model"
)

type NovelRepository struct{}

func NewNovelRepository() *NovelRepository {
    return &NovelRepository{}
}

func (r *NovelRepository) Create(novel *model.Novel) error {
    return database.DB.Create(novel).Error
}

func (r *NovelRepository) GetByID(id uint) (*model.Novel, error) {
    var novel model.Novel
    err := database.DB.Preload("Author").First(&novel, id).Error
    return &novel, err
}

func (r *NovelRepository) GetByAuthor(authorID uint) ([]model.Novel, error) {
    var novels []model.Novel
    err := database.DB.Where("author_id = ?", authorID).Find(&novels).Error
    return novels, err
}
```

### 🎯 Step 4: Repeat for Each Domain

Follow the same pattern for chapters, then genres when ready.

## Development Workflow

### 1. Starting Development
```bash
# Pull latest changes
git pull origin main

# Create feature branch based on current phase
git checkout -b feature/phase2-novels

# Start development server
air
```

### 2. Phase-Based Feature Development

**For New Domain in Current Phase**:
1. Check if required migrations are active
2. Create domain structure (model, repository, service, handler)
3. Implement business logic
4. Add routes and test endpoints
5. Write tests

**For New Phase Activation**:
1. Uncomment next phase migrations
2. Test migration runs successfully  
3. Create new domain structures
4. Build incrementally
5. Test thoroughly before next phase

### 3. Committing Changes
```bash
# Stage changes
git add .

# Use descriptive commits with phase context
git commit -m "feat(phase2): add novel CRUD operations"
git commit -m "migration: activate phase 2 (novels + chapters)"

# Push to remote branch
git push origin feature/phase2-novels
```

## Code Standards

### Go Code Style
Follow standard Go conventions:

#### Package Naming
```go
// Good
package user
package middleware

// Avoid
package userService
package User
```

#### Function Naming
```go
// Exported functions (public)
func CreateUser(user *User) error {}

// Unexported functions (private)
func validateEmail(email string) bool {}
```

#### Variable Naming
```go
// Good
var userCount int
var dbConnection *gorm.DB

// Avoid
var u int
var db *gorm.DB // too generic in global scope
```

#### Error Handling
```go
// Always handle errors explicitly
user, err := userService.GetUser(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Use error wrapping for context
if err := db.Create(&user).Error; err != nil {
    return fmt.Errorf("creating user in database: %w", err)
}
```

### Project Structure Standards

#### File Naming
- Use lowercase with underscores: `user_service.go`
- Test files: `user_service_test.go`
- Interface files: `user_repository_interface.go`

#### Directory Organization
```
internal/
└── domain/
    ├── dto/           # Request/response structures
    ├── handler/       # HTTP handlers
    ├── model/         # Domain models
    ├── repository/    # Data access layer
    └── service/       # Business logic
```

### Documentation Standards

#### Function Documentation
```go
// CreateUser creates a new user in the system.
// It validates the user data and returns an error if validation fails.
func CreateUser(user *User) error {
    // implementation
}
```

#### Struct Documentation
```go
// User represents a user in the system.
// It contains basic profile information and authentication details.
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"unique;not null"`
}
```

## Expanding to New Domains

### Overview: Domain-Driven Development

Our project follows **Domain-Driven Design (DDD)** where each business domain (user, novel, chapter, genre) has its own:
- **Model**: Data structure and business rules
- **Repository**: Data access layer
- **Service**: Business logic layer  
- **Handler**: HTTP API layer
- **DTO**: Data Transfer Objects for API

### Step-by-Step: Adding a New Domain

#### 1. Planning & Phase Check

Before adding a domain, determine:
- **Does it need new database tables?** → Check if migration is needed
- **What relationships does it have?** → Ensure dependent tables exist
- **Which phase does it belong to?** → Current or future phase

Example Questions:
```
Adding "Chapter" domain:
✅ Depends on Novel → Novel exists? Check migrations
✅ Chapter → Novel relationship → Simple foreign key
✅ Can be added in same phase as Novel
```

#### 2. Create Domain Structure

```bash
# Standard domain structure
mkdir -p internal/newdomain/{model,repository,service,handler,dto}
```

#### 3. Implementation Order (Critical for Success)

**Always follow this order to avoid circular imports:**

##### A. Model First
```go
// internal/newdomain/model/newdomain.go
package model

import (
    "time"
    "gorm.io/gorm"
    // Import other models CAREFULLY to avoid circular imports
)

type NewDomain struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Your fields here
    Name      string `json:"name" gorm:"not null"`
    
    // Relationships (only to existing domains)
    UserID    uint   `json:"user_id" gorm:"not null"`
    // User   userModel.User `json:"user" gorm:"foreignKey:UserID"`
}
```

##### B. Repository (Data Layer)
```go
// internal/newdomain/repository/newdomain_repository.go
package repository

import (
    "github.com/FeisalDy/nogo/internal/database"
    "github.com/FeisalDy/nogo/internal/newdomain/model"
)

type NewDomainRepository struct{}

func NewNewDomainRepository() *NewDomainRepository {
    return &NewDomainRepository{}
}

func (r *NewDomainRepository) Create(item *model.NewDomain) error {
    return database.DB.Create(item).Error
}

func (r *NewDomainRepository) GetByID(id uint) (*model.NewDomain, error) {
    var item model.NewDomain
    err := database.DB.First(&item, id).Error
    return &item, err
}

func (r *NewDomainRepository) GetAll() ([]model.NewDomain, error) {
    var items []model.NewDomain
    err := database.DB.Find(&items).Error
    return items, err
}

func (r *NewDomainRepository) Update(item *model.NewDomain) error {
    return database.DB.Save(item).Error
}

func (r *NewDomainRepository) Delete(id uint) error {
    return database.DB.Delete(&model.NewDomain{}, id).Error
}
```

##### C. DTO (Data Transfer Objects)
```go
// internal/newdomain/dto/newdomain_dto.go
package dto

// Request DTOs
type CreateNewDomainRequest struct {
    Name   string `json:"name" binding:"required"`
    UserID uint   `json:"user_id" binding:"required"`
}

type UpdateNewDomainRequest struct {
    Name string `json:"name" binding:"required"`
}

// Response DTOs
type NewDomainResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    UserID    uint   `json:"user_id"`
    CreatedAt string `json:"created_at"`
}
```

##### D. Service (Business Logic)
```go
// internal/newdomain/service/newdomain_service.go
package service

import (
    "fmt"
    "github.com/FeisalDy/nogo/internal/newdomain/model"
    "github.com/FeisalDy/nogo/internal/newdomain/repository"
    "github.com/FeisalDy/nogo/internal/newdomain/dto"
)

type NewDomainService struct {
    repository *repository.NewDomainRepository
}

func NewNewDomainService(repo *repository.NewDomainRepository) *NewDomainService {
    return &NewDomainService{repository: repo}
}

func (s *NewDomainService) Create(req *dto.CreateNewDomainRequest) (*model.NewDomain, error) {
    // Business validation
    if req.Name == "" {
        return nil, fmt.Errorf("name is required")
    }

    item := &model.NewDomain{
        Name:   req.Name,
        UserID: req.UserID,
    }

    if err := s.repository.Create(item); err != nil {
        return nil, fmt.Errorf("failed to create: %w", err)
    }

    return item, nil
}

func (s *NewDomainService) GetByID(id uint) (*model.NewDomain, error) {
    return s.repository.GetByID(id)
}

func (s *NewDomainService) Update(id uint, req *dto.UpdateNewDomainRequest) (*model.NewDomain, error) {
    item, err := s.repository.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("item not found: %w", err)
    }

    item.Name = req.Name
    
    if err := s.repository.Update(item); err != nil {
        return nil, fmt.Errorf("failed to update: %w", err)
    }

    return item, nil
}

func (s *NewDomainService) Delete(id uint) error {
    return s.repository.Delete(id)
}
```

##### E. Handler (HTTP API)
```go
// internal/newdomain/handler/newdomain_handler.go
package handler

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/FeisalDy/nogo/internal/newdomain/service"
    "github.com/FeisalDy/nogo/internal/newdomain/dto"
)

type NewDomainHandler struct {
    service *service.NewDomainService
}

func NewNewDomainHandler(service *service.NewDomainService) *NewDomainHandler {
    return &NewDomainHandler{service: service}
}

func (h *NewDomainHandler) Create(c *gin.Context) {
    var req dto.CreateNewDomainRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    item, err := h.service.Create(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, item)
}

func (h *NewDomainHandler) GetByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    item, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.JSON(http.StatusOK, item)
}

func (h *NewDomainHandler) Update(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    var req dto.UpdateNewDomainRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    item, err := h.service.Update(uint(id), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, item)
}

func (h *NewDomainHandler) Delete(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    if err := h.service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
```

##### F. Routes Integration
```go
// internal/newdomain/routes.go
package newdomain

import (
    "github.com/gin-gonic/gin"
    "github.com/FeisalDy/nogo/internal/newdomain/handler"
    "github.com/FeisalDy/nogo/internal/newdomain/repository"
    "github.com/FeisalDy/nogo/internal/newdomain/service"
)

func SetupRoutes(router *gin.RouterGroup) {
    // Initialize dependencies
    repo := repository.NewNewDomainRepository()
    svc := service.NewNewDomainService(repo)
    h := handler.NewNewDomainHandler(svc)

    // Setup routes
    routes := router.Group("/newdomain")
    {
        routes.POST("/", h.Create)
        routes.GET("/:id", h.GetByID)
        routes.PUT("/:id", h.Update)
        routes.DELETE("/:id", h.Delete)
    }
}
```

#### 4. Integration with Main Application

Update `cmd/server/main.go` or router:

```go
import (
    // ... other imports
    "github.com/FeisalDy/nogo/internal/newdomain"
)

func main() {
    // ... existing setup
    
    // Setup API routes
    api := r.Group("/api/v1")
    {
        // Existing routes
        user.SetupRoutes(api)
        
        // Add new domain
        newdomain.SetupRoutes(api)
    }
}
```

### Real Example: Adding Novel Domain to Phase 2

Following the exact pattern for our novel reading platform:

1. **Activate Migration** (if needed):
   ```go
   // Uncomment in migrations.go
   Migration002CreateNovels(),
   ```

2. **Create Structure**:
   ```bash
   mkdir -p internal/novel/{model,repository,service,handler,dto}
   ```

3. **Build in Order**: Model → Repository → DTO → Service → Handler → Routes

This approach ensures:
- ✅ No circular dependencies
- ✅ Clear separation of concerns
- ✅ Testable components
- ✅ Consistent API patterns
- ✅ Easy to maintain and extend

## Testing Guidelines

### Unit Tests
Create test files alongside source files:

```go
// user_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    // Arrange
    service := NewUserService(mockRepository)
    user := &model.User{Name: "Test", Email: "test@example.com"}

    // Act
    err := service.CreateUser(user)

    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### Integration Tests
Test complete workflows:

```go
func TestCreateUserEndpoint(t *testing.T) {
    // Setup test database
    // Create test server
    // Make HTTP request
    // Verify response
}
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/user/service/
```

## Database Management

### Migration System

This project uses a **custom migration system** instead of GORM AutoMigrate for better control:

#### How It Works
```go
// Each migration is a struct with Up/Down functions
type Migration struct {
    ID          string                    // "001_create_users"
    Description string                    // Human readable
    Up          func(*gorm.DB) error     // Apply changes
    Down        func(*gorm.DB) error     // Rollback changes
}
```

#### Migration Files Location
```
internal/database/migrations/
├── migrations.go              # Framework + migration list
├── 001_create_users.go       # Phase 1: Users
├── 002_create_novels.go      # Phase 2: Novels (ready)
├── 003_create_chapters.go    # Phase 2: Chapters (ready)
├── 004_create_genres.go      # Phase 3: Genres (ready)
└── 005_add_novel_genres.go   # Phase 3: Novel-Genre relations (ready)
```

#### Adding New Migration

1. **Create Migration File**:
   ```go
   // internal/database/migrations/006_add_bookmarks.go
   package migrations
   
   import "gorm.io/gorm"
   
   type Bookmark struct {
       ID      uint `gorm:"primaryKey"`
       UserID  uint `gorm:"not null"`
       NovelID uint `gorm:"not null"`
       // ... other fields
   }
   
   func Migration006AddBookmarks() Migration {
       return Migration{
           ID:          "006_add_bookmarks",
           Description: "Create bookmarks table for user-novel relationships",
           Up: func(db *gorm.DB) error {
               return db.AutoMigrate(&Bookmark{})
           },
           Down: func(db *gorm.DB) error {
               return db.Migrator().DropTable(&Bookmark{})
           },
       }
   }
   ```

2. **Add to Migration List**:
   ```go
   // internal/database/migrations/migrations.go
   func GetAllMigrations() []Migration {
       return []Migration{
           Migration001CreateUsers(),
           Migration002CreateNovels(),
           Migration003CreateChapters(),
           Migration004CreateGenres(),
           Migration005AddNovelGenres(),
           Migration006AddBookmarks(),  // Add here
       }
   }
   ```

3. **Run Migration**:
   ```bash
   go run cmd/migrate/main.go
   # or restart the app (migrations run automatically)
   ```

#### Migration Commands

```bash
# Run all pending migrations
go run cmd/migrate/main.go

# Check migration status
psql -d boiler_dev -c "SELECT * FROM migration_histories ORDER BY applied_at;"

# Manual rollback (emergency only)
psql -d boiler_dev -c "DROP TABLE IF EXISTS bookmarks; DELETE FROM migration_histories WHERE migration_id = '006_add_bookmarks';"
```

### Database Seeding

Create seed data for development and testing:

```go
// internal/database/seeds/seeds.go
package seeds

import (
    "log"
    "gorm.io/gorm"
    userModel "github.com/FeisalDy/nogo/internal/user/model"
    novelModel "github.com/FeisalDy/nogo/internal/novel/model"
)

func SeedDatabase(db *gorm.DB) error {
    log.Println("Seeding database...")
    
    if err := seedUsers(db); err != nil {
        return err
    }
    
    if err := seedNovels(db); err != nil {
        return err
    }
    
    log.Println("Database seeding completed")
    return nil
}

func seedUsers(db *gorm.DB) error {
    users := []userModel.User{
        {
            Name:     "John Author",
            Email:    "john@example.com",
            Password: "hashed_password", // Hash in real implementation
            Role:     "author",
            IsActive: true,
        },
        {
            Name:     "Jane Reader", 
            Email:    "jane@example.com",
            Password: "hashed_password",
            Role:     "user",
            IsActive: true,
        },
    }
    
    for _, user := range users {
        // Only create if doesn't exist
        var existing userModel.User
        if err := db.Where("email = ?", user.Email).First(&existing).Error; err == gorm.ErrRecordNotFound {
            if err := db.Create(&user).Error; err != nil {
                return err
            }
            log.Printf("Created user: %s", user.Email)
        }
    }
    
    return nil
}

func seedNovels(db *gorm.DB) error {
    // Only seed if novels table exists (Phase 2 active)
    if !db.Migrator().HasTable(&novelModel.Novel{}) {
        log.Println("Novels table doesn't exist, skipping novel seeds")
        return nil
    }
    
    // Get author user
    var author userModel.User
    if err := db.Where("email = ?", "john@example.com").First(&author).Error; err != nil {
        log.Println("Author not found, skipping novel seeds")
        return nil
    }
    
    novels := []novelModel.Novel{
        {
            Title:       "The Great Adventure",
            Description: "An epic tale of heroes and magic",
            Status:      "published",
            AuthorID:    author.ID,
        },
        {
            Title:       "Mystery of the Old Castle",
            Description: "A thrilling mystery novel",
            Status:      "draft",
            AuthorID:    author.ID,
        },
    }
    
    for _, novel := range novels {
        var existing novelModel.Novel
        if err := db.Where("title = ? AND author_id = ?", novel.Title, novel.AuthorID).First(&existing).Error; err == gorm.ErrRecordNotFound {
            if err := db.Create(&novel).Error; err != nil {
                return err
            }
            log.Printf("Created novel: %s", novel.Title)
        }
    }
    
    return nil
}
```

#### Add Seeding to Main App

```go
// cmd/server/main.go or separate command
import "github.com/FeisalDy/nogo/internal/database/seeds"

func main() {
    // After migrations
    database.Init(cfg.DB)
    
    // Seed database (development only)
    if cfg.App.Environment == "development" {
        if err := seeds.SeedDatabase(database.DB); err != nil {
            log.Printf("Warning: Database seeding failed: %v", err)
        }
    }
}
```

### Environment-Specific Database Config

```bash
# .env.development
DB_NAME=boiler_dev
DB_HOST=localhost
DB_PORT=5432

# .env.testing  
DB_NAME=boiler_test
DB_HOST=localhost
DB_PORT=5432

# .env.production
DB_NAME=boiler_prod
DB_HOST=your-production-host
DB_PORT=5432
```

## Troubleshooting

### Common Issues

#### Database Connection Errors
```
Error: failed to connect to database
```
**Solutions**:
1. Verify PostgreSQL is running
2. Check database credentials in `.env`
3. Ensure database exists
4. Check network connectivity

#### Port Already in Use
```
Error: listen tcp :8080: bind: address already in use
```
**Solutions**:
1. Kill process using port: `lsof -ti:8080 | xargs kill -9`
2. Use different port in `.env`: `PORT=8081`

#### Module Import Errors
```
Error: package github.com/FeisalDy/nogo/internal/user not found
```
**Solutions**:
1. Run `go mod tidy`
2. Verify module path in `go.mod`
3. Check import paths in source files

### Debug Mode
Enable detailed logging by setting environment:
```bash
export GIN_MODE=debug
export LOG_LEVEL=debug
```

### Database Debugging
Enable SQL logging in GORM:
```go
DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### Performance Monitoring
Add timing middleware for development:
```go
func TimingMiddleware() gin.HandlerFunc {
    return gin.Logger()
}
```

## Best Practices

### Security
- Never commit secrets to version control
- Use environment variables for configuration
- Validate all input data
- Implement proper error handling

### Performance
- Use database indexes appropriately
- Implement connection pooling
- Cache frequently accessed data
- Monitor query performance

### Maintainability
- Keep functions small and focused
- Use meaningful variable names
- Write comprehensive tests
- Document complex business logic
- Follow consistent code formatting