# Simple Makefile for a Go project

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

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

swagger:
	@swag init -g ./internal/server/routes.go

run_local_db:
	@docker run --name yorgut-dev -e POSTGRES_USER=localdev -e POSTGRES_PASSWORD=localdev -e POSTGRES_DB=yorgut-dev -p 5432:5432 -d postgres

migration_up:
	@migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose up

migration_down:
	@migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose down

migration_fix:
	@migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" force VERSION

.PHONY: all build run test clean migration_up migration_down migration_fix watch
