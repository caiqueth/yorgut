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

.PHONY: all build run seed test clean swagger run_local_db

# Regenerate the Swagger spec, then the frontend API client (types + TanStack
# Query hooks in frontend/src/api/generated). Needs swag (go install
# github.com/swaggo/swag/cmd/swag@v1.16.4) and the frontend node_modules.
# --requiredByDefault marks every non-omitempty field required, so generated
# response types are not all optional.
gen-api:
	@swag init -g ./internal/server/routes.go --requiredByDefault
	@cd ../frontend && yarn gen:api

.PHONY: gen-api
