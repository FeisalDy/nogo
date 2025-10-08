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
