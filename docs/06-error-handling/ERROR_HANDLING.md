# Error Handling System

This document explains the error handling system implemented in the application.

## Overview

The error handling system provides:

- **Unique error codes** for different domains (USER, AUTH, NOVEL, CHAPTER, UPLOAD, DB, VAL, GEN)
- **Human-readable error messages** with proper formatting
- **Validation error formatting** that converts validator errors into readable messages
- **Standardized API responses** for both success and error cases

## Directory Structure

```
internal/
  common/
    errors/
      errors.go          # Error definitions and formatting
    utils/
      response.go        # Response helpers
```

## Error Codes

Error codes follow the pattern: `DOMAIN###` where:

- **DOMAIN** is a 3-10 letter prefix (USER, AUTH, NOVEL, etc.)
- **###** is a 3-digit number (001-099 per domain)

### Error Code Ranges

| Domain  | Code Range            | Description                         |
| ------- | --------------------- | ----------------------------------- |
| USER    | USER001-USER099       | User-related errors                 |
| AUTH    | AUTH001-AUTH099       | Authentication/Authorization errors |
| NOVEL   | NOVEL001-NOVEL099     | Novel-related errors                |
| CHAPTER | CHAPTER001-CHAPTER099 | Chapter-related errors              |
| UPLOAD  | UPLOAD001-UPLOAD099   | File upload errors                  |
| DB      | DB001-DB099           | Database errors                     |
| VAL     | VAL001-VAL099         | Validation errors                   |
| GEN     | GEN001-GEN099         | General errors                      |

## Usage

### 1. Using Predefined Errors

```go
import (
    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/utils"
)

func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.UserService.GetUser(id)
    if err != nil {
        // Use predefined error
        utils.RespondWithAppError(c, errors.ErrUserNotFound)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, user)
}
```

### 2. Creating Custom Errors with Details

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    if err := h.UserService.CreateUser(&user); err != nil {
        // Add custom details to error
        appError := errors.ErrUserCreationFailed.WithDetails(map[string]any{
            "reason": err.Error(),
            "timestamp": time.Now(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    utils.RespondSuccess(c, http.StatusCreated, user, "User created successfully")
}
```

### 3. Handling Validation Errors

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var registerDTO dto.RegisterUserDTO
    if err := c.ShouldBindJSON(&registerDTO); err != nil {
        // Automatically formats validation errors
        utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
        return
    }

    validate := validator.New()
    if err := validate.Struct(registerDTO); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
        return
    }
}
```

### 4. Creating New Error Types

```go
// In errors.go, add new error code constant
const (
    ErrCodeUserSuspended = "USER008"
)

// Add predefined error
var (
    ErrUserSuspended = NewAppError(ErrCodeUserSuspended, "User account is suspended")
)

// In response.go, add status code mapping
func GetStatusCodeFromErrorCode(code string) int {
    switch {
    case code == errors.ErrCodeUserSuspended:
        return http.StatusForbidden
    // ... other cases
    }
}
```

## Response Formats

### Success Response

```json
{
  "success": true,
  "data": {
    "id": "123",
    "username": "john_doe",
    "email": "john@example.com"
  },
  "message": "User created successfully"
}
```

### Error Response (Simple)

```json
{
  "success": false,
  "error": {
    "code": "USER001",
    "message": "User not found"
  }
}
```

### Error Response (With Details)

```json
{
  "success": false,
  "error": {
    "code": "USER003",
    "message": "Failed to create user",
    "details": {
      "reason": "duplicate email address"
    }
  }
}
```

### Validation Error Response

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        {
          "field": "username",
          "message": "username is required",
          "tag": "required",
          "value": ""
        },
        {
          "field": "email",
          "message": "email must be a valid email address",
          "tag": "email",
          "value": "invalid-email"
        },
        {
          "field": "password",
          "message": "password must be at least 8 characters long",
          "tag": "min",
          "value": "123"
        },
        {
          "field": "confirm password",
          "message": "confirm password must match password",
          "tag": "eqfield",
          "value": ""
        }
      ],
      "summary": "username: username is required; email: email must be a valid email address; password: password must be at least 8 characters long; confirm password: confirm password must match password"
    }
  }
}
```

## Response Helper Functions

### `RespondSuccess`

Sends a successful response with optional message.

```go
utils.RespondSuccess(c, http.StatusOK, data)
utils.RespondSuccess(c, http.StatusCreated, data, "Created successfully")
```

### `RespondError`

Sends an error response with custom status code.

```go
utils.RespondError(c, http.StatusBadRequest, appError)
```

### `RespondValidationError`

Formats and sends validation errors (automatically uses 400 Bad Request).

```go
utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
```

### `RespondWithAppError`

Automatically determines the appropriate HTTP status code based on error code.

```go
utils.RespondWithAppError(c, errors.ErrUserNotFound) // Returns 404
utils.RespondWithAppError(c, errors.ErrAuthUnauthorized) // Returns 401
```

## Validation Messages

The system automatically converts validation tags to human-readable messages:

| Validation Tag | Example Message                                |
| -------------- | ---------------------------------------------- |
| `required`     | "username is required"                         |
| `email`        | "email must be a valid email address"          |
| `min`          | "password must be at least 8 characters long"  |
| `max`          | "username must be at most 20 characters long"  |
| `eqfield`      | "confirm password must match password"         |
| `len`          | "code must be exactly 6 characters long"       |
| `gte`          | "age must be greater than or equal to 18"      |
| `alpha`        | "name must contain only alphabetic characters" |
| `url`          | "website must be a valid URL"                  |

## Best Practices

1. **Always use predefined errors** when possible
2. **Use meaningful error codes** that identify the domain and type of error
3. **Add details** when you need to provide more context
4. **Use validation error helper** for struct validation
5. **Don't expose sensitive information** in error messages
6. **Keep error messages user-friendly** and actionable
7. **Use appropriate HTTP status codes** (automatically handled by `RespondWithAppError`)

## Example: Complete Handler with Error Handling

```go
package handler

