# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY ../../shared ../shared

RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o marketing-service ./cmd/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary
COPY --from=builder /app/marketing-service .
COPY --from=builder /app/migrations ./migrations

# Expose ports
EXPOSE 8089 9089

CMD ["./marketing-service"]
