# Go Boilerplate

A production-ready Go REST API boilerplate built with the Gin framework, PostgreSQL, and Clean Architecture principles.

## Features

- 🏗️ **Clean Architecture** - Organized by domains with clear separation of concerns
- 🚀 **Gin Framework** - Fast HTTP web framework
- 🗄️ **PostgreSQL + GORM** - Robust database integration with ORM
- 🔧 **Environment Configuration** - Easy setup with environment variables
- 🐳 **Docker Support** - Containerized deployment ready
- 🔥 **Hot Reloading** - Development server with Air
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

## Project Structure

```
boiler/
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
- `POST /users` - Create a new user
- `GET /users/:id` - Get user by ID

For detailed API documentation, see [API.md](docs/API.md).

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
```
This will create a binary file named `my-app` in the root directory.

### Deployment

To deploy the application, you can simply run the binary file:
```bash
./my-app
```

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