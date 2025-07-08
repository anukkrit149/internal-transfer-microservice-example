# Variables
APP_NAME := internal-transfer-microservice
DOCKER_REPO := anukkrit149/$(APP_NAME)
IMG_TAG := latest
VERSION := 1.0.0
GO_BUILD_FLAGS := -v

# Go related variables
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

# Docker related variables
DOCKER_BUILD_FLAGS := --no-cache

.PHONY: all build clean run migrate docker-build docker-push

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build $(GO_BUILD_FLAGS) -o $(APP_NAME) .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(APP_NAME)
	go clean

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	go run main.go api

# Run with custom config
run-with-config:
	@echo "Running $(APP_NAME) with custom config..."
	go run main.go api --config config/env.yaml

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run main.go migrate

# Run database migrations with custom config
migrate-with-config:
	@echo "Running database migrations with custom config..."
	go run main.go migrate --config config/env.yaml

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_REPO):$(VERSION)..."
	docker build $(DOCKER_BUILD_FLAGS) -t $(DOCKER_REPO):$(VERSION) .
	docker tag $(DOCKER_REPO):$(VERSION) $(DOCKER_REPO):latest

# Push Docker image to registry
docker-push:
	@echo "Pushing Docker image $(DOCKER_REPO):$(VERSION) to registry..."
	docker push $(DOCKER_REPO):$(VERSION)
	docker push $(DOCKER_REPO):latest

# Run the application in Docker
docker-run:
	@echo "Running $(APP_NAME) in Docker..."
	docker run -p 3000:3000 --name $(APP_NAME) $(DOCKER_REPO):$(IMG_TAG)

# Stop and remove Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# Help
help:
	@echo "Available targets:"
	@echo "  all               - Default target, builds the application"
	@echo "  build             - Build the application"
	@echo "  clean             - Clean build artifacts"
	@echo "  run               - Run the application"
	@echo "  run-with-config   - Run the application with custom config"
	@echo "  migrate           - Run database migrations"
	@echo "  migrate-with-config - Run database migrations with custom config"
	@echo "  test              - Run tests"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-push       - Push Docker image to registry"
	@echo "  docker-run        - Run the application in Docker"
	@echo "  docker-stop       - Stop and remove Docker container"
	@echo "  help              - Show this help message"