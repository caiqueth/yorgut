# Simple Makefile for a Go project

# DATABASE_DSN etc. from the local, gitignored .env (if present).
-include .env

GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.26.0
MIGRATIONS_DIR := internal/database/migrations

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Fill the database with test data
seed:
	@go run cmd/seed/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Migration and schema tests against a real Postgres (embedded-postgres
# downloads it on the first run; no Docker needed)
test_pg:
	@go test -tags pgtest ./internal/database/ -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

swagger:
	@swag init -g ./internal/server/routes.go --requiredByDefault

# Local Postgres (docker-compose.yml); see README "Database"
db_up:
	@docker compose up -d --wait

db_down:
	@docker compose down

# New SQL migration: make migration_new name=add_something
migration_new:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) -s create $(name) sql

# Applied and pending migrations of the Postgres at DATABASE_DSN
migration_status:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_DSN)" status

# Regenerate the Swagger spec, then the frontend API client (types + TanStack
# Query hooks in frontend/src/api/generated). Needs swag (go install
# github.com/swaggo/swag/cmd/swag@v1.16.4) and the frontend node_modules.
# --requiredByDefault marks every non-omitempty field required, so generated
# response types are not all optional.
gen-api: swagger
	@cd ../frontend && yarn gen:api

.PHONY: all build run seed test test_pg clean swagger gen-api db_up db_down migration_new migration_status