import (
    "net/http"

    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/FeisalDy/nogo/internal/user/dto"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type UserHandler struct {
    UserService *service.UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    // 1. Bind and validate JSON
    var registerDTO dto.RegisterUserDTO
    if err := c.ShouldBindJSON(&registerDTO); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
        return
    }

    // 2. Validate struct
    validate := validator.New()
    if err := validate.Struct(registerDTO); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
        return
    }

    // 3. Business logic validation
    if registerDTO.Password != registerDTO.ConfirmPassword {
        utils.RespondWithAppError(c, errors.ErrAuthPasswordMismatch)
        return
    }

    // 4. Call service
    user, err := h.UserService.CreateUser(&registerDTO)
    if err != nil {
        // Handle different error types
        if err == service.ErrEmailExists {
            utils.RespondWithAppError(c, errors.ErrUserAlreadyExists)
            return
        }

        // Generic error with details
        appError := errors.ErrUserCreationFailed.WithDetails(map[string]interface{}{
            "reason": err.Error(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    // 5. Success response
    utils.RespondSuccess(c, http.StatusCreated, user, "User created successfully")
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")

    user, err := h.UserService.GetUser(id)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrUserNotFound)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    id := c.Param("id")

    var updateDTO dto.UpdateUserDTO
    if err := c.ShouldBindJSON(&updateDTO); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
        return
    }

    user, err := h.UserService.UpdateUser(id, &updateDTO)
    if err != nil {
        if err == service.ErrUserNotFound {
            utils.RespondWithAppError(c, errors.ErrUserNotFound)
            return
        }

        appError := errors.ErrUserUpdateFailed.WithDetails(map[string]interface{}{
            "reason": err.Error(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, user, "User updated successfully")
}
```

## Adding New Error Domains

When adding a new feature domain (e.g., comments, reviews), follow these steps:

1. **Add error code constants** in `errors.go`:

```go
const (
    ErrCodeCommentNotFound = "COMMENT001"
    ErrCodeCommentCreationFailed = "COMMENT002"
    // ...
)
```

2. **Add predefined errors**:

```go
var (
    ErrCommentNotFound = NewAppError(ErrCodeCommentNotFound, "Comment not found")
    ErrCommentCreationFailed = NewAppError(ErrCodeCommentCreationFailed, "Failed to create comment")
    // ...
)
```

3. **Add status code mapping** in `response.go`:

```go
case code == errors.ErrCodeCommentNotFound:
    return http.StatusNotFound
case code == errors.ErrCodeCommentCreationFailed:
    return http.StatusInternalServerError
```

4. **Use in handlers**:

```go
utils.RespondWithAppError(c, errors.ErrCommentNotFound)
```
