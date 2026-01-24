# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service
COPY services/manufacturing-service/ ./services/manufacturing-service/

WORKDIR /app/services/manufacturing-service

# Download dependencies
RUN go mod download

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /manufacturing-service ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /manufacturing-service .

# Copy migrations
COPY services/manufacturing-service/migrations ./migrations

EXPOSE 8087 9087

CMD ["./manufacturing-service"]
