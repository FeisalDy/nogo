# Authentication Testing Guide

This guide provides test cases and examples for testing the authentication system.

## Test Scenarios

### 1. User Registration

#### Success Case

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

**Expected Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "status": "active"
    }
  },
  "message": "Registration successful"
}
```

#### Error Cases

**Missing Fields:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{}'
```

Response (400):

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        { "field": "username", "message": "username is required" },
        { "field": "email", "message": "email is required" },
        { "field": "password", "message": "password is required" },
        {
          "field": "confirm password",
          "message": "confirm password is required"
        }
      ]
    }
  }
}
```

**Invalid Email:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "invalid-email",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

Response (400):

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        { "field": "email", "message": "email must be a valid email address" }
      ]
    }
  }
}
```

**Password Too Short:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123",
    "confirm_password": "123"
  }'
```

Response (400):

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
          "message": "password must be at least 8 characters long"
        }
      ]
    }
  }
}
```

**Password Mismatch:**

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "different123"
  }'
```

Response (400):

```json
{
  "success": false,
  "error": {
    "code": "AUTH006",
    "message": "Password and confirm password do not match"
  }
}
```

**Duplicate Email:**

```bash
# Register first user
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'

# Try to register with same email
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user2",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

Response (409):

```json
{
  "success": false,
  "error": {
    "code": "USER002",
    "message": "User already exists"
  }
}
```

### 2. User Login

#### Success Case

**Request:**

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Expected Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "status": "active"
    }
  },
  "message": "Login successful"
}
```

#### Error Cases

**Invalid Credentials:**

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }'
```

Response (401):

```json
{
  "success": false,
  "error": {
    "code": "USER006",
    "message": "Invalid username or password"
  }
}
```

**Non-existent User:**

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@example.com",
    "password": "password123"
  }'
```

Response (401):

```json
{
  "success": false,
  "error": {
    "code": "USER006",
    "message": "Invalid username or password"
  }
}
```

### 3. Get Current User (Protected Route)

#### Success Case

**Request:**

```bash
# First, login to get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  | jq -r '.data.token')

# Then use token to access protected route
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": "active"
  }
}
```

#### Error Cases

**Missing Token:**

```bash
curl -X GET http://localhost:8080/api/users/me
```

Response (401):

```json
{
  "success": false,
  "error": {
    "code": "AUTH003",
    "message": "Authentication token is missing"
  }
}
```

**Invalid Token Format:**

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: invalid-token"
```

Response (401):

```json
{
  "success": false,
  "error": {
    "code": "AUTH001",
    "message": "Invalid authentication token"
  }
}
```

**Invalid Token:**

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer invalid.token.here"
```

Response (401):

```json
{
  "success": false,
  "error": {
    "code": "AUTH001",
    "message": "Invalid authentication token"
  }
}
```

## Test Script

Create `scripts/test_auth.sh`:

```bash
#!/usr/bin/env fish

set BASE_URL "http://localhost:8080"
set GREEN '\033[0;32m'
set RED '\033[0;31m'
set BLUE '\033[0;34m'
set NC '\033[0m' # No Color

echo -e "$BLUE=== Authentication Testing ===$NC\n"

