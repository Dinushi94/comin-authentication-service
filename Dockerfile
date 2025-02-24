# Stage 1: Build the application
FROM golang:1.23.4-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application from the correct directory
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/

# Stage 2: Create the final minimal image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application will run on
EXPOSE 8080

# Run the application
CMD ["./main"]