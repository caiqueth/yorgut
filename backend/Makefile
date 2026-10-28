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

# Migration, schema and service tests against a real MySQL: the local one
# (make db_up first), or the server TEST_MYSQL_DSN names (a user that may
# create and drop databases; never production)
test_mysql:
	@go test -count=1 -tags mysqltest -p 1 ./internal/database/ ./internal/service/ -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

swagger:
	@swag init -g ./internal/server/routes.go --requiredByDefault

# The repo-root docker-compose.yml
COMPOSE := docker compose -f ../docker-compose.yml

# Local MySQL; see README "Database"
db_up:
	@$(COMPOSE) up -d --wait db

db_down:
	@$(COMPOSE) stop db

# Local SeaweedFS, an S3 bucket for uploads; see README "Uploads"
storage_up:
	@$(COMPOSE) up -d --wait storage

storage_down:
	@$(COMPOSE) stop storage

# Copy the pictures in UPLOADS_DIR into the S3 bucket (S3_* in .env), keeping
# their keys; run it before switching STORAGE from disk to s3. Safe to rerun.
uploads_to_s3:
	@go run ./cmd/uploads-to-s3

# New SQL migration: make migration_new name=add_something
migration_new:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) -s create $(name) sql

# Applied and pending migrations of the MySQL at DATABASE_DSN (in the Go
# driver's form, user:pass@tcp(host:3306)/db, not a mysql:// URL)
migration_status:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DATABASE_DSN)" status

# Regenerate the Swagger spec, then the frontend API client (types + TanStack
# Query hooks in frontend/src/api/generated). Needs swag (go install
# github.com/swaggo/swag/cmd/swag@v1.16.4) and the frontend node_modules.
# --requiredByDefault marks every non-omitempty field required, so generated
# response types are not all optional.
gen-api: swagger
	@cd ../frontend && yarn gen:api

.PHONY: all build run seed test test_mysql clean swagger gen-api db_up db_down storage_up storage_down uploads_to_s3 migration_new migration_status
