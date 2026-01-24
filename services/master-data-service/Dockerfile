FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
COPY ../../shared ./../../shared

RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /master-data-service cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /master-data-service .
COPY migrations ./migrations

EXPOSE 8083 9083

CMD ["./master-data-service"]
