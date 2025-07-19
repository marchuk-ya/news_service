.PHONY: up down clean deps mock test-unit test-integration test-coverage build run swagger

DOCKER_COMPOSE = docker-compose -f docker/docker-compose.yml

up:
	$(DOCKER_COMPOSE) up -d --build

down:
	$(DOCKER_COMPOSE) down

clean:
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

deps:
	go mod download
	go mod tidy

mock:
	~/go/bin/mockery --dir internal/domain --name NewsRepository --output internal/mocks --outpkg mocks

swagger:
	~/go/bin/swag init -g cmd/main.go -o docs

test-unit:
	go test -v -coverprofile=coverage.out ./internal/... -short
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-integration:
	go test -v ./tests/...

test-coverage:
	go tool cover -html=coverage.out

build:
	go build -o bin/news_service cmd/main.go

run:
	go run cmd/main.go