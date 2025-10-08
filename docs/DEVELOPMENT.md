# Development Guide

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Setup Instructions](#setup-instructions)
3. [Development Workflow](#development-workflow)
4. [Code Standards](#code-standards)
5. [Adding New Features](#adding-new-features)
6. [Testing Guidelines](#testing-guidelines)
7. [Database Management](#database-management)
8. [Troubleshooting](#troubleshooting)

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

### 5. Verify Setup
Test the API:
```bash
curl http://localhost:8080/ping
# Expected response: {"message":"pong"}
```

## Development Workflow

### 1. Starting Development
```bash
# Pull latest changes
git pull origin main

# Create feature branch
git checkout -b feature/your-feature-name

# Start development server
air
```

### 2. Making Changes
- Follow the established project structure
- Write tests for new functionality
- Update documentation as needed
- Run tests before committing

### 3. Committing Changes
```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add user profile endpoint"

# Push to remote branch
git push origin feature/your-feature-name
```

### 4. Code Review Process
1. Create Pull Request
2. Ensure all tests pass
3. Request code review
4. Address feedback
5. Merge to main branch

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

## Adding New Features

### 1. Planning
- Define the feature requirements
- Design the API endpoints
- Plan database schema changes
- Consider impact on existing code

### 2. Domain Structure
Create the complete domain structure:

```bash
mkdir -p internal/newdomain/{dto,handler,model,repository,service}
```

### 3. Implementation Order
1. **Model**: Define the domain entity
2. **Repository**: Implement data access
3. **Service**: Add business logic
4. **Handler**: Create HTTP endpoints
5. **DTO**: Define request/response structures

### 4. Integration
Update `main.go` to wire up the new domain:

```go
// Initialize domain
newRepository := newRepository.NewRepository()
newService := newService.NewService(newRepository)
newHandler := newHandler.NewHandler(newService)

// Register routes
newRoutes := r.Group("/newdomain")
{
    newRoutes.POST("/", newHandler.Create)
    newRoutes.GET("/:id", newHandler.Get)
}
```

### Example: Adding Product Domain

1. **Create Structure**:
```bash
mkdir -p internal/product/{dto,handler,model,repository,service}
```

2. **Model** (`internal/product/model/product.go`):
```go
package model

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name        string  `json:"name" gorm:"not null"`
    Description string  `json:"description"`
    Price       float64 `json:"price" gorm:"not null"`
    UserID      uint    `json:"user_id" gorm:"not null"`
}
```

3. **Repository** (`internal/product/repository/product_repository.go`):
```go
package repository

import (
    "boiler/internal/database"
    "boiler/internal/product/model"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
    return &ProductRepository{}
}

func (r *ProductRepository) CreateProduct(product *model.Product) error {
    return database.DB.Create(product).Error
}

func (r *ProductRepository) GetProduct(id string) (*model.Product, error) {
    var product model.Product
    if err := database.DB.First(&product, id).Error; err != nil {
        return nil, err
    }
    return &product, nil
}
```

4. **Service** (`internal/product/service/product_service.go`):
```go
package service

import (
    "boiler/internal/product/model"
    "boiler/internal/product/repository"
)

type ProductService struct {
    ProductRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
    return &ProductService{ProductRepository: productRepository}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
    // Add business logic here (validation, etc.)
    return s.ProductRepository.CreateProduct(product)
}

func (s *ProductService) GetProduct(id string) (*model.Product, error) {
    return s.ProductRepository.GetProduct(id)
}
```

5. **Handler** (`internal/product/handler/product_handler.go`):
```go
package handler

import (
    "net/http"
    "boiler/internal/product/model"
    "boiler/internal/product/service"
    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
    return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var product model.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.ProductService.CreateProduct(&product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
    id := c.Param("id")

    product, err := h.ProductService.GetProduct(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}
```

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

### Migrations
Currently using GORM AutoMigrate. For production, implement proper migrations:

```go
// In main.go or separate migration file
func migrate() {
    database.DB.AutoMigrate(&model.User{})
    database.DB.AutoMigrate(&model.Product{})
}
```

### Database Seeding
Create seed data for development:

```go
func seedDatabase() {
    users := []model.User{
        {Name: "John Doe", Email: "john@example.com"},
        {Name: "Jane Smith", Email: "jane@example.com"},
    }
    
    for _, user := range users {
        database.DB.FirstOrCreate(&user, model.User{Email: user.Email})
    }
}
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
Error: package boiler/internal/user not found
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