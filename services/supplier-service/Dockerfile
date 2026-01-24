# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy shared package
COPY shared/ /app/shared/

# Copy service files
COPY services/supplier-service/ /app/services/supplier-service/

WORKDIR /app/services/supplier-service

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /supplier-service cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /supplier-service /app/supplier-service

# Expose ports
EXPOSE 8084 9084

# Run the service
CMD ["/app/supplier-service"]
