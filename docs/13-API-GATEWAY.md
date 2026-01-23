# 13 - API GATEWAY

## TỔNG QUAN

API Gateway là điểm vào duy nhất (single entry point) cho tất cả client requests vào hệ thống ERP. Gateway chịu trách nhiệm routing, authentication, authorization, rate limiting, và request/response transformation.

### Responsibilities

✅ **Request Routing**: Forward requests đến đúng backend service  
✅ **Authentication**: Validate JWT access tokens  
✅ **Authorization**: Check permissions trước khi cho phép truy cập  
✅ **Rate Limiting**: Prevent abuse và protect backend services  
✅ **Request/Response Logging**: Centralized logging  
✅ **Circuit Breaking**: Handle service failures gracefully  
✅ **CORS**: Cross-Origin Resource Sharing  
✅ **API Versioning**: Support multiple API versions

### Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **Cache**: Redis (rate limiting, auth cache)
- **Load Balancing**: Round-robin to service replicas

### Port

- HTTP: `8080`

---

## ARCHITECTURE

```mermaid
graph TB
    Client[Client<br/>Web/Mobile] --> Nginx[Nginx<br/>SSL/TLS Termination]
    Nginx --> Gateway[API Gateway :8080]
    
    subgraph "Middleware Chain"
        Gateway --> M1[CORS]
        M1 --> M2[Request Logger]
        M2 --> M3[Recovery]
        M3 --> M4[Rate Limiter]
        M4 --> M5{Auth Required?}
        M5 -->|Yes| M6[JWT Validator]
        M5 -->|No| M7[Route to Service]
        M6 --> M8[Permission Checker]
        M8 --> M7
    end
    
    M7 --> AuthSvc[Auth Service :8081]
    M7 --> UserSvc[User Service :8082]
    M7 --> MasterSvc[Master Data :8083]
    M7 --> SupplierSvc[Supplier :8084]
    M7 --> ProcSvc[Procurement :8085]
    M7 --> WMSSvc[WMS :8086]
    M7 --> MfgSvc[Manufacturing :8087]
    M7 --> SalesSvc[Sales :8088]
    M7 --> MarketingSvc[Marketing :8089]
    M7 --> NotifSvc[Notification :8090]
    M7 --> FileSvc[File :8091]
    M7 --> ReportSvc[Reporting :8092]
    M7 --> AISvc[AI :8093]
    
    Gateway -.Cache.-> Redis[Redis]
    M6 -.gRPC.-> AuthSvc
    M8 -.gRPC.-> AuthSvc
```

---

## MIDDLEWARE CHAIN

### 1. CORS Middleware

Xử lý Cross-Origin requests.

```go
// Config
cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000",        // Dev frontend
        "https://erp.company.com",      // Production
    },
    AllowMethods: []string{
        "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
    },
    AllowHeaders: []string{
        "Authorization",
        "Content-Type",
        "X-Request-ID",
    },
    ExposeHeaders: []string{
        "X-Total-Count",
        "X-Page",
        "X-Per-Page",
    },
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
}
```

---

### 2. Request Logger Middleware

Log mỗi request để audit và debugging.

```go
// Log format
{
  "timestamp": "2024-01-23T15:00:00Z",
  "request_id": "req-uuid",
  "method": "POST",
  "path": "/api/v1/procurement/purchase-orders",
  "user_id": "user-uuid",
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "status": 201,
  "latency_ms": 145,
  "body_size": 1024
}
```

---

### 3. Recovery Middleware

Catch panics và trả về error response thay vì crash.

```go
// Response khi panic
{
  "error": "Internal Server Error",
  "request_id": "req-uuid",
  "timestamp": "2024-01-23T15:00:00Z"
}
```

---

### 4. Rate Limiter Middleware

Limit số requests per user/IP.

**Strategy**: Token bucket algorithm với Redis

```go
// Per-user rate limit
Key: "ratelimit:user:{user_id}"
Limit: 100 requests / minute
Burst: 20 requests

// Per-IP rate limit (cho unauthenticated requests)
Key: "ratelimit:ip:{ip_address}"
Limit: 30 requests / minute
Burst: 10 requests
```

**Response khi exceed limit (429 Too Many Requests)**:
```json
{
  "error": "Rate limit exceeded",
  "retry_after": 30,
  "limit": 100,
  "remaining": 0,
  "reset_at": "2024-01-23T15:01:00Z"
}
```

**Headers**:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1706018460
Retry-After: 30
```

---

### 5. JWT Validator Middleware

Validate access token từ Authorization header.

```go
// Extract token from header
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