# Test 1: Register User
echo -e "$GREEN[TEST 1] Register New User$NC"
set REGISTER_RESPONSE (curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_'(date +%s)'",
    "email": "test_'(date +%s)'@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }')

echo $REGISTER_RESPONSE | jq '.'

if test (echo $REGISTER_RESPONSE | jq -r '.success') = "true"
    echo -e "$GREEN✓ Registration successful$NC\n"
    set TOKEN (echo $REGISTER_RESPONSE | jq -r '.data.token')
    set USER_EMAIL (echo $REGISTER_RESPONSE | jq -r '.data.user.email')
else
    echo -e "$RED✗ Registration failed$NC\n"
    exit 1
end

# Test 2: Login with Created User
echo -e "$GREEN[TEST 2] Login with Registered User$NC"
set LOGIN_RESPONSE (curl -s -X POST "$BASE_URL/api/users/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$USER_EMAIL'",
    "password": "password123"
  }')

echo $LOGIN_RESPONSE | jq '.'

if test (echo $LOGIN_RESPONSE | jq -r '.success') = "true"
    echo -e "$GREEN✓ Login successful$NC\n"
    set TOKEN (echo $LOGIN_RESPONSE | jq -r '.data.token')
else
    echo -e "$RED✗ Login failed$NC\n"
    exit 1
end

# Test 3: Access Protected Route
echo -e "$GREEN[TEST 3] Access Protected Route (/me)$NC"
set ME_RESPONSE (curl -s -X GET "$BASE_URL/api/users/me" \
  -H "Authorization: Bearer $TOKEN")

echo $ME_RESPONSE | jq '.'

if test (echo $ME_RESPONSE | jq -r '.success') = "true"
    echo -e "$GREEN✓ Protected route access successful$NC\n"
else
    echo -e "$RED✗ Protected route access failed$NC\n"
    exit 1
end

# Test 4: Invalid Credentials
echo -e "$GREEN[TEST 4] Login with Invalid Credentials$NC"
set INVALID_LOGIN (curl -s -X POST "$BASE_URL/api/users/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$USER_EMAIL'",
    "password": "wrongpassword"
  }')

echo $INVALID_LOGIN | jq '.'

if test (echo $INVALID_LOGIN | jq -r '.success') = "false"
    echo -e "$GREEN✓ Correctly rejected invalid credentials$NC\n"
else
    echo -e "$RED✗ Should have rejected invalid credentials$NC\n"
end

# Test 5: Missing Token
echo -e "$GREEN[TEST 5] Access Protected Route Without Token$NC"
set NO_TOKEN (curl -s -X GET "$BASE_URL/api/users/me")

echo $NO_TOKEN | jq '.'

if test (echo $NO_TOKEN | jq -r '.success') = "false"
    echo -e "$GREEN✓ Correctly rejected missing token$NC\n"
else
    echo -e "$RED✗ Should have rejected missing token$NC\n"
end

# Test 6: Validation Error
echo -e "$GREEN[TEST 6] Registration with Validation Errors$NC"
set VALIDATION_ERROR (curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "invalid-email",
    "password": "123",
    "confirm_password": "456"
  }')

echo $VALIDATION_ERROR | jq '.'

if test (echo $VALIDATION_ERROR | jq -r '.error.code') = "USER007"
    echo -e "$GREEN✓ Correctly returned validation errors$NC\n"
else
    echo -e "$RED✗ Should have returned validation errors$NC\n"
end

echo -e "$BLUE=== All Tests Completed ===$NC"
```

Make it executable:

```bash
chmod +x scripts/test_auth.sh
```

Run tests:

```bash
./scripts/test_auth.sh
```

## Postman Collection

### Setup

1. Create a new collection called "NoGo API - Auth"
2. Add environment variables:
   - `base_url`: `http://localhost:8080`
   - `token`: (will be set automatically)

### Requests

#### 1. Register User

**Request:**

- Method: POST
- URL: `{{base_url}}/api/users/register`
- Body (JSON):

```json
{
  "username": "{{$randomUserName}}",
  "email": "{{$randomEmail}}",
  "password": "password123",
  "confirm_password": "password123"
}
```

**Tests Script:**

```javascript
pm.test("Status code is 201", () => {
  pm.response.to.have.status(201);
});

pm.test("Response has token", () => {
  const data = pm.response.json();
  pm.expect(data.data).to.have.property("token");
  pm.environment.set("token", data.data.token);
  pm.environment.set("user_email", data.data.user.email);
});
```

#### 2. Login User

**Request:**

- Method: POST
- URL: `{{base_url}}/api/users/login`
- Body (JSON):

```json
{
  "email": "{{user_email}}",
  "password": "password123"
}
```

**Tests Script:**

```javascript
pm.test("Status code is 200", () => {
  pm.response.to.have.status(200);
});

pm.test("Response has token", () => {
  const data = pm.response.json();
  pm.expect(data.data).to.have.property("token");
  pm.environment.set("token", data.data.token);
});
```

#### 3. Get Current User

**Request:**

- Method: GET
- URL: `{{base_url}}/api/users/me`
- Headers:
  - `Authorization`: `Bearer {{token}}`

**Tests Script:**

```javascript
pm.test("Status code is 200", () => {
  pm.response.to.have.status(200);
});

pm.test("Response has user data", () => {
  const data = pm.response.json();
  pm.expect(data.data).to.have.property("id");
  pm.expect(data.data).to.have.property("email");
});
```

## Integration Tests (Go)

Create `internal/user/handler/auth_test.go`:

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/FeisalDy/nogo/internal/user/dto"
    "github.com/FeisalDy/nogo/internal/user/handler"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name           string
        payload        dto.RegisterUserDTO
        expectedStatus int
        expectError    bool
    }{
        {
            name: "Successful registration",
            payload: dto.RegisterUserDTO{
                Username:        "testuser",
                Email:           "test@example.com",
                Password:        "password123",
                ConfirmPassword: "password123",
            },
            expectedStatus: http.StatusCreated,
            expectError:    false,
        },
        {
            name: "Password too short",
            payload: dto.RegisterUserDTO{
                Username:        "testuser",
                Email:           "test@example.com",
                Password:        "123",
                ConfirmPassword: "123",
            },
            expectedStatus: http.StatusBadRequest,
            expectError:    true,
        },
        {
            name: "Password mismatch",
            payload: dto.RegisterUserDTO{
                Username:        "testuser",
                Email:           "test@example.com",
                Password:        "password123",
                ConfirmPassword: "different123",
            },
            expectedStatus: http.StatusBadRequest,
            expectError:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            body, _ := json.Marshal(tt.payload)
            c.Request = httptest.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
            c.Request.Header.Set("Content-Type", "application/json")

            // Execute handler
            // handler.Register(c)

            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}
```

Run tests:

```bash
go test ./internal/user/handler/... -v
```

## Performance Testing

Use Apache Bench to test authentication performance:

```bash
# Test registration endpoint
ab -n 100 -c 10 -p register.json -T application/json \
  http://localhost:8080/api/users/register

# Test login endpoint
ab -n 1000 -c 50 -p login.json -T application/json \
  http://localhost:8080/api/users/login
```

## Security Testing

### 1. Test JWT Expiration

```bash
# Generate token and wait 24+ hours, then try to use it
# Should return AUTH002 error
```

### 2. Test Token Manipulation

```bash
# Modify token and try to use it
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.modified.signature"
```

Should return AUTH001 error.

### 3. Test SQL Injection

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password\" OR 1=1--"
  }'
```

Should safely handle and return invalid credentials.

## Troubleshooting Tests

If tests fail, check:

1. **Database is running** and accessible
2. **Server is running** on the correct port
3. **Email uniqueness** - use dynamic emails in tests
4. **Token expiration** - generate fresh tokens for each test
5. **JWT secret** - ensure it's consistent

## Continuous Integration

Add to `.github/workflows/test.yml`:

```yaml
name: Test Authentication

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:12
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21

      - name: Run tests
        run: go test ./... -v
        env:
          DB_HOST: localhost
          DB_USER: postgres
          DB_PASSWORD: postgres
          DB_NAME: test_db
```
