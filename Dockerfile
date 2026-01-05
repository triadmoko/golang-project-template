# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git and ca-certificates (required for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go.work and all go.mod/go.sum files first for better caching
COPY go.work go.work.sum ./
COPY libs/config/go.mod libs/config/
COPY libs/crypto/go.mod libs/crypto/go.sum libs/crypto/
COPY libs/database/go.mod libs/database/go.sum libs/database/
COPY libs/domain/go.mod libs/domain/go.sum libs/domain/
COPY libs/errors/go.mod libs/errors/
COPY libs/httputil/go.mod libs/httputil/go.sum libs/httputil/
COPY libs/jwt/go.mod libs/jwt/go.sum libs/jwt/
COPY libs/logger/go.mod libs/logger/go.sum libs/logger/
COPY libs/pagination/go.mod libs/pagination/go.sum libs/pagination/
COPY services/api-gateway/go.mod services/api-gateway/go.sum services/api-gateway/

# Download dependencies
RUN go mod download

# Copy source code
COPY libs/ libs/
COPY services/api-gateway/ services/api-gateway/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api-gateway \
    ./services/api-gateway/cmd

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezone
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/bin/api-gateway .

# Copy migrations if needed
COPY --from=builder /app/services/api-gateway/migrations ./migrations

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./api-gateway"]