// Steps:
1. Extract token from "Authorization: Bearer <token>"
2. Call Auth Service via gRPC: ValidateToken(token)
3. If valid, extract user_id and add to context
4. If invalid, return 401 Unauthorized
5. Check token in blacklist (Redis)
```

**Error Responses**:

```json
// Missing token
{
  "error": "Authorization header required",
  "code": "MISSING_TOKEN"
}

// Invalid token
{
  "error": "Invalid or expired token",
  "code": "INVALID_TOKEN"
}

// Blacklisted token (user logged out)
{
  "error": "Token has been revoked",
  "code": "REVOKED_TOKEN"
}
```

---

### 6. Permission Checker Middleware

Check permissions cho từng route.

```go
// Route configuration
{
  Path: "/api/v1/procurement/purchase-orders",
  Method: "POST",
  Permission: "procurement:po:create",
  Service: "procurement-service:8085"
}

// Check permission
1. Get user_id from context (set by JWT Validator)
2. Call Auth Service gRPC: CheckPermission(user_id, "procurement:po:create")
3. If allowed, proceed to service
4. If denied, return 403 Forbidden
```

**Error Response (403 Forbidden)**:
```json
{
  "error": "Permission denied",
  "code": "FORBIDDEN",
  "required_permission": "procurement:po:create"
}
```

---

## ROUTE CONFIGURATION

### Configuration File (YAML)

```yaml
api_gateway:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  
  cors:
    allowed_origins:
      - http://localhost:3000
      - https://erp.company.com
    allowed_methods:
      - GET
      - POST
      - PUT
      - DELETE
      - PATCH
    allowed_headers:
      - Authorization
      - Content-Type
    allow_credentials: true
    max_age: 43200 # 12 hours
  
  rate_limit:
    enabled: true
    requests_per_minute: 100
    burst: 20
    ip_limit: 30 # for unauthenticated
  
  circuit_breaker:
    enabled: true
    threshold: 5 # failures before opening circuit
    timeout: 10s # how long circuit stays open
  
