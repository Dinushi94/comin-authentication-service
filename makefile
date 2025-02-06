# Binary name
BINARY_NAME=comin-auth
DATABASE_URL="postgresql://comin_owner:Ye5rfjcIB7FX@ep-flat-shadow-a8onelva.eastus2.azure.neon.tech/comin?sslmode=require"

.PHONY: build run test clean deps tidy migrate-* docker-* test-setup test-teardown test-integration test-all

# Build commands
build:
	go build -o bin/$(BINARY_NAME) cmd/auth/main.go

run:
	go run cmd/auth/main.go

# Test commands
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-race:
	go test -race -v ./...

test-setup:
	docker-compose -f docker-compose.test.yml up -d
	sleep 5

test-teardown:
	docker-compose -f docker-compose.test.yml down -v

test-integration: test-setup
	go test -v ./... -tags=integration
	make test-teardown

test-db-up:
	docker-compose -f docker-compose.test.yml up -d

test-db-down:
	docker-compose -f docker-compose.test.yml down -v

test-all: test test-integration

# Clean up commands
clean:
	rm -f bin/$(BINARY_NAME)
	go clean

# Dependency management
deps:
	go mod download

tidy:
	go mod tidy

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Development tools
dev:
	air -c .air.toml

lint:
	golangci-lint run

swagger:
	swag init -g cmd/auth/main.go

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build files"
	@echo "  deps          - Download dependencies"
	@echo "  tidy          - Tidy go.mod"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  migrate-create- Create new migration file"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  dev          - Run with hot reload"
	@echo "  lint          - Run linter"
	@echo "  swagger       - Generate swagger documentation"