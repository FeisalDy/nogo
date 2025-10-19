# Authentication System

This document explains the authentication system implemented in the application using JWT (JSON Web Tokens) and bcrypt password hashing.

## Overview

The authentication system provides:

- **Secure password hashing** using bcrypt
- **JWT-based authentication** with token expiration
- **Protected routes** using middleware
- **User registration and login** endpoints
- **Token validation and refresh** capabilities

## Components

### 1. Password Hashing (`internal/common/utils/password.go`)

Functions for secure password handling:

```go
// Hash a plain text password
hashedPassword, err := utils.HashPassword("mypassword123")

// Compare plain text with hashed password
isValid := utils.ComparePassword(hashedPassword, "mypassword123")
```

### 2. JWT Utilities (`internal/common/utils/jwt.go`)

Functions for JWT token management:

```go
// Generate a token
token, err := utils.GenerateToken(userID, email, username)

// Validate a token
claims, err := utils.ValidateToken(tokenString)

// Refresh a token
newToken, err := utils.RefreshToken(oldToken)
```

**Token Configuration:**

- Secret key: `JWTSecret` (should be loaded from environment variables in production)
- Expiration: `24 hours` (configurable via `TokenExpiration`)

### 3. Auth Middleware (`internal/common/middleware/auth.go`)

Protects routes by validating JWT tokens:

```go
// Require authentication
protected := router.Group("/")
protected.Use(middleware.AuthMiddleware())
{
    protected.GET("/me", handler.GetMe)
}

// Optional authentication
router.Use(middleware.OptionalAuthMiddleware())
```

Helper functions to get user info from context:

```go
userID, exists := middleware.GetUserID(c)
email, exists := middleware.GetUserEmail(c)
username, exists := middleware.GetUserUsername(c)
```

## API Endpoints

### Register User

**Endpoint:** `POST /api/users/register`

**Request Body:**

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepass123",
  "confirm_password": "securepass123"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "status": "active"
    }
  },
  "message": "Registration successful"
}
```

**Validation Rules:**

- `username`: Required
- `email`: Required, must be valid email format
- `password`: Required, minimum 8 characters
- `confirm_password`: Required, must match password

**Error Responses:**

_Email already exists (409 Conflict):_

```json
{
  "success": false,
  "error": {
    "code": "USER002",
    "message": "User already exists"
  }
}
```

_Password mismatch (400 Bad Request):_

```json
{
  "success": false,
  "error": {
    "code": "AUTH006",
    "message": "Password and confirm password do not match"
  }
}
```

### Login User

**Endpoint:** `POST /api/users/login`

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "status": "active"
    }
  },
  "message": "Login successful"
}
```

**Error Responses:**

_Invalid credentials (401 Unauthorized):_

```json
{
  "success": false,
  "error": {
    "code": "USER006",
    "message": "Invalid username or password"
  }
}
```

### Get Current User

**Endpoint:** `GET /api/users/me`

**Headers:**

```
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "avatar_url": "https://example.com/avatar.jpg",
    "bio": "Software developer",
    "status": "active"
  }
}
```

**Error Responses:**

_Missing token (401 Unauthorized):_

```json
{
  "success": false,
  "error": {
    "code": "AUTH003",
    "message": "Authentication token is missing"
  }
}
```

_Invalid token (401 Unauthorized):_

```json
{
  "success": false,
  "error": {
    "code": "AUTH001",
    "message": "Invalid authentication token"
  }
}
```

## Usage Examples

### cURL Examples

**Register:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepass123",
    "confirm_password": "securepass123"
  }'
```

**Login:**

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Get Current User:**

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <your-token-here>"
```

### JavaScript/Fetch Examples

**Register:**

```javascript
const response = await fetch("http://localhost:8080/api/users/register", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    username: "john_doe",
    email: "john@example.com",
    password: "securepass123",
    confirm_password: "securepass123",
  }),
});

const data = await response.json();
if (data.success) {
  // Store token
  localStorage.setItem("token", data.data.token);
  localStorage.setItem("user", JSON.stringify(data.data.user));
}
```

**Login:**

```javascript
const response = await fetch("http://localhost:8080/api/users/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    email: "john@example.com",
    password: "securepass123",
  }),
});

const data = await response.json();
if (data.success) {
  localStorage.setItem("token", data.data.token);
  localStorage.setItem("user", JSON.stringify(data.data.user));
}
```

**Make Authenticated Request:**

```javascript
const token = localStorage.getItem("token");

const response = await fetch("http://localhost:8080/api/users/me", {
  headers: {
    Authorization: `Bearer ${token}`,
  },
});

const data = await response.json();
if (!data.success && data.error.code === "AUTH001") {
  // Token expired or invalid - redirect to login
  window.location.href = "/login";
}
```

## Implementation in Handlers

### Protecting Routes

