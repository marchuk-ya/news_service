.PHONY: build run test docker-build docker-up docker-down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Build Docker image
docker-build:
	docker-compose -f docker/docker-compose.yml build

# Start Docker containers
docker-up:
	docker-compose -f docker/docker-compose.yml up -d

# Stop Docker containers
docker-down:
	docker-compose -f docker/docker-compose.yml down

# Install dependencies
deps:
	go mod download
	go mod tidy 