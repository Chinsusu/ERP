# API Gateway

API Gateway for the ERP Cosmetics System - single entry point for all client requests.

## Features

- 🔀 **Routing** - Dynamic routing to 15+ backend services
- 🔐 **Authentication** - JWT validation with blacklist check
- ⏱️ **Rate Limiting** - Redis-based sliding window (100 req/min per user)
- 🔌 **Circuit Breaker** - Fault tolerance with automatic recovery
- 📝 **Request Logging** - Structured logging with request tracing
- 🌐 **CORS** - Cross-origin resource sharing support

## Port

- **HTTP**: 8080

## Quick Start

```bash
# Run locally
cd /opt/ERP/services/api-gateway
make run

# Health check
curl http://localhost:8080/health
```

## Middleware Chain

```
Request → RequestID → CORS → Logger → Recovery → RateLimit → Auth → CircuitBreaker → Proxy → Response
```

## Routes

| Prefix | Service | Auth Required |
|--------|---------|---------------|
| `/api/v1/auth` | auth-service:8081 | ❌ |
| `/api/v1/users` | user-service:8082 | ✅ |
| `/api/v1/departments` | user-service:8082 | ✅ |
| `/api/v1/categories` | master-data-service:8083 | ✅ |
| `/api/v1/units` | master-data-service:8083 | ✅ |
| `/api/v1/materials` | master-data-service:8083 | ✅ |
| `/api/v1/products` | master-data-service:8083 | ✅ |
| `/api/v1/suppliers` | supplier-service:8084 | ✅ |
| `/api/v1/procurement` | procurement-service:8085 | ✅ |
| `/api/v1/warehouse` | wms-service:8086 | ✅ |
| `/api/v1/manufacturing` | manufacturing-service:8087 | ✅ |
| `/api/v1/sales` | sales-service:8088 | ✅ |
| `/api/v1/marketing` | marketing-service:8089 | ✅ |
| `/api/v1/files` | file-service:8091 | ✅ |
| `/api/v1/reports` | reporting-service:8092 | ✅ |

## Health Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Aggregate health of all services |
| `GET /ready` | Readiness probe |
| `GET /live` | Liveness probe |
| `GET /health/:service` | Specific service health |

## Rate Limiting

- **Authenticated users**: 100 requests/minute
- **Unauthenticated (by IP)**: 30 requests/minute
- **Sliding window algorithm** using Redis

### Response Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1706018460
```

### Rate Limit Exceeded (429)

```json
{
  "error": "Rate limit exceeded",
  "retry_after": 60,
  "limit": 100,
  "remaining": 0
}
```

## Circuit Breaker

- **Threshold**: 5 consecutive failures → OPEN
- **Timeout**: 30 seconds in OPEN state
- **Half-open**: Allow 1 test request

### States

```
CLOSED ──[5 failures]──► OPEN ──[30s]──► HALF-OPEN ──[success]──► CLOSED
                            ▲                   │
                            └──[failure]────────┘
```

## Authentication

JWT token in Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Error Responses

| Code | Error |
|------|-------|
| 401 | Missing/invalid/expired token |
| 403 | Permission denied |

## Request Headers Added

The gateway adds these headers to downstream requests:

| Header | Value |
|--------|-------|
| `X-Request-ID` | Unique request UUID |
| `X-User-ID` | User ID from JWT |
| `X-Forwarded-For` | Client IP |
| `X-Real-IP` | Client IP |

## Response Headers Added

| Header | Value |
|--------|-------|
| `X-Response-Time` | Request latency |
| `X-Service` | Backend service name |
| `X-Gateway-Version` | Gateway version |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP port |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `JWT_SECRET` | - | JWT signing key |
| `RATE_LIMIT_ENABLED` | true | Enable rate limiting |
| `RATE_LIMIT_PER_MIN` | 100 | Requests per minute |
| `CIRCUIT_BREAKER_ENABLED` | true | Enable circuit breaker |
| `CIRCUIT_BREAKER_THRESHOLD` | 5 | Failures before open |

## Project Structure

```
api-gateway/
├── cmd/main.go
├── config/routes.yaml
├── internal/
│   ├── config/config.go
│   ├── middleware/
│   │   ├── request_id.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── recovery.go
│   │   ├── rate_limiter.go
│   │   ├── auth.go
│   │   └── circuit_breaker.go
│   ├── proxy/
│   │   ├── handler.go
│   │   └── service_registry.go
│   ├── health/handler.go
│   └── router/router.go
├── go.mod
├── Makefile
└── Dockerfile
```

---

**Port**: 8080  
**Status**: Ready for testing
