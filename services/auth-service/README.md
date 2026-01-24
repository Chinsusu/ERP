# Auth Service

Authentication and Authorization service for the ERP system.

## Features

- 🔐 JWT-based authentication (access & refresh tokens)
- 👥 Role-Based Access Control (RBAC)
- 🔑 Fine-grained permissions (`service:resource:action`)
- 🛡️ Account security (bcrypt, rate limiting, lockout)
- 🔄 Token rotation for enhanced security
- 📊 Redis caching for permissions
- 📡 Event publishing via NATS
- 🚀 gRPC API for internal services

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP), gRPC
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: NATS JetStream

## Ports

- **HTTP**: 8081
- **gRPC**: 9081

## Quick Start

### Development

```bash
# Install dependencies
make deps

# Run migrations
make migrate-up

# Run service
make run
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## API Endpoints

### REST API

- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh access token
- `GET /auth/me` - Get current user info
- `GET /auth/permissions` - Get user permissions
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password
- `GET /health` - Health check

### gRPC Methods

- `ValidateToken` - Verify JWT token
- `GetUserPermissions` - Get all user permissions
- `CheckPermission` - Check specific permission

## Database Schema

8 tables:
- `roles` - User roles
- `permissions` - Fine-grained permissions
- `role_permissions` - Role-permission mapping
- `user_credentials` - User authentication data
- `user_roles` - User-role assignments
- `refresh_tokens` - Refresh token storage
- `sessions` - Active sessions
- `password_reset_tokens` - Password reset tokens

## Default Credentials

**Admin User:**
- Email: `admin@company.vn`
- Password: `Admin@123`

**⚠️ Change this password immediately in production!**

## Environment Variables

Key variables (see `.env.example` for full list):

- `SERVICE_NAME` - Service identifier
- `PORT` - HTTP port (default: 8081)
- `GRPC_PORT` - gRPC port (default: 9081)
- `DB_HOST`, `DB_PORT`, `DB_NAME` - PostgreSQL connection
- `REDIS_HOST`, `REDIS_PORT` - Redis connection
- `NATS_URL` - NATS connection
- `JWT_SECRET` - JWT signing secret
- `JWT_ACCESS_TOKEN_EXPIRE` - Access token TTL (default: 15m)
- `JWT_REFRESH_TOKEN_EXPIRE` - Refresh token TTL (default: 7d)

## Security Features

### Password Security
- Bcrypt hashing with cost factor 12
- Minimum 8 characters required

### Account Protection
- 5 failed login attempts = 30-minute lockout
- Account can be manually disabled

### Token Security
- Access tokens: 15 minutes (short-lived)
- Refresh tokens: 7 days (stored in DB, can be revoked)
- Token rotation on refresh
- Token blacklisting via Redis

### Permission System

Format: `{service}:{resource}:{action}`

**Examples:**
- `user:user:read` - Read users
- `procurement:po:approve` - Approve purchase orders
- `wms:*:read` - Read all WMS resources
- `*:*:*` - Full access (super admin)

## Events Published

- `auth.user.logged_in` - User successfully logged in
- `auth.user.logged_out` - User logged out
- `auth.password.changed` - Password was changed
- `auth.account.locked` - Account locked due to failed attempts
- `auth.permission.granted` - Permission assigned to user
- `auth.permission.revoked` - Permission removed from user

## Development Guide

### Project Structure

```
auth-service/
├── cmd/main.go                 # Entry point
├── internal/
│   ├── domain/                 # Business entities & interfaces
│   ├── usecase/                # Business logic
│   ├── delivery/               # HTTP & gRPC handlers
│   └── infrastructure/         # DB, cache, events
├── migrations/                 # SQL migrations
├── proto/                      # gRPC definitions
├── Dockerfile                  # Production build
├── Dockerfile.dev              # Development build
└── Makefile                    # Build automation
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

### Database Migrations

```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### Generate Protobuf

```bash
make proto
```

## Production Deployment

1. Set strong `JWT_SECRET` in environment
2. Change default admin password
3. Configure Redis password
4. Enable SSL/TLS for gRPC
5. Set up monitoring and logging
6. Configure backup schedule

## License

Proprietary - All rights reserved
