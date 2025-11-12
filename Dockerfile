# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag CLI for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /build/app .

# Copy .env file (in production, use docker secrets instead)
COPY .env.example .env

EXPOSE 4000

CMD ["./app"]
