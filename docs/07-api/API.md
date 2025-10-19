# API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, no authentication is implemented. This is planned for future releases.

## Content Type
All requests and responses use `application/json` content type unless specified otherwise.

## Error Response Format
All error responses follow this structure:
```json
{
  "error": "Error description"
}
```

## Status Codes
- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Health Check

### Check API Health
**Endpoint**: `GET /ping`

**Description**: Check if the API is running and responsive.

**Request**: No parameters required

**Response**:
```json
{
  "message": "pong"
}
```

**Example**:
```bash
curl -X GET http://localhost:8080/ping
```

---

## Users API

### Create User
**Endpoint**: `POST /users`

**Description**: Create a new user in the system.

**Request Body**:
```json
{
  "name": "string",     // Required: User's full name
  "email": "string"     // Required: User's email address (must be unique)
}
```

**Success Response** (201 Created):
```json
{
  "ID": 1,
  "CreatedAt": "2025-10-06T10:00:00Z",
  "UpdatedAt": "2025-10-06T10:00:00Z",
  "DeletedAt": null,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error Responses**:

*400 Bad Request* - Invalid JSON or missing required fields:
```json
{
  "error": "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

*500 Internal Server Error* - Database error (e.g., duplicate email):
```json
{
  "error": "UNIQUE constraint failed: users.email"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### Get User by ID
**Endpoint**: `GET /users/{id}`

**Description**: Retrieve a user by their unique ID.

**Path Parameters**:
- `id` (integer): The unique identifier of the user

**Success Response** (200 OK):
```json
{
  "ID": 1,
  "CreatedAt": "2025-10-06T10:00:00Z",
  "UpdatedAt": "2025-10-06T10:00:00Z",
  "DeletedAt": null,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error Response**:

*404 Not Found* - User does not exist:
```json
{
  "error": "user not found"
}
```

**Example**:
```bash
curl -X GET http://localhost:8080/users/1
```

---

## Data Models

### User Model
```json
{
  "ID": "integer",           // Auto-generated unique identifier
  "CreatedAt": "datetime",   // Timestamp when user was created
  "UpdatedAt": "datetime",   // Timestamp when user was last updated
  "DeletedAt": "datetime",   // Timestamp when user was soft-deleted (null if active)
  "name": "string",          // User's full name
  "email": "string"          // User's email address (unique constraint)
}
```

### Field Constraints
- **name**: Required, string, no length limit currently
- **email**: Required, string, must be unique across all users
- **ID**: Auto-generated, primary key
- **CreatedAt/UpdatedAt**: Auto-managed by GORM
- **DeletedAt**: Used for soft deletion (GORM feature)

---

## Future API Endpoints

The following endpoints are planned for future releases:

### User Management
- `PUT /users/{id}` - Update user information
- `DELETE /users/{id}` - Delete user (soft delete)
- `GET /users` - List all users with pagination
- `GET /users/search?q={query}` - Search users by name or email

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh authentication token

### User Profile
- `GET /users/{id}/profile` - Get user profile details
- `PUT /users/{id}/profile` - Update user profile
- `POST /users/{id}/avatar` - Upload user avatar

---

## Testing the API

### Using cURL

1. **Test API Health**:
```bash
curl http://localhost:8080/ping
```

2. **Create a User**:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Smith", "email": "alice@example.com"}'
```

3. **Get User by ID**:
```bash
curl http://localhost:8080/users/1
```

### Using HTTPie

1. **Test API Health**:
```bash
http GET localhost:8080/ping
```

2. **Create a User**:
```bash
http POST localhost:8080/users name="Bob Johnson" email="bob@example.com"
```

3. **Get User by ID**:
```bash
http GET localhost:8080/users/1
```

### Using Postman

Import the following collection to test all endpoints:

```json
{
  "info": {
    "name": "Go Boilerplate API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{baseUrl}}/ping",
          "host": ["{{baseUrl}}"],
          "path": ["ping"]
        }
      }
    },
    {
      "name": "Create User",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"John Doe\",\n  \"email\": \"john@example.com\"\n}"
        },
        "url": {
          "raw": "{{baseUrl}}/users",
          "host": ["{{baseUrl}}"],
          "path": ["users"]
        }
      }
    },
    {
      "name": "Get User",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{baseUrl}}/users/1",
          "host": ["{{baseUrl}}"],
          "path": ["users", "1"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080",
      "type": "string"
    }
  ]
}
```

---

## Rate Limiting
Currently not implemented. In production, consider implementing rate limiting to prevent abuse.

## CORS
Currently not configured. For web applications, configure CORS middleware as needed.

## Monitoring
Consider implementing:
- Request logging
- Performance metrics
- Error tracking
- Health check endpoints for infrastructure monitoring