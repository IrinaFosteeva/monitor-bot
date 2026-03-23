ifneq (,$(wildcard .env))
include .env
export
endif

DB_SSLMODE ?= disable
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_DIR ?= ./migrations

.PHONY: help db-up db-down migrate-up migrate-down migrate-version test run-api run-bot run-worker

help:
	@echo "Available targets:"
	@echo "  make db-up            - start local infrastructure via docker-compose"
	@echo "  make db-down          - stop local infrastructure"
	@echo "  make migrate-up       - apply all database migrations"
	@echo "  make migrate-down     - rollback the last database migration"
	@echo "  make migrate-version  - show current migration version"
	@echo "  make test             - run Go tests"
	@echo "  make run-api          - run REST API"
	@echo "  make run-bot          - run Telegram bot"
	@echo "  make run-worker       - run monitoring worker"

db-up:
	docker-compose up -d

db-down:
	docker-compose down

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

test:
	go test ./...

run-api:
	go run ./cmd/api

run-bot:
	go run ./cmd/bot

run-worker:
	go run ./cmd/worker
