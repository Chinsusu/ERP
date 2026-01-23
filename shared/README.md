# Shared Libraries

This directory contains shared Go packages used across all microservices in the ERP system.

## Packages

### `config`
Configuration management using Viper. Loads environment variables and config files.

```go
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}
```

### `database`
PostgreSQL connection with GORM. Includes connection pooling and transaction helpers.

```go
db, err := database.Connect(database.NewDefaultConfig(cfg.GetDSN()))
if err != nil {
    log.Fatal(err)
}
```

### `logger`
Structured logging with Zap. Supports JSON and console output formats.

```go
logger.Init(&logger.Config{
    Level:  "info",
    Format: "json",
})

logger.Info("Server started", zap.Int("port", 8080))
```

### `jwt`
JWT token generation and verification for authentication.

```go
jwtManager := jwt.NewManager(secret, accessTokenExpire, refreshTokenExpire)
token, err := jwtManager.GenerateAccessToken(userID, email, roleIDs)
```

### `middleware`
Common HTTP middlewares for Gin framework:
- CORS
- Request ID
- Logger
- Recovery
- Auth (JWT validation)

```go
r := gin.New()
r.Use(middleware.RequestID())
r.Use(middleware.Logger(logger.Get()))
r.Use(middleware.Recovery(logger.Get()))
r.Use(middleware.CORS("*"))
```

### `errors`
Custom error types with HTTP status codes.

```go
return errors.NotFound("User")
return errors.Unauthorized("Invalid credentials")
return errors.Internal(err)
```

### `validator`
Struct validation using go-playground/validator with user-friendly error messages.

```go
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

if err := validator.Validate(&req); err != nil {
    // returns ValidationErrors
}
```

### `response`
Standard HTTP response helpers with consistent structure.

```go
// Success response
response.Success(c, data)

// Success with pagination
meta := response.NewMeta(page, pageSize, totalItems)
response.SuccessWithMeta(c, data, meta)

// Error response
response.Error(c, errors.NotFound("Resource"))
```

### `nats`
NATS JetStream client for event-driven architecture.

```go
client, err := nats.NewClient(&nats.Config{
    URL:    "nats://localhost:4222",
    Logger: logger.Get(),
})

// Publish event
client.Publish("user.created", userData)

// Subscribe to events
client.Subscribe("user.*", "user-queue", handler)
```

### `grpc`
gRPC server and client helpers with logging and recovery interceptors.

```go
// Server
server := grpc.NewServer("9090", logger.Get())
pb.RegisterUserServiceServer(server.GetServer(), &userService{})
server.Start()

// Client
client, err := grpc.NewClient("localhost:9090", logger.Get())
conn := client.GetConn()
```

## Usage

Import the packages in your service:

```go
import (
    "github.com/erp-cosmetics/shared/pkg/config"
    "github.com/erp-cosmetics/shared/pkg/database"
    "github.com/erp-cosmetics/shared/pkg/logger"
    "github.com/erp-cosmetics/shared/pkg/middleware"
)
```

## Proto Files

The `proto/` directory contains shared protobuf definitions:

- `common.proto`: Common message types used across services

### Generating Proto Code

```bash
# From the root directory
make proto

# Or manually
cd shared/proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       *.proto
```

## Dependencies

All dependencies are managed in `go.mod`. Key dependencies:

- **gin-gonic/gin**: HTTP framework
- **gorm.io/gorm**: ORM
- **uber.org/zap**: Logging
- **spf13/viper**: Configuration
- **golang-jwt/jwt**: JWT tokens
- **nats-io/nats.go**: NATS client
- **google.golang.org/grpc**: gRPC

## Development

When adding new shared packages:

1. Create the package directory under `pkg/`
2. Implement the package with tests
3. Update this README
4. Run `go mod tidy` in the shared directory

## Testing

```bash
cd shared
go test ./...
```
