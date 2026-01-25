# Multi-stage build for production
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy shared module first
COPY shared/ ./shared/

# Copy service code
COPY services/auth-service/ ./services/auth-service/

WORKDIR /build/services/auth-service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=builder /build/services/auth-service/main .

# Copy migrations if exists
COPY --from=builder /build/services/auth-service/migrations ./migrations/ 2>/dev/null || true

# Expose ports
EXPOSE 8081 9081

CMD ["./main"]