```go
package novel

import (
    "github.com/FeisalDy/nogo/internal/common/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
    // Public routes
    router.GET("/", handler.ListNovels)
    router.GET("/:id", handler.GetNovel)

    // Protected routes
    protected := router.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.POST("/", handler.CreateNovel)
        protected.PUT("/:id", handler.UpdateNovel)
        protected.DELETE("/:id", handler.DeleteNovel)
    }
}
```

### Using User Info in Handlers

```go
func (h *NovelHandler) CreateNovel(c *gin.Context) {
    // Get authenticated user ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
        return
    }

    // Use userID in business logic
    novel := &model.Novel{
        Title:  dto.Title,
        AuthorID: userID,
    }

    // ... rest of handler
}
```

## Security Best Practices

### 1. Environment Variables

**IMPORTANT:** In production, load JWT secret from environment variables:

```go
// config/config.go
type Config struct {
    JWTSecret string
    JWTExpiration time.Duration
}

func LoadConfig() *Config {
    return &Config{
        JWTSecret: os.Getenv("JWT_SECRET"),
        JWTExpiration: 24 * time.Hour,
    }
}
```

Then update `internal/common/utils/jwt.go`:

```go
var JWTSecret = []byte(config.Get().JWTSecret)
```

**.env file:**

```env
JWT_SECRET=your-very-secure-random-secret-key-here
JWT_EXPIRATION=24h
```

### 2. Password Requirements

Current validation:

- Minimum 8 characters

**Recommended enhancements:**

```go
type RegisterUserDTO struct {
    Username        string `json:"username" validate:"required,min=3,max=20"`
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=8,max=72"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
```

Add custom validation for password strength:

```go
// internal/common/utils/validation.go
func ValidatePasswordStrength(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    hasUpper := false
    hasLower := false
    hasNumber := false
    hasSpecial := false

    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    if !hasUpper || !hasLower || !hasNumber {
        return errors.New("password must contain uppercase, lowercase, and numbers")
    }

    return nil
}
```

### 3. Rate Limiting

Add rate limiting to prevent brute force attacks:

```go
// internal/common/middleware/rate_limit.go
import "github.com/ulule/limiter/v3"

func RateLimitMiddleware() gin.HandlerFunc {
    rate := limiter.Rate{
        Period: 1 * time.Hour,
        Limit:  5, // 5 attempts per hour
    }

    store := memory.NewStore()
    instance := limiter.New(store, rate)

    return func(c *gin.Context) {
        limiterCtx, err := instance.Get(c, c.ClientIP())
        if err != nil {
            utils.RespondWithAppError(c, errors.ErrInternalServer)
            c.Abort()
            return
        }

        if limiterCtx.Reached {
            utils.RespondWithAppError(c, errors.NewAppError(
                "AUTH009",
                "Too many login attempts. Please try again later.",
            ))
            c.Abort()
            return
        }

        c.Next()
    }
}
```

Apply to login route:

```go
router.POST("/login", middleware.RateLimitMiddleware(), userHandler.Login)
```

### 4. Token Refresh

Implement token refresh for better security:

```go
// Handler
func (h *UserHandler) RefreshToken(c *gin.Context) {
    oldToken := c.GetHeader("Authorization")
    oldToken = strings.TrimPrefix(oldToken, "Bearer ")

    newToken, err := utils.RefreshToken(oldToken)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrAuthInvalidToken)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, gin.H{"token": newToken})
}
```

### 5. HTTPS Only

In production, always use HTTPS to protect tokens in transit.

### 6. Token Blacklist

For logout functionality, implement token blacklist:

```go
// Store revoked tokens in Redis or database
type TokenBlacklist struct {
    Token     string
    RevokedAt time.Time
}

// Check if token is blacklisted in middleware
if IsTokenBlacklisted(tokenString) {
    utils.RespondWithAppError(c, errors.ErrAuthInvalidToken)
    c.Abort()
    return
}
```

## Testing

See [AUTH_TESTING.md](AUTH_TESTING.md) for comprehensive testing examples.

## Troubleshooting

### "Invalid authentication token"

1. Check if token is properly formatted in Authorization header: `Bearer <token>`
2. Verify token hasn't expired (24 hours by default)
3. Ensure JWT secret matches between token generation and validation

### "User already exists"

Email addresses must be unique. Use a different email or implement password reset.

### Password not hashing

Ensure you're using the `Register` method, not the deprecated `CreateUser` method.

## Migration from Old System

If you have existing users with plain text passwords, run a migration:

```go
// One-time migration script
func MigratePasswords() {
    var users []model.User
    database.DB.Find(&users)

    for _, user := range users {
        if user.Password != nil && !isHashed(*user.Password) {
            hashed, _ := utils.HashPassword(*user.Password)
            user.Password = &hashed
            database.DB.Save(&user)
        }
    }
}

func isHashed(password string) bool {
    return strings.HasPrefix(password, "$2a$") // bcrypt hash prefix
}
```
