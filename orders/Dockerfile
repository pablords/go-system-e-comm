FROM golang:1.25-alpine AS base

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create tmp directory with proper permissions
RUN mkdir -p tmp && chmod -R 777 tmp

# Expose port
EXPOSE 8080

# Use air for live reload
CMD ["air", "-c", ".air.toml"]
