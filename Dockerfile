# Stage 1: Build binary
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
# COPY file .env vào trong image
COPY .env .env
# Build the Go binary
RUN go build -o main .

# Stage 2: Minimal image
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy file .env từ builder stage
COPY --from=builder /app/.env .env

# Expose the application port (chỉnh theo app bạn)
EXPOSE 8080

# Run the binary
CMD ["./main"]
