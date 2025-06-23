# Use the official Golang image as a builder
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Run go generate
RUN go generate ./...

# Build the Go application
RUN CGO_ENABLED=0 go build -o main ./cmd/main.go

# Use a minimal image for the runtime
FROM alpine:latest

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# An entry point to run the application
ENTRYPOINT ["./main"]
