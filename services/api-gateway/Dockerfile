FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy shared module first
COPY shared/ ./shared/

# Copy service code
COPY services/api-gateway/ ./services/api-gateway/

WORKDIR /build/services/api-gateway

# Download dependencies
ENV GOTOOLCHAIN=auto
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

# Copy the binary
COPY --from=builder /build/services/api-gateway/main .

# Copy config
COPY --from=builder /build/services/api-gateway/config ./config

EXPOSE 8080

CMD ["./main"]
