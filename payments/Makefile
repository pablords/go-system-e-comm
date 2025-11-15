.PHONY: proto run test clean migrate build docker-build docker-run

# Generate proto files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/payment.proto

# Run the application
run:
	go run cmd/grpc/main.go

# Build the application
build:
	go build -o bin/payments-service cmd/grpc/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean generated files and binaries
clean:
	rm -rf bin/ coverage.out coverage.html
	find . -name "*.pb.go" -delete

# Run migrations
migrate:
	@echo "Running migrations..."
	@mysql -h ${DB_HOST} -P ${DB_PORT} -u ${DB_USER} -p${DB_PASSWORD} ${DB_NAME} < migrations/001_create_tables.sql

# Docker build
docker-build:
	docker build -t payments-service .

# Docker run
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run
