# Testing the Error Handling System

This guide shows how to test the new error handling system with example requests.

## Setup

Make sure your server is running:

```bash
go run cmd/server/main.go
```

## Test Cases

### 1. Missing All Required Fields

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response (400 Bad Request):**

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
          "message": "email is required",
          "tag": "required",
          "value": ""
        },
        {
          "field": "password",
          "message": "password is required",
          "tag": "required",
          "value": ""
        },
        {
          "field": "confirm password",
          "message": "confirm password is required",
          "tag": "required",
          "value": ""
        }
      ],
      "summary": "username: username is required; email: email is required; password: password is required; confirm password: confirm password is required"
    }
  }
}
```

### 2. Invalid Email Format

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "not-an-email",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

**Expected Response (400 Bad Request):**

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        {
          "field": "email",
          "message": "email must be a valid email address",
          "tag": "email",
          "value": "not-an-email"
        }
      ],
      "summary": "email: email must be a valid email address"
    }
  }
}
```

### 3. Password Too Short

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "123",
    "confirm_password": "123"
  }'
```

**Expected Response (400 Bad Request):**

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        {
          "field": "password",
          "message": "password must be at least 8 characters long",
          "tag": "min",
          "value": "123"
        }
      ],
      "summary": "password: password must be at least 8 characters long"
    }
  }
}
```

### 4. Password Mismatch

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "confirm_password": "different_password"
  }'
```

**Expected Response (400 Bad Request):**

```json
{
  "success": false,
  "error": {
    "code": "AUTH006",
    "message": "Password and confirm password do not match"
  }
}
```

### 5. Successful Registration

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

**Expected Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2025-10-19T10:30:00Z"
  },
  "message": "User created successfully"
}
```

### 6. User Not Found

**Request:**

```bash
curl -X GET http://localhost:8080/api/users/nonexistent-id
```

**Expected Response (404 Not Found):**

```json
{
  "success": false,
  "error": {
    "code": "USER001",
    "message": "User not found"
  }
}
```

### 7. Multiple Validation Errors

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "invalid-email",
    "password": "123",
    "confirm_password": "456"
  }'
```

**Expected Response (400 Bad Request):**

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
          "value": "456"
        }
      ],
      "summary": "username: username is required; email: email must be a valid email address; password: password must be at least 8 characters long; confirm password: confirm password must match password"
    }
  }
}
```

## Testing with Postman

### Collection Structure

```
NoGo API
├── Users
│   ├── Register User (Success)
│   ├── Register User (Missing Fields)
│   ├── Register User (Invalid Email)
│   ├── Register User (Short Password)
│   ├── Register User (Password Mismatch)
│   └── Get User (Not Found)
└── ...
```

### Test Scripts

Add these to Postman's "Tests" tab to verify responses:

#### For Validation Error Tests:

```javascript
pm.test("Status code is 400", function () {
  pm.response.to.have.status(400);
});

pm.test("Response has error structure", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData).to.have.property("success", false);
  pm.expect(jsonData).to.have.property("error");
  pm.expect(jsonData.error).to.have.property("code");
  pm.expect(jsonData.error).to.have.property("message");
});

pm.test("Error code is USER007", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData.error.code).to.eql("USER007");
});

pm.test("Has validation details", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData.error.details).to.have.property("fields");
  pm.expect(jsonData.error.details.fields).to.be.an("array");
});
```

#### For Success Tests:

```javascript
pm.test("Status code is 201", function () {
  pm.response.to.have.status(201);
});

pm.test("Response has success structure", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData).to.have.property("success", true);
  pm.expect(jsonData).to.have.property("data");
  pm.expect(jsonData).to.have.property("message");
});

pm.test("User data is present", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData.data).to.have.property("id");
  pm.expect(jsonData.data).to.have.property("username");
  pm.expect(jsonData.data).to.have.property("email");
});
```

## Integration Tests (Go)

Create a test file `internal/user/handler/user_handler_test.go`:

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/FeisalDy/nogo/internal/common/utils"
    "github.com/FeisalDy/nogo/internal/user/dto"
    "github.com/FeisalDy/nogo/internal/user/handler"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser_ValidationErrors(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name           string
        payload        dto.RegisterUserDTO
        expectedCode   string
        expectedStatus int
    }{
        {
            name: "Missing all fields",
            payload: dto.RegisterUserDTO{},
            expectedCode: "USER007",
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "Invalid email",
            payload: dto.RegisterUserDTO{
                Username: "john_doe",
                Email: "not-an-email",
                Password: "password123",
                ConfirmPassword: "password123",
            },
            expectedCode: "USER007",
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "Password too short",
            payload: dto.RegisterUserDTO{
                Username: "john_doe",
                Email: "john@example.com",
                Password: "123",
                ConfirmPassword: "123",
            },
            expectedCode: "USER007",
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            body, _ := json.Marshal(tt.payload)
            c.Request = httptest.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
            c.Request.Header.Set("Content-Type", "application/json")

            // Create handler and execute
            // h := handler.NewUserHandler(mockService)
            // h.CreateUser(c)

            // Assertions
            assert.Equal(t, tt.expectedStatus, w.Code)

            var response utils.ErrorResponse
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.False(t, response.Success)
            assert.Equal(t, tt.expectedCode, response.Error.Code)
        })
    }
}
```

## Automated Testing Script

Create `scripts/test_api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing Error Handling System"
echo "=============================="

# Test 1: Missing all fields
echo -e "\n${GREEN}Test 1: Missing all required fields${NC}"
curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.'

# Test 2: Invalid email
echo -e "\n${GREEN}Test 2: Invalid email format${NC}"
curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "not-an-email",
    "password": "password123",
    "confirm_password": "password123"
  }' | jq '.'

# Test 3: Password too short
echo -e "\n${GREEN}Test 3: Password too short${NC}"
curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "123",
    "confirm_password": "123"
  }' | jq '.'

# Test 4: Password mismatch
echo -e "\n${GREEN}Test 4: Password mismatch${NC}"
curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "confirm_password": "different_password"
  }' | jq '.'

echo -e "\n${GREEN}Tests completed!${NC}"
```

Make it executable:

```bash
chmod +x scripts/test_api.sh
```

Run tests:

```bash
./scripts/test_api.sh
```
