.PHONY: test test-coverage test-unit test-integration run build docker-up docker-down clean

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run only unit tests
test-unit:
	go test -v -short ./internal/domain/entity/... ./internal/usecase/...

# Run tests for a specific package
test-package:
	go test -v ./$(PKG)

# Run application locally
run:
	go run cmd/api/main.go

# Run application with live reload (air)
dev:
	~/go/bin/air -c .air.toml

# Build application
build:
	go build -o bin/orders-api cmd/api/main.go

# Start Docker containers
docker-up:
	docker-compose up --build

# Stop Docker containers
docker-down:
	docker-compose down

# Clean generated files
clean:
	rm -rf tmp/ bin/ coverage.out coverage.html docs/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate Swagger documentation
swagger:
	~/go/bin/swag init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated in docs/"

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...

# Show test coverage by package
coverage-by-package:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Run tests in watch mode (requires entr)
test-watch:
	find . -name "*.go" | entr -c go test ./...