routes:
  # ============================================
  # Auth Service (Public endpoints)
  # ============================================
  - path: /api/v1/auth/login
    method: POST
    service: auth-service:8081
    timeout: 5s
    auth_required: false
    rate_limit: 10/min # stricter for login
  
  - path: /api/v1/auth/refresh
    method: POST
    service: auth-service:8081
    timeout: 5s
    auth_required: false
  
  - path: /api/v1/auth/forgot-password
    method: POST
    service: auth-service:8081
    timeout: 5s
    auth_required: false
    rate_limit: 5/min
  
  - path: /api/v1/auth/reset-password
    method: POST
    service: auth-service:8081
    timeout: 5s
    auth_required: false
  
  # Auth - Protected endpoints
  - path: /api/v1/auth/*
    methods: [GET, POST, PUT, DELETE, PATCH]
    service: auth-service:8081
    timeout: 10s
    auth_required: true
    # Permissions checked internally by service
  
  # ============================================
  # User Service
  # ============================================
  - path: /api/v1/users
    method: GET
    service: user-service:8082
    timeout: 10s
    auth_required: true
    permission: user:user:read
  
  - path: /api/v1/users
    method: POST
    service: user-service:8082
    timeout: 10s
    auth_required: true
    permission: user:user:create
  
  - path: /api/v1/users/:id
    methods: [GET, PUT, DELETE]
    service: user-service:8082
    timeout: 10s
    auth_required: true
    permissions:
      GET: user:user:read
      PUT: user:user:update
      DELETE: user:user:delete
  
  - path: /api/v1/departments/*
    methods: [GET, POST, PUT, DELETE]
    service: user-service:8082
    timeout: 10s
    auth_required: true
    # Permissions vary by method
  
  # ============================================
  # Master Data Service
  # ============================================
  - path: /api/v1/master-data/materials
    methods: [GET, POST]
    service: master-data-service:8083
    timeout: 10s
    auth_required: true
    permissions:
      GET: master_data:material:read
      POST: master_data:material:create
  
  - path: /api/v1/master-data/products/*
    methods: [GET, POST, PUT, DELETE]
    service: master-data-service:8083
    timeout: 10s
    auth_required: true
  
  # ============================================
  # Supplier Service
  # ============================================
  - path: /api/v1/suppliers/*
    methods: [GET, POST, PUT, DELETE, PATCH]
    service: supplier-service:8084
    timeout: 10s
    auth_required: true
  
  # ============================================
  # Procurement Service
  # ============================================
  - path: /api/v1/procurement/requisitions
    methods: [GET, POST]
    service: procurement-service:8085
    timeout: 15s
    auth_required: true
    permissions:
      GET: procurement:pr:read
      POST: procurement:pr:create
  
  - path: /api/v1/procurement/purchase-orders
    methods: [GET, POST]
    service: procurement-service:8085
    timeout: 15s
    auth_required: true
    permissions:
      GET: procurement:po:read
      POST: procurement:po:create
  
  # ============================================
  # WMS Service
  # ============================================
  - path: /api/v1/warehouse/stock
    method: GET
    service: wms-service:8086
    timeout: 10s
    auth_required: true
    permission: wms:stock:read
  
  - path: /api/v1/warehouse/grn
    methods: [GET, POST]
    service: wms-service:8086
    timeout: 15s
    auth_required: true
    permissions:
      GET: wms:grn:read
      POST: wms:grn:create
  
  # ============================================
  # Manufacturing Service
  # ============================================
  - path: /api/v1/manufacturing/boms
    methods: [GET, POST]
    service: manufacturing-service:8087
    timeout: 10s
    auth_required: true
    permissions:
      GET: manufacturing:bom:read
      POST: manufacturing:bom:create
  
  - path: /api/v1/manufacturing/work-orders/*
    methods: [GET, POST, PUT, PATCH]
    service: manufacturing-service:8087
    timeout: 15s
    auth_required: true
  
  # ============================================
  # Sales Service
  # ============================================
  - path: /api/v1/sales/*
    methods: [GET, POST, PUT, DELETE, PATCH]
    service: sales-service:8088
    timeout: 10s
    auth_required: true
  
  # ============================================
  # Marketing Service
  # ============================================
  - path: /api/v1/marketing/*
    methods: [GET, POST, PUT, DELETE]
    service: marketing-service:8089
    timeout: 10s
    auth_required: true
  
  # ============================================
  # File Service
  # ============================================
  - path: /api/v1/files/upload
    method: POST
    service: file-service:8091
    timeout: 60s # longer for file uploads
    auth_required: true
    max_body_size: 10MB
  
  - path: /api/v1/files/*
    methods: [GET, DELETE]
    service: file-service:8091
    timeout: 30s
    auth_required: true
  
  # ============================================
  # Reporting Service
  # ============================================
  - path: /api/v1/reports/*
    method: GET
    service: reporting-service:8092
    timeout: 30s # reports may take longer
    auth_required: true
  
  # ============================================
  # AI Service
  # ============================================
  - path: /api/v1/ai/*
    methods: [GET, POST]
    service: ai-service:8093
    timeout: 60s # AI operations may be slow
    auth_required: true
```

---

## SERVICE DISCOVERY

API Gateway discovers backend services via **DNS** trong Docker network.

```go
// Service URLs
services := map[string]string{
    "auth-service": "http://auth-service:8081",
    "user-service": "http://user-service:8082",
    "master-data-service": "http://master-data-service:8083",
    // ...
}

// Health check every 30s
// If service is down, mark as unhealthy
// Circuit breaker opens after 5 consecutive failures
```

---

## CIRCUIT BREAKER PATTERN

Protect backend services khi có issues.

```
States:
- CLOSED: Normal traffic flow
- OPEN: Block all requests, return 503 Service Unavailable
- HALF-OPEN: Allow limited requests to test if service recovered

Transitions:
CLOSED --[5 failures]--> OPEN
OPEN --[after 10s]--> HALF-OPEN
HALF-OPEN --[success]--> CLOSED
HALF-OPEN --[failure]--> OPEN
```

**Response khi circuit is OPEN (503)**:
```json
{
  "error": "Service temporarily unavailable",
  "code": "SERVICE_UNAVAILABLE",
  "service": "procurement-service",
  "retry_after": 10
}
```

---

## REQUEST/RESPONSE TRANSFORMATION

### Request Enrichment

Gateway adds custom headers trước khi forward:

```
X-Request-ID: unique-request-uuid
X-User-ID: user-uuid (from JWT)
X-Forwarded-For: client-ip
X-Real-IP: client-ip
```

### Response Headers

Gateway adds metadata vào response:

```
X-Response-Time: 145ms
X-Service: procurement-service
X-Gateway-Version: 1.0.0
```

### Pagination Headers

Cho list endpoints:

```
X-Total-Count: 150
X-Page: 1
X-Per-Page: 20
X-Total-Pages: 8
Link: <http://...?page=2>; rel="next", <http://...?page=8>; rel="last"
```

---

## ERROR HANDLING

### Standard Error Response

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "request_id": "req-uuid",
  "timestamp": "2024-01-23T15:00:00Z",
  "details": {
    "field": "email",
    "reason": "Invalid email format"
  }
}
```

### HTTP Status Codes

| Code | Meaning | When |
|------|---------|------|
| 200 | OK | Success |
| 201 | Created | Resource created |
| 204 | No Content | Delete success |
| 400 | Bad Request | Invalid input |
| 401 | Unauthorized | Invalid/missing token |
| 403 | Forbidden | No permission |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Duplicate resource |
| 422 | Unprocessable Entity | Validation error |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 502 | Bad Gateway | Backend service error |
| 503 | Service Unavailable | Service down |
| 504 | Gateway Timeout | Backend timeout |

---

## CACHING STRATEGY

### Redis Cache Keys

```
# User permissions cache
Key: "cache:permissions:{user_id}"
TTL: 15 minutes

# User info cache
Key: "cache:user:{user_id}"
TTL: 15 minutes

# Rate limit counters
Key: "ratelimit:user:{user_id}"
TTL: 60 seconds

# Token blacklist
Key: "blacklist:token:{jti}"
TTL: Until token expiry
```

### Cache Invalidation

- On logout: Add token to blacklist
- On permission change: Clear user permissions cache
- On user update: Clear user info cache

---

## MONITORING & METRICS

### Prometheus Metrics

```
# Request metrics
api_gateway_requests_total{method, path, status}
api_gateway_request_duration_seconds{method, path}

# Error metrics
api_gateway_errors_total{type, service}

# Rate limiting
api_gateway_rate_limit_exceeded_total{user_id}

# Circuit breaker
api_gateway_circuit_breaker_state{service} # 0=closed, 1=open, 2=half-open

# Backend health
api_gateway_backend_health{service} # 0=down, 1=up
```

### Health Check Endpoint

```
GET /health

Response 200:
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 86400,
  "services": {
    "auth-service": "healthy",
    "user-service": "healthy",
    "procurement-service": "degraded",
    "wms-service": "healthy"
  }
}
```

---

## CONFIGURATION

### Environment Variables

```bash
# Server
API_GATEWAY_PORT=8080
API_GATEWAY_ENV=production

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=<secret>

# Auth Service Connection
AUTH_SERVICE_URL=http://auth-service:8081
AUTH_GRPC_URL=auth-service:9081

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_PER_MINUTE=100
RATE_LIMIT_BURST=20

# Circuit Breaker
CIRCUIT_BREAKER_ENABLED=true
CIRCUIT_BREAKER_THRESHOLD=5
CIRCUIT_BREAKER_TIMEOUT=10s

# Timeouts
DEFAULT_TIMEOUT=30s
UPLOAD_TIMEOUT=60s
REPORT_TIMEOUT=60s

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://erp.company.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## SECURITY CONSIDERATIONS

### TLS/SSL

- Nginx terminates SSL
- API Gateway communicates with Nginx over HTTP (internal network)
- Backend services communicate over HTTP (internal network)

### API Key (Future Enhancement)

Cho external integrations:

```
X-API-Key: api-key-value
```

### IP Whitelisting (Future)

Restrict access based on IP ranges.

### Request Size Limits

```
Max Header Size: 1MB
Max Body Size: 10MB (configurable per route)
Max Upload Size: 100MB (file uploads)
```

---

## DEPLOYMENT

### Docker Compose

```yaml
api-gateway:
  build: ./services/api-gateway
  ports:
    - "8080:8080"
  environment:
    - API_GATEWAY_PORT=8080
    - REDIS_HOST=redis
    - AUTH_SERVICE_URL=http://auth-service:8081
    - AUTH_GRPC_URL=auth-service:9081
    - RATE_LIMIT_PER_MINUTE=100
  depends_on:
    - redis
    - auth-service
  networks:
    - erp-network
  restart: unless-stopped
  deploy:
    replicas: 2 # For high availability
    resources:
      limits:
        cpus: '1'
        memory: 512M
```

---

## FUTURE ENHANCEMENTS

- [ ] GraphQL support
- [ ] WebSocket support cho real-time updates
- [ ] Request/Response caching
- [ ] API versioning (v2, v3)
- [ ] OpenAPI/Swagger documentation
- [ ] Advanced analytics & metrics
- [ ] A/B testing support
- [ ] Blue-green deployment support

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
