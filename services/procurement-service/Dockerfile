# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy shared package first
COPY shared/ /app/shared/

# Copy go mod files
COPY services/procurement-service/go.mod services/procurement-service/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY services/procurement-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Expose ports
EXPOSE 8085 9085

# Run the application
CMD ["./main"]
