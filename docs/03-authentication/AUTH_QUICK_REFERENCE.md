# Authentication Quick Reference

## 🚀 Quick Start

### 1. Register a User

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

### 2. Login

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

Save the token from response:

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Access Protected Routes

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

## 📝 Code Examples

### Hash Password Before Saving

```go
import "github.com/FeisalDy/nogo/internal/common/utils"

// Hash password
hashedPassword, err := utils.HashPassword("plaintextpassword")
if err != nil {
    // handle error
}

// Save to database
user.Password = &hashedPassword
```

### Verify Password

```go
// Compare during login
isValid := utils.ComparePassword(user.Password, loginPassword)
if !isValid {
    // Invalid credentials
}
```

### Generate JWT Token

```go
import "github.com/FeisalDy/nogo/internal/common/utils"

token, err := utils.GenerateToken(user.ID, user.Email, username)
if err != nil {
    // handle error
}

// Return token to client
```

### Validate JWT Token

```go
claims, err := utils.ValidateToken(tokenString)
if err != nil {
    // Invalid token
}

// Access claims
userID := claims.UserID
email := claims.Email
```

### Protect Routes

```go
import "github.com/FeisalDy/nogo/internal/common/middleware"

// Protected route group
protected := router.Group("/")
protected.Use(middleware.AuthMiddleware())
{
    protected.GET("/me", handler.GetMe)
    protected.PUT("/profile", handler.UpdateProfile)
}
```

### Get User from Context

```go
import "github.com/FeisalDy/nogo/internal/common/middleware"

func (h *Handler) YourHandler(c *gin.Context) {
    // Get user ID
    userID, exists := middleware.GetUserID(c)
    if !exists {
        utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
        return
    }

    // Use userID in your logic
    resource.OwnerID = userID
}
```

### Complete Handler with Auth

```go
package handler

import (
    "net/http"

    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/middleware"
    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

func (h *Handler) CreateResource(c *gin.Context) {
    // 1. Get authenticated user
    userID, exists := middleware.GetUserID(c)
    if !exists {
        utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
        return
    }

    // 2. Validate input
    var dto CreateDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeValidation)
        return
    }

    validate := validator.New()
    if err := validate.Struct(dto); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeValidation)
        return
    }

    // 3. Create resource with authenticated user
    resource := &model.Resource{
        Title:   dto.Title,
        OwnerID: userID,
    }

    if err := h.service.Create(resource); err != nil {
        utils.RespondWithAppError(c, errors.ErrCreationFailed)
        return
    }

    // 4. Success
    utils.RespondSuccess(c, http.StatusCreated, resource, "Resource created")
}
```

## 🔑 Environment Configuration

**IMPORTANT:** Change JWT secret in production!

Create or update `.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database

# JWT Configuration
JWT_SECRET=your-very-secure-random-secret-key-change-this
JWT_EXPIRATION=24h

# Server
PORT=8080
```

Then update `internal/common/utils/jwt.go`:

```go
import (
    "os"
    "time"
)

var (
    JWTSecret = []byte(getEnv("JWT_SECRET", "default-secret-key"))
    TokenExpiration = 24 * time.Hour // or parse from JWT_EXPIRATION
)

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

## 📊 Response Formats

### Successful Registration/Login

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

### Error Responses

| Error               | Code    | Status | Description                  |
| ------------------- | ------- | ------ | ---------------------------- |
| Token missing       | AUTH003 | 401    | No Authorization header      |
| Invalid token       | AUTH001 | 401    | Malformed or expired token   |
| Token expired       | AUTH002 | 401    | Token past expiration        |
| Unauthorized        | AUTH004 | 401    | No permission                |
| Invalid credentials | USER006 | 401    | Wrong email/password         |
| User exists         | USER002 | 409    | Email already registered     |
| Password mismatch   | AUTH006 | 400    | Password != confirm_password |

## 🧪 Testing

### Fish Shell Script

```fish
#!/usr/bin/env fish

# Register
set RESPONSE (curl -s -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }')

# Extract token
set TOKEN (echo $RESPONSE | jq -r '.data.token')

# Use token
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### JavaScript/Fetch

```javascript
// Register
const registerResponse = await fetch("/api/users/register", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    username: "john_doe",
    email: "john@example.com",
    password: "securepass123",
    confirm_password: "securepass123",
  }),
});

const { data } = await registerResponse.json();

// Store token
localStorage.setItem("token", data.token);

// Use token in requests
const response = await fetch("/api/users/me", {
  headers: {
    Authorization: `Bearer ${localStorage.getItem("token")}`,
  },
});
```

## 🛡️ Security Checklist

- [ ] Change JWT_SECRET in production
- [ ] Use HTTPS in production
- [ ] Set appropriate token expiration
- [ ] Implement rate limiting on login/register
- [ ] Add password strength requirements
- [ ] Implement token refresh mechanism
- [ ] Add logout/token blacklist
- [ ] Enable CORS properly
- [ ] Sanitize all inputs
- [ ] Log authentication failures
- [ ] Implement account lockout after failed attempts
- [ ] Add 2FA (optional)

## 🔍 Common Issues

### "Invalid authentication token"

**Cause:** Token expired, malformed, or wrong secret key

**Solution:**

1. Check token format: `Bearer <token>`
2. Verify token hasn't expired (24h default)
3. Ensure JWT_SECRET matches

### "User already exists"

**Cause:** Email already registered

**Solution:**

- Use different email
- Implement password reset flow

### Password not hashing

**Cause:** Using old `CreateUser` method

**Solution:**

- Use `Register` method which hashes passwords automatically
- Or manually hash before saving:
  ```go
  hashedPass, _ := utils.HashPassword(plainPassword)
  user.Password = &hashedPass
  ```

### Can't access protected route

**Cause:** Missing or invalid token

**Solution:**

1. Login to get valid token
2. Add `Authorization: Bearer <token>` header
3. Check token hasn't expired

## 📚 Full Documentation

- [Authentication Guide](AUTHENTICATION.md) - Complete documentation
- [Authentication Testing](AUTH_TESTING.md) - Test cases and examples
- [Error Handling Guide](ERROR_HANDLING.md) - Error handling system

## 💡 Tips

1. **Store tokens securely:**

   - Browser: `localStorage` or `httpOnly` cookies
   - Mobile: Secure storage (Keychain, KeyStore)

2. **Handle token expiration:**

   ```javascript
   if (response.status === 401) {
     // Token expired - redirect to login
     localStorage.removeItem("token");
     window.location.href = "/login";
   }
   ```

3. **Refresh tokens:**

   ```go
   newToken, err := utils.RefreshToken(oldToken)
   ```

4. **Optional auth:**

   ```go
   // Route accessible with or without token
   router.Use(middleware.OptionalAuthMiddleware())
   ```

5. **Test with different users:**
   ```bash
   # Create multiple test users
   for i in {1..5}; do
     curl -X POST http://localhost:8080/api/users/register \
       -H "Content-Type: application/json" \
       -d "{\"username\":\"user$i\",\"email\":\"user$i@test.com\",\"password\":\"pass123\",\"confirm_password\":\"pass123\"}"
   done
   ```
