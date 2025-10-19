# Go Boilerplate

A production-ready Go REST API boilerplate built with the Gin framework, PostgreSQL, and Clean Architecture principles.

## Features

- 🏗️ **Clean Architecture** - Organized by domains with clear separation of concerns
- 🚀 **Gin Framework** - Fast HTTP web framework
- 🗄️ **PostgreSQL + GORM** - Robust database integration with ORM
- � **JWT Authentication** - Secure token-based authentication with bcrypt password hashing
- 🛡️ **Auth Middleware** - Protected routes with automatic token validation
- �🔧 **Environment Configuration** - Easy setup with environment variables
- 🐳 **Docker Support** - Containerized deployment ready
- 🔥 **Hot Reloading** - Development server with Air
- ⚠️ **Advanced Error Handling** - Structured errors with unique codes and readable messages
- ✅ **Validation Error Formatting** - Beautiful, user-friendly validation error responses
- 📚 **Comprehensive Documentation** - Architecture, API, and development guides

## Quick Start

1. **Clone and Setup**

   ```bash
   git clone <repository-url>
   cd boiler
   go mod tidy
   ```

2. **Configure Environment**

   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Run with Hot Reloading**

   ```bash
   # Install Air (if not already installed)
   go install github.com/cosmtrek/air@latest

   # Start development server
   air
   ```

4. **Test the API**
   ```bash
   curl http://localhost:8080/ping
   # Response: {"message":"pong"}
   ```

## Documentation

- 📖 **[Architecture Guide](docs/ARCHITECTURE.md)** - Detailed explanation of the application structure and flow
- 🔗 **[API Documentation](docs/API.md)** - Complete API endpoint reference
- 👨‍💻 **[Development Guide](docs/DEVELOPMENT.md)** - Setup instructions and coding standards
- 🔐 **[Authentication Guide](docs/AUTHENTICATION.md)** - JWT authentication and password hashing
- 🧪 **[Authentication Testing](docs/AUTH_TESTING.md)** - Testing guide for auth endpoints
- ⚠️ **[Error Handling Guide](docs/ERROR_HANDLING.md)** - Complete error handling system documentation
- ⚡ **[Error Handling Quick Reference](docs/ERROR_HANDLING_QUICK_REFERENCE.md)** - Quick copy-paste examples
- 🧪 **[Error Handling Tests](docs/ERROR_HANDLING_TESTS.md)** - Testing guide for error responses

## Project Structure

```
github.com/FeisalDy/nogo/
├── cmd/server/           # Application entry point
├── config/              # Configuration management
├── docs/                # Documentation files
├── internal/            # Private application code
│   ├── common/          # Shared utilities and middleware
│   ├── database/        # Database connection
│   └── user/           # User domain (example)
│       ├── dto/        # Data transfer objects
│       ├── handler/    # HTTP handlers
│       ├── model/      # Domain models
│       ├── repository/ # Data access layer
│       └── service/    # Business logic
├── pkg/                # Public libraries
├── scripts/            # Helper scripts
├── Dockerfile          # Docker configuration
└── README.md          # This file
```

## API Endpoints

### Health Check

- `GET /ping` - API health check

### User Management

- `POST /api/users/register` - Register a new user (returns JWT token)
- `POST /api/users/login` - Login user (returns JWT token)
- `GET /api/users/me` - Get current authenticated user (requires token)
- `GET /api/users/:id` - Get user by ID (requires token)

For detailed API documentation, see [API.md](docs/API.md).

## Authentication

The application features JWT-based authentication with secure password hashing:

### Register a New User

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

Response includes JWT token:

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

### Login

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### Access Protected Routes

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Security Features:**

- Passwords hashed with bcrypt
- JWT tokens with 24-hour expiration
- Protected routes via middleware
- Validation on all inputs

For detailed authentication documentation, see:

- [Authentication Guide](docs/AUTHENTICATION.md) - Complete authentication system
- [Authentication Testing](docs/AUTH_TESTING.md) - Testing examples and scenarios

## Error Handling

The application features a comprehensive error handling system with:

- **Unique error codes** per domain (USER, AUTH, NOVEL, CHAPTER, etc.)
- **Human-readable error messages** with automatic validation error formatting
- **Standardized API responses** for both success and error cases
- **Automatic HTTP status code mapping** based on error type

### Example Error Response

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
          "value": "invalid-email"
        }
      ]
    }
  }
}
```

### Example Success Response

```json
{
  "success": true,
  "data": {
    "id": "123",
    "username": "john_doe"
  },
  "message": "User created successfully"
}
```

For detailed documentation, see:

- [Error Handling Guide](docs/ERROR_HANDLING.md) - Complete documentation
- [Quick Reference](docs/ERROR_HANDLING_QUICK_REFERENCE.md) - Copy-paste examples
- [Testing Guide](docs/ERROR_HANDLING_TESTS.md) - How to test error responses

## Development

See the [Development Guide](docs/DEVELOPMENT.md) for detailed setup instructions, coding standards, and best practices.

### Quick Development Setup

1. **Prerequisites**: Go 1.21+, PostgreSQL 12+
2. **Install dependencies**: `go mod tidy`
3. **Setup database**: Create PostgreSQL database and configure `.env`
4. **Install Air**: `go install github.com/cosmtrek/air@latest`
5. **Start development**: `air`

### Environment Variables

Create a `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=boiler_dev
PORT=8080
```

## Building and Deployment

### Building

To build the application, run the following command:
go build -o my-app cmd/server/main.go

````
This will create a binary file named `my-app` in the root directory.

### Deployment

To deploy the application, you can simply run the binary file:
```bash
./my-app
````

You can also use a process manager like `systemd` or `supervisor` to run the application in the background.

For a more robust deployment, you can use Docker to containerize the application. Here is an example `Dockerfile`:

```Dockerfile
# Start from the official Go image
FROM golang:1.21-alpine

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o my-app cmd/server/main.go

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./my-app"]
```

You can then build the Docker image and run the container:

```bash
docker build -t my-app .
docker run -p 8080:8080 my-app
```
