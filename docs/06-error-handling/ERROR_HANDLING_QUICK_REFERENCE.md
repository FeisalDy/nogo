# Error Handling Quick Reference

## Quick Copy-Paste Examples

### 1. Basic Handler Pattern

```go
import (
    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/gin-gonic/gin"
)

func (h *Handler) YourHandler(c *gin.Context) {
    // Validation
    var dto YourDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeYourDomainValidation)
        return
    }

    // Business logic
    result, err := h.Service.DoSomething(&dto)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrYourError)
        return
    }

    // Success
    utils.RespondSuccess(c, http.StatusOK, result)
}
```

### 2. Validation with Validator

```go
validate := validator.New()
if err := validate.Struct(dto); err != nil {
    utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
    return
}
```

### 3. Error with Custom Details

```go
if err != nil {
    appError := errors.ErrUserCreationFailed.WithDetails(map[string]interface{}{
        "reason": err.Error(),
        "user_id": userId,
    })
    utils.RespondWithAppError(c, appError)
    return
}
```

### 4. Success Response

```go
// Without message
utils.RespondSuccess(c, http.StatusOK, data)

// With message
utils.RespondSuccess(c, http.StatusCreated, data, "Resource created successfully")
```

## All Error Codes at a Glance

| Code                  | Description               | HTTP Status |
| --------------------- | ------------------------- | ----------- |
| **USER Domain**       |                           |             |
| USER001               | User not found            | 404         |
| USER002               | User already exists       | 409         |
| USER003               | User creation failed      | 500         |
| USER004               | User update failed        | 500         |
| USER005               | User deletion failed      | 500         |
| USER006               | Invalid credentials       | 401         |
| USER007               | User validation failed    | 400         |
| **AUTH Domain**       |                           |             |
| AUTH001               | Invalid token             | 401         |
| AUTH002               | Token expired             | 401         |
| AUTH003               | Token missing             | 401         |
| AUTH004               | Unauthorized              | 401         |
| AUTH005               | Forbidden                 | 403         |
| AUTH006               | Password mismatch         | 400         |
| AUTH007               | Registration failed       | 400         |
| AUTH008               | Login failed              | 401         |
| **NOVEL Domain**      |                           |             |
| NOVEL001              | Novel not found           | 404         |
| NOVEL002              | Novel creation failed     | 500         |
| NOVEL003              | Novel update failed       | 500         |
| NOVEL004              | Novel deletion failed     | 500         |
| NOVEL005              | Novel validation failed   | 400         |
| **CHAPTER Domain**    |                           |             |
| CHAPTER001            | Chapter not found         | 404         |
| CHAPTER002            | Chapter creation failed   | 500         |
| CHAPTER003            | Chapter update failed     | 500         |
| CHAPTER004            | Chapter deletion failed   | 500         |
| CHAPTER005            | Chapter validation failed | 400         |
| **UPLOAD Domain**     |                           |             |
| UPLOAD001             | Invalid file              | 400         |
| UPLOAD002             | Upload failed             | 500         |
| UPLOAD003             | File too large            | 400         |
| UPLOAD004             | Invalid file type         | 400         |
| UPLOAD005             | No file provided          | 400         |
| **DATABASE Domain**   |                           |             |
| DB001                 | Connection error          | 500         |
| DB002                 | Query error               | 500         |
| DB003                 | Transaction error         | 500         |
| **VALIDATION Domain** |                           |             |
| VAL001                | Validation failed         | 400         |
| VAL002                | Invalid input             | 400         |
| VAL003                | Missing field             | 400         |
| **GENERAL Domain**    |                           |             |
| GEN001                | Internal server error     | 500         |
| GEN002                | Bad request               | 400         |
| GEN003                | Not found                 | 404         |

## Response Helper Functions

```go
// Success responses
utils.RespondSuccess(c, statusCode, data)
utils.RespondSuccess(c, statusCode, data, "Optional message")

// Error responses
utils.RespondError(c, statusCode, appError)
utils.RespondValidationError(c, err, errorCode)
utils.RespondWithAppError(c, appError) // Auto-determines status code
```

## Common Validation Tags

```go
`validate:"required"`           // Field is required
`validate:"email"`              // Must be valid email
`validate:"min=8"`              // Minimum length
`validate:"max=100"`            // Maximum length
`validate:"eqfield=Password"`   // Must equal another field
`validate:"gte=0"`              // Greater than or equal
`validate:"lte=100"`            // Less than or equal
`validate:"oneof=active inactive"` // Must be one of values
`validate:"url"`                // Must be valid URL
`validate:"uuid"`               // Must be valid UUID
`validate:"alphanum"`           // Alphanumeric only
```

## Adding New Error Codes

### Step 1: Add constant in `errors/errors.go`

```go
const (
    ErrCodeYourNewError = "DOMAIN###"
)
```

### Step 2: Add predefined error

```go
var (
    ErrYourNewError = NewAppError(ErrCodeYourNewError, "Your error message")
)
```

### Step 3: Add status code mapping in `utils/response.go`

```go
case code == errors.ErrCodeYourNewError:
    return http.StatusBadRequest
```

### Step 4: Use in handler

```go
utils.RespondWithAppError(c, errors.ErrYourNewError)
```

## Example: Complete CRUD Handler

```go
package handler

import (
    "net/http"

    "github.com/FeisalDy/nogo/internal/common/errors"
    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type ResourceHandler struct {
    service *service.ResourceService
}

// Create
func (h *ResourceHandler) Create(c *gin.Context) {
    var dto CreateDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeResourceValidation)
        return
    }

    validate := validator.New()
    if err := validate.Struct(dto); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeResourceValidation)
        return
    }

    resource, err := h.service.Create(&dto)
    if err != nil {
        appError := errors.ErrResourceCreationFailed.WithDetails(map[string]interface{}{
            "reason": err.Error(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    utils.RespondSuccess(c, http.StatusCreated, resource, "Resource created successfully")
}

// Read
func (h *ResourceHandler) Get(c *gin.Context) {
    id := c.Param("id")

    resource, err := h.service.Get(id)
    if err != nil {
        utils.RespondWithAppError(c, errors.ErrResourceNotFound)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, resource)
}

// List
func (h *ResourceHandler) List(c *gin.Context) {
    resources, err := h.service.List()
    if err != nil {
        appError := errors.ErrDatabaseQuery.WithDetails(map[string]interface{}{
            "reason": err.Error(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, resources)
}

// Update
func (h *ResourceHandler) Update(c *gin.Context) {
    id := c.Param("id")

    var dto UpdateDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeResourceValidation)
        return
    }

    resource, err := h.service.Update(id, &dto)
    if err != nil {
        if err == service.ErrNotFound {
            utils.RespondWithAppError(c, errors.ErrResourceNotFound)
            return
        }

        appError := errors.ErrResourceUpdateFailed.WithDetails(map[string]interface{}{
            "reason": err.Error(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, resource, "Resource updated successfully")
}

// Delete
func (h *ResourceHandler) Delete(c *gin.Context) {
    id := c.Param("id")

    if err := h.service.Delete(id); err != nil {
        if err == service.ErrNotFound {
            utils.RespondWithAppError(c, errors.ErrResourceNotFound)
            return
        }

        appError := errors.ErrResourceDeletionFailed.WithDetails(map[string]interface{}{
            "reason": err.Error(),
        })
        utils.RespondWithAppError(c, appError)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, gin.H{"id": id}, "Resource deleted successfully")
}
```
