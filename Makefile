.PHONY: help build build-tools run test tidy vet swag docker-build docker-up docker-down migrate-up migrate-down migrate-version seed

APP ?= ./cmd/api
BINARY ?= bin/api

help:
	@echo "Targets:"
	@echo "  make build         - build API binary to $(BINARY)"
	@echo "  make build-tools   - build bin/api, bin/migrate, bin/seed"
	@echo "  make run           - go run ./cmd/api"
	@echo "  make test          - go test ./..."
	@echo "  make tidy          - go mod tidy"
	@echo "  make vet           - go vet ./..."
	@echo "  make swag          - regenerate Swagger docs under docs/"
	@echo "  make migrate-up    - go run ./cmd/migrate up"
	@echo "  make migrate-down  - go run ./cmd/migrate down"
	@echo "  make migrate-version - show applied migration version"
	@echo "  make seed          - go run ./cmd/seed (needs SEED_SUPER_ADMIN_* env)"
	@echo "  make docker-build  - docker compose build api"
	@echo "  make docker-up     - docker compose up --build"
	@echo "  make docker-down   - docker compose down"

build:
	@mkdir -p bin
	go build -o $(BINARY) $(APP)

build-tools:
	@mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/seed ./cmd/seed

run:
	go run $(APP)

test:
	go test ./...

tidy:
	go mod tidy

vet:
	go vet ./...

swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal

docker-build:
	docker compose build api

docker-up:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-version:
	go run ./cmd/migrate version

seed:
	go run ./cmd/seed
