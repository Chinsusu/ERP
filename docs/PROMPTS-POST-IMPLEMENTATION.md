# 🧪 PROMPTS FOR POST-IMPLEMENTATION PHASES

> Prompts cho Testing, Documentation, Deployment, Maintenance
> Sử dụng sau khi hoàn thành Implementation (Phase 1-6)

---

## 📋 MỤC LỤC

1. [Phase 7: Testing & QA](#phase-7-testing--qa)
2. [Phase 9: Documentation & Training](#phase-9-documentation--training)
3. [Phase 10: Deployment & Go-Live](#phase-10-deployment--go-live)
4. [Phase 11: Monitoring & Maintenance](#phase-11-monitoring--maintenance)

---

# PHASE 7: TESTING & QA

## PROMPT 7.1: Unit Tests - Auth & User Services

```markdown
# CONTEXT
Hệ thống ERP mỹ phẩm đã implement xong. Cần viết Unit Tests.

## Tech Stack
- Go 1.22+
- Testing: testify/assert, testify/mock, testify/suite
- Mocking: mockery hoặc manual mocks
- Coverage tool: go test -cover

## YÊU CẦU
Viết Unit Tests cho Auth Service và User Service với coverage > 80%.

### 1. Auth Service Tests

#### 1.1 LoginUseCase Tests
```go
// Test cases cần cover:
// - Login success với valid credentials
// - Login fail với wrong password
// - Login fail với non-existent email
// - Login fail khi account locked
// - Login fail khi account inactive
// - Account lock sau 5 failed attempts
// - Rate limiting check
```

#### 1.2 Token Tests
```go
// Test cases:
// - Generate access token success
// - Generate refresh token success
// - Validate token success
// - Validate expired token → error
// - Validate invalid signature → error
// - Refresh token success
// - Refresh with revoked token → error
```

#### 1.3 Permission Tests
```go
// Test cases:
// - Check permission với exact match
// - Check permission với wildcard (wms:*:read)
// - Check permission với super wildcard (*:*:*)
// - User không có permission → denied
// - Role hierarchy check
```

### 2. User Service Tests

#### 2.1 UserUseCase Tests
```go
// Test cases:
// - Create user success
// - Create user với duplicate email → error
// - Update user success
// - Delete user (soft delete)
// - Get user by ID
// - List users với pagination
// - Assign role to user
// - Remove role from user
```

#### 2.2 DepartmentUseCase Tests
```go
// Test cases:
// - Create department success
// - Create nested department (với parent)
// - Update department
// - Delete department (check có users không)
// - Get department tree
// - Move department (change parent)
```

### 3. Mock Interfaces

```go
// Cần mock:
type MockUserRepository interface {
    Create(user *User) error
    FindByID(id string) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id string) error
    List(filter UserFilter) ([]User, int64, error)
}

type MockAuthClient interface {
    ValidateToken(token string) (*UserInfo, error)
    CheckPermission(userID, permission string) bool
}
```

### 4. Test Structure

```
services/auth-service/
├── internal/
│   ├── usecase/
│   │   ├── login_usecase.go
│   │   ├── login_usecase_test.go      ← tạo file này
│   │   ├── token_usecase.go
│   │   ├── token_usecase_test.go      ← tạo file này
│   │   └── ...
│   └── domain/
│       └── repository/
│           └── mocks/                  ← tạo folder này
│               ├── user_repository_mock.go
│               └── token_repository_mock.go
```

## OUTPUT
- Tất cả test files với đầy đủ test cases
- Mock files
- Makefile target: `make test-auth`, `make test-user`
- Coverage report command
```

---

## PROMPT 7.2: Unit Tests - WMS Service (CRITICAL)

```markdown
# CONTEXT
WMS Service là critical service với FEFO logic. Cần test kỹ.

## YÊU CẦU
Viết Unit Tests cho WMS Service, đặc biệt FEFO logic.

### 1. FEFO Logic Tests (CRITICAL)

```go
// Test cases cho IssueStockFEFO:

func TestIssueStockFEFO_SingleLot_Success(t *testing.T) {
    // Given: 1 lot với 100 units, expiry 2025-12-31
    // When: Issue 50 units
    // Then: Issue từ lot đó, remaining = 50
}

func TestIssueStockFEFO_MultipleLots_IssueFromEarliestExpiry(t *testing.T) {
    // Given:
    //   - Lot A: 50 units, expiry 2025-06-30 (sớm hơn)
    //   - Lot B: 100 units, expiry 2025-12-31
    // When: Issue 70 units
    // Then: 
    //   - Issue 50 từ Lot A (hết)
    //   - Issue 20 từ Lot B
    //   - Lot B remaining = 80
}

func TestIssueStockFEFO_SkipExpiredLots(t *testing.T) {
    // Given:
    //   - Lot A: 50 units, expiry 2024-01-01 (ĐÃ HẾT HẠN)
    //   - Lot B: 100 units, expiry 2025-12-31
    // When: Issue 30 units
    // Then: Issue từ Lot B (skip Lot A)
}

func TestIssueStockFEFO_InsufficientStock_Error(t *testing.T) {
    // Given: Total available = 100 units
    // When: Issue 150 units
    // Then: Return ErrInsufficientStock
}

func TestIssueStockFEFO_ReservedStockNotIssued(t *testing.T) {
    // Given:
    //   - Lot A: quantity=100, reserved=30, available=70
    // When: Issue 80 units
    // Then: Error (chỉ có 70 available)
}

func TestIssueStockFEFO_MultipleLocations(t *testing.T) {
    // Given: Same lot ở multiple locations
    // When: Issue
    // Then: Issue từ tất cả locations của lot đó
}
```

### 2. Stock Reservation Tests

```go
// Test cases:
// - Reserve success
// - Reserve khi không đủ stock → error
// - Reserve đã tồn tại cho cùng reference → error hoặc update
// - Release reservation success
// - Release non-existent reservation → error
// - Auto-expire reservation sau timeout
```

### 3. GRN Tests

```go
// Test cases:
// - Create GRN từ PO success
// - Create GRN với quantity > PO ordered → warning/error
// - Complete GRN → stock tăng
// - Complete GRN → lot created
// - Complete GRN → event published
// - GRN với QC failed → stock vào quarantine
```

### 4. Stock Movement Tests

```go
// Test cases:
// - Movement IN → stock tăng
// - Movement OUT → stock giảm
// - Movement TRANSFER → from giảm, to tăng
// - Movement ADJUSTMENT → stock update
// - Movement audit trail created
```

### 5. Lot Expiry Tests

```go
// Test cases:
// - Lot expiring trong 90 ngày → alert
// - Lot expiring trong 30 ngày → daily alert
// - Lot expired → auto block
// - Get expiring lots query
```

### 6. Test Data Builder

```go
// Helper để tạo test data
type LotBuilder struct {
    lot *Lot
}

func NewLotBuilder() *LotBuilder {
    return &LotBuilder{
        lot: &Lot{
            ID:         uuid.New().String(),
            LotNumber:  "LOT-TEST-001",
            Status:     "AVAILABLE",
            QCStatus:   "PASSED",
        },
    }
}

func (b *LotBuilder) WithMaterial(id string) *LotBuilder {
    b.lot.MaterialID = id
    return b
}

func (b *LotBuilder) WithQuantity(qty float64) *LotBuilder {
    b.lot.Quantity = qty
    return b
}

func (b *LotBuilder) WithExpiry(date time.Time) *LotBuilder {
    b.lot.ExpiryDate = date
    return b
}

func (b *LotBuilder) Build() *Lot {
    return b.lot
}
```

## OUTPUT
- Đầy đủ test files cho WMS Service
- Test data builders
- Coverage > 90% cho FEFO logic
```

---

## PROMPT 7.3: Unit Tests - Manufacturing Service

```markdown
# CONTEXT
Manufacturing Service có BOM encryption và traceability. Cần test kỹ.

## YÊU CẦU

### 1. BOM Encryption Tests

```go
// Test cases:
// - Encrypt formula details success
// - Decrypt formula details success
// - Decrypt với wrong key → error
// - Encrypted data không readable
// - Re-encrypt với new key
```

### 2. BOM Access Control Tests

```go
// Test cases:
// - User với permission "bom:formula_view" → xem full BOM
// - User với permission "bom:quantity_view" → xem quantities, không formula
// - User với permission "bom:read" only → xem materials list only
// - User không có permission → denied
// - Audit log created khi access BOM
```

### 3. Work Order Tests

```go
// Test cases:
// - Create WO từ BOM success
// - WO lifecycle: PLANNED → RELEASED → IN_PROGRESS → COMPLETED
// - Start WO → materials reserved (call WMS)
// - Issue material to WO → stock giảm
// - Complete WO → finished goods created
// - Complete WO → yield calculated
// - Cancel WO → reservations released
```

### 4. Traceability Tests

```go
// Test cases:
func TestBackwardTrace_Success(t *testing.T) {
    // Given: Product lot BATCH-001 produced from WO-001
    //        WO-001 used: Lot-A (Material 1), Lot-B (Material 2)
    // When: Trace backward from BATCH-001
    // Then: Return Lot-A, Lot-B với quantities
}

func TestForwardTrace_Success(t *testing.T) {
    // Given: Material Lot-A used in WO-001, WO-002
    //        WO-001 → BATCH-001
    //        WO-002 → BATCH-002
    // When: Trace forward from Lot-A
    // Then: Return BATCH-001, BATCH-002
}

func TestTraceability_CompleteChain(t *testing.T) {
    // Given: Supplier lot → Internal lot → Product lot
    // When: Full trace
    // Then: Complete chain từ supplier đến finished goods
}
```

### 5. QC Tests

```go
// Test cases:
// - QC inspection pass → WO continues
// - QC inspection fail → NCR created
// - QC inspection fail → notification sent
// - QC approval workflow
```

## OUTPUT
- Đầy đủ test files cho Manufacturing Service
- Mock cho WMS gRPC client
- Mock cho encryption service
```

---

## PROMPT 7.4: Integration Tests (API Tests)

```markdown
# CONTEXT
Cần viết Integration Tests để test các API endpoints end-to-end.

## Tech Stack
- Postman/Newman hoặc
- Go httptest + testcontainers
- K6 cho load testing

## YÊU CẦU

### 1. Test Environment Setup

```yaml
# docker-compose.test.yml
version: '3.9'
services:
  postgres-test:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: test
    tmpfs:
      - /var/lib/postgresql/data
  
  redis-test:
    image: redis:7-alpine
  
  nats-test:
    image: nats:latest
    command: "-js"
```

### 2. Auth API Tests

```
POST /api/v1/auth/login
  ✓ 200: Valid credentials
  ✓ 401: Wrong password
  ✓ 401: User not found
  ✓ 423: Account locked
  ✓ 429: Rate limited

POST /api/v1/auth/refresh
  ✓ 200: Valid refresh token
  ✓ 401: Expired refresh token
  ✓ 401: Revoked refresh token

GET /api/v1/auth/me
  ✓ 200: With valid token
  ✓ 401: Without token
  ✓ 401: With expired token
```

### 3. WMS API Tests

```
GET /api/v1/stock
  ✓ 200: List stock với pagination
  ✓ 200: Filter by warehouse
  ✓ 200: Filter by material
  ✓ 401: Unauthorized

POST /api/v1/grn
  ✓ 201: Create GRN success
  ✓ 400: Invalid PO reference
  ✓ 400: Quantity exceeds PO
  ✓ 403: No permission

POST /api/v1/goods-issue
  ✓ 200: Issue with FEFO
  ✓ 400: Insufficient stock
  ✓ 200: Partial issue từ multiple lots

GET /api/v1/stock/expiring?days=30
  ✓ 200: Return expiring lots
```

### 4. Manufacturing API Tests

```
GET /api/v1/boms/:id
  ✓ 200: Full BOM với formula (có permission)
  ✓ 200: Partial BOM (không có formula permission)
  ✓ 403: Không có permission

POST /api/v1/work-orders
  ✓ 201: Create WO success
  ✓ 400: Invalid BOM
  ✓ 400: Product không có active BOM

PATCH /api/v1/work-orders/:id/start
  ✓ 200: Start success, materials reserved
  ✓ 400: Insufficient materials
  ✓ 400: WO not in RELEASED status
```

### 5. Cross-Service Integration Tests

```
# Test: PO → GRN → Stock → WO → Finished Goods

1. Create PO (Procurement)
2. Confirm PO → Event to WMS
3. Create GRN (WMS)
4. Complete GRN → Stock updated
5. Create WO (Manufacturing)
6. Start WO → Materials reserved
7. Issue materials → Stock reduced (FEFO)
8. Complete WO → Finished goods created
9. Verify traceability
```

### 6. Postman Collection Structure

```
ERP-Cosmetics-API-Tests/
├── Auth/
│   ├── Login.json
│   ├── Refresh.json
│   └── Permissions.json
├── Users/
├── MasterData/
├── Suppliers/
├── Procurement/
├── WMS/
│   ├── Stock.json
│   ├── GRN.json
│   ├── GoodsIssue.json
│   └── FEFO-Tests.json
├── Manufacturing/
│   ├── BOM.json
│   ├── WorkOrder.json
│   └── Traceability.json
├── Sales/
├── Marketing/
└── E2E-Flows/
    ├── Procurement-to-Stock.json
    ├── Order-to-Delivery.json
    └── Material-to-Product.json
```

## OUTPUT
- Postman collection JSON files
- Newman CI script
- docker-compose.test.yml
- Test data seeding scripts
```

---

## PROMPT 7.5: Performance & Load Testing

```markdown
# CONTEXT
Cần test performance để đảm bảo hệ thống handle được 100 concurrent users.

## Tech Stack
- K6 (load testing)
- Grafana (visualize results)

## YÊU CẦU

### 1. Performance Requirements

| Metric | Target |
|--------|--------|
| API Response Time (p95) | < 200ms |
| API Response Time (p99) | < 500ms |
| Concurrent Users | 100 |
| Requests/Second | 500+ |
| Error Rate | < 0.1% |
| Database Query Time (p95) | < 50ms |

### 2. K6 Test Scripts

```javascript
// tests/load/stock-inquiry.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 20 },   // Ramp up
    { duration: '3m', target: 50 },   // Stay at 50
    { duration: '2m', target: 100 },  // Peak load
    { duration: '1m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<200', 'p(99)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  // Login và get token
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, {
    email: 'test@company.vn',
    password: 'Test@123',
  });
  return { token: loginRes.json('access_token') };
}

export default function(data) {
  const headers = {
    'Authorization': `Bearer ${data.token}`,
    'Content-Type': 'application/json',
  };

  // Test stock inquiry
  const stockRes = http.get(`${BASE_URL}/api/v1/stock?page=1&limit=20`, { headers });
  
  check(stockRes, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(1);
}
```

### 3. Critical Endpoints to Test

```javascript
// Endpoints cần load test:

// 1. Stock Inquiry (high frequency)
GET /api/v1/stock

// 2. FEFO Goods Issue (complex logic)
POST /api/v1/goods-issue

// 3. BOM Retrieval (với decryption)
GET /api/v1/boms/:id

// 4. Work Order Operations
POST /api/v1/work-orders/:id/start
POST /api/v1/work-orders/:id/issue-material

// 5. Report Queries (heavy DB)
GET /api/v1/reports/stock-summary
GET /api/v1/reports/expiring-lots
```

### 4. Database Performance Tests

```sql
-- Queries cần optimize và test:

-- 1. Stock by material (thường xuyên)
EXPLAIN ANALYZE
SELECT * FROM stock 
WHERE material_id = $1 
AND available_quantity > 0;

-- 2. FEFO query (critical)
EXPLAIN ANALYZE
SELECT l.*, s.available_quantity
FROM lots l
JOIN stock s ON l.id = s.lot_id
WHERE l.material_id = $1
AND l.status = 'AVAILABLE'
AND l.expiry_date > NOW()
ORDER BY l.expiry_date ASC;

-- 3. Traceability query
EXPLAIN ANALYZE
SELECT * FROM batch_traceability
WHERE finished_lot_id = $1;
```

### 5. Stress Test Scenarios

```javascript
// Scenario 1: End of month closing
// - 50 users doing stock inquiry
// - 20 users creating GRNs
// - 10 users issuing goods
// - 5 users running reports

// Scenario 2: Production peak
// - 30 work orders running simultaneously
// - Material issues every 5 seconds
// - QC inspections

// Scenario 3: Sales rush
// - 100 sales orders in 1 hour
// - Stock reservations
// - Delivery processing
```

## OUTPUT
- K6 test scripts cho tất cả critical endpoints
- Grafana dashboard JSON
- Performance test report template
- Optimization recommendations
```

---

## PROMPT 7.6: Security Testing

```markdown
# CONTEXT
Cần security audit cho hệ thống ERP trước khi go-live.

## YÊU CẦU

### 1. OWASP Top 10 Checklist

```markdown
## A01: Broken Access Control
[ ] Test horizontal privilege escalation (user A access user B data)
[ ] Test vertical privilege escalation (staff access admin functions)
[ ] Test IDOR (Insecure Direct Object References)
[ ] Test missing function level access control
[ ] Test CORS misconfiguration

## A02: Cryptographic Failures
[ ] Check TLS configuration (TLS 1.2+)
[ ] Check password hashing (bcrypt cost >= 12)
[ ] Check JWT signing algorithm (RS256, not HS256)
[ ] Check BOM encryption (AES-256-GCM)
[ ] Check sensitive data in logs

## A03: Injection
[ ] SQL Injection tests
[ ] NoSQL Injection tests
[ ] Command Injection tests
[ ] LDAP Injection tests

## A04: Insecure Design
[ ] Check business logic flaws
[ ] Check race conditions (stock reservation)
[ ] Check negative quantity handling

## A05: Security Misconfiguration
[ ] Check default credentials removed
[ ] Check unnecessary features disabled
[ ] Check error messages (no stack traces)
[ ] Check security headers

## A06: Vulnerable Components
[ ] Check Go dependencies (govulncheck)
[ ] Check npm dependencies (npm audit)
[ ] Check Docker base images

## A07: Authentication Failures
[ ] Brute force protection
[ ] Session management
[ ] Password policy enforcement
[ ] MFA implementation (nếu có)

## A08: Software and Data Integrity
[ ] Check CI/CD pipeline security
[ ] Check dependency integrity

## A09: Logging and Monitoring
[ ] Check sensitive data không log
[ ] Check audit logging đầy đủ
[ ] Check log injection prevention

## A10: SSRF
[ ] Check server-side request forgery
```

### 2. Security Test Cases

```go
// Authentication Tests
func TestBruteForceProtection(t *testing.T) {
    // Attempt 10 failed logins
    // Verify account locked after 5
    // Verify lockout duration
}

func TestJWTManipulation(t *testing.T) {
    // Test with modified payload
    // Test with different algorithm
    // Test with expired token
    // Test with token from different env
}

// Authorization Tests
func TestHorizontalPrivilegeEscalation(t *testing.T) {
    // User A tries to access User B's data
    // User A tries to modify User B's records
}

func TestVerticalPrivilegeEscalation(t *testing.T) {
    // Staff tries to access admin endpoints
    // Staff tries to assign admin role
}

// Injection Tests
func TestSQLInjection(t *testing.T) {
    // Test with: ' OR '1'='1
    // Test with: '; DROP TABLE users;--
    // Test in search, filter, sort params
}

// Business Logic Tests
func TestRaceCondition_StockReservation(t *testing.T) {
    // 10 concurrent requests to reserve same stock
    // Verify no over-reservation
}

func TestNegativeQuantity(t *testing.T) {
    // Try to issue negative quantity
    // Try to receive negative quantity
}
```

### 3. Security Headers Check

```go
// Required headers:
// X-Content-Type-Options: nosniff
// X-Frame-Options: DENY
// X-XSS-Protection: 1; mode=block
// Strict-Transport-Security: max-age=31536000; includeSubDomains
// Content-Security-Policy: default-src 'self'
// Referrer-Policy: strict-origin-when-cross-origin
```

### 4. Penetration Test Scenarios

```markdown
## Scenario 1: External Attacker
- Port scanning
- Service enumeration
- Authentication bypass attempts
- API fuzzing

## Scenario 2: Malicious User (Staff)
- Access other department data
- Escalate privileges
- Data exfiltration
- Audit log tampering

## Scenario 3: Compromised Admin
- Access BOM formulas
- Export customer data
- Create backdoor account
- Disable security features
```

### 5. BOM Security Specific Tests

```go
// BOM là trade secret, cần test kỹ:

func TestBOMFormulaNotInLogs(t *testing.T) {
    // Access BOM
    // Check logs không chứa formula content
}

func TestBOMFormulaNotInResponse_NoPermission(t *testing.T) {
    // User không có formula_view permission
    // Response không chứa formula_details
}

func TestBOMEncryptionAtRest(t *testing.T) {
    // Query database directly
    // Verify formula_details is encrypted bytes
}

func TestBOMAuditLog(t *testing.T) {
    // Access BOM
    // Verify audit log created với user info
}
```

## OUTPUT
- Security test scripts
- OWASP checklist completed
- Vulnerability report template
- Remediation recommendations
```

---

# PHASE 9: DOCUMENTATION & TRAINING

## PROMPT 9.1: User Manual - Vietnamese

```markdown
# CONTEXT
Tạo User Manual tiếng Việt cho hệ thống ERP mỹ phẩm.

## YÊU CẦU

### 1. Document Structure

```
HƯỚNG DẪN SỬ DỤNG HỆ THỐNG ERP
================================

Mục lục:
1. Giới thiệu
2. Đăng nhập & Tài khoản
3. Quản lý Nguyên vật liệu
4. Quản lý Nhà cung cấp
5. Mua hàng (PR/PO)
6. Quản lý Kho (WMS)
7. Sản xuất
8. Bán hàng
9. Marketing
10. Báo cáo
11. Quản trị hệ thống
12. Câu hỏi thường gặp
```

### 2. Chapter Template

```markdown
## [TÊN MODULE]

### Giới thiệu
- Mục đích của module
- Các chức năng chính

### Màn hình chính
[Screenshot với chú thích]

### Các thao tác cơ bản

#### Thêm mới [đối tượng]
1. Bước 1: Click nút "Thêm mới"
2. Bước 2: Điền thông tin
   - Trường A (bắt buộc): Mô tả
   - Trường B: Mô tả
3. Bước 3: Click "Lưu"

[Screenshot minh họa]

#### Chỉnh sửa [đối tượng]
...

#### Xóa [đối tượng]
...

### Các tình huống thường gặp

#### Tình huống 1: [Mô tả]
**Vấn đề**: ...
**Giải pháp**: ...

### Lưu ý quan trọng
⚠️ Lưu ý 1: ...
⚠️ Lưu ý 2: ...
```

### 3. Module-Specific Content

#### 3.1 Quản lý Kho (WMS)
```markdown
## QUẢN LÝ KHO (WMS)

### Giới thiệu
Module WMS quản lý toàn bộ hoạt động kho:
- Nhập kho (GRN)
- Xuất kho
- Tồn kho
- Theo dõi Lot/Batch
- Hạn sử dụng

### Nguyên tắc FEFO
⚠️ **QUAN TRỌNG**: Hệ thống áp dụng nguyên tắc FEFO (First Expired First Out)
- Hàng SẮP HẾT HẠN sẽ được xuất TRƯỚC
- Không cần chọn lot thủ công, hệ thống tự động chọn

### Nhập kho (GRN)

#### Tạo phiếu nhập kho từ PO
1. Vào menu: Kho > Nhập kho > Tạo mới
2. Chọn PO cần nhập
3. Hệ thống load danh sách items từ PO
4. Nhập thông tin cho từng dòng:
   - Số lượng thực nhận
   - Lot nhà cung cấp
   - Ngày sản xuất
   - Hạn sử dụng ⚠️ BẮT BUỘC
   - Vị trí lưu kho
5. Click "Lưu" → Trạng thái: Chờ QC
6. Sau khi QC pass → Click "Hoàn thành"
7. Hàng tự động vào kho chính

[Screenshot: Màn hình tạo GRN]

#### Kiểm tra hàng sắp hết hạn
1. Vào menu: Kho > Báo cáo > Hàng sắp hết hạn
2. Chọn số ngày (mặc định: 90 ngày)
3. Xem danh sách lots sắp hết hạn
4. Ưu tiên xuất những lot này trước

[Screenshot: Báo cáo hàng sắp hết hạn]
```

#### 3.2 Sản xuất (BOM)
```markdown
## QUẢN LÝ SẢN XUẤT

### Bảo mật công thức (BOM)
⚠️ **LƯU Ý BẢO MẬT**: 
- Công thức sản phẩm là TÀI SẢN MẬT của công ty
- Chỉ người có quyền mới xem được chi tiết công thức
- KHÔNG chụp màn hình, in, hoặc sao chép công thức
- Mọi truy cập được ghi log

### Quyền xem BOM
| Vai trò | Xem nguyên liệu | Xem số lượng | Xem công thức |
|---------|-----------------|--------------|---------------|
| Staff | ✓ | ✗ | ✗ |
| Production Manager | ✓ | ✓ | ✗ |
| R&D Manager | ✓ | ✓ | ✓ |

### Tạo lệnh sản xuất (Work Order)
1. Vào menu: Sản xuất > Lệnh sản xuất > Tạo mới
2. Chọn sản phẩm cần sản xuất
3. Nhập số lượng kế hoạch
4. Hệ thống tự động tính nguyên liệu cần
5. Kiểm tra tồn kho đủ không
6. Click "Tạo" → Trạng thái: Đã lên kế hoạch
7. Click "Phát hành" → Sẵn sàng sản xuất
8. Click "Bắt đầu" → Nguyên liệu được giữ kho
9. Xuất nguyên liệu cho từng công đoạn
10. Hoàn thành → Nhập thành phẩm

### Truy xuất nguồn gốc
Từ thành phẩm → Nguyên liệu:
1. Vào: Sản xuất > Truy xuất > Truy xuất ngược
2. Nhập mã lot thành phẩm
3. Xem danh sách nguyên liệu đã dùng

Từ nguyên liệu → Thành phẩm:
1. Vào: Sản xuất > Truy xuất > Truy xuất xuôi
2. Nhập mã lot nguyên liệu
3. Xem các thành phẩm đã sử dụng lot này
```

### 4. FAQ Section

```markdown
## CÂU HỎI THƯỜNG GẶP

### Đăng nhập
**Q: Quên mật khẩu?**
A: Click "Quên mật khẩu" → Nhập email → Check email để reset.

**Q: Tài khoản bị khóa?**
A: Tài khoản bị khóa sau 5 lần nhập sai. Liên hệ Admin để mở khóa.

### Kho
**Q: Tại sao không chọn được lot khi xuất kho?**
A: Hệ thống tự động chọn theo FEFO. Lot sắp hết hạn được ưu tiên.

**Q: Hàng đã nhập nhưng không thấy tồn kho?**
A: Kiểm tra GRN đã được QC duyệt và hoàn thành chưa.

### Sản xuất
**Q: Không tạo được lệnh sản xuất?**
A: Kiểm tra sản phẩm có BOM đã được duyệt chưa.

**Q: Không xem được công thức?**
A: Bạn không có quyền xem công thức. Liên hệ R&D Manager.
```

## OUTPUT
- User Manual hoàn chỉnh dạng DOCX (hoặc Markdown)
- Có screenshots mockup
- Ngôn ngữ tiếng Việt, dễ hiểu
- Format cho printing
```

---

## PROMPT 9.2: API Documentation (OpenAPI)

```markdown
# CONTEXT
Tạo API Documentation theo chuẩn OpenAPI 3.0 cho developers.

## YÊU CẦU

### 1. OpenAPI Specification Structure

```yaml
openapi: 3.0.3
info:
  title: ERP Cosmetics API
  description: |
    API documentation cho hệ thống ERP mỹ phẩm.
    
    ## Authentication
    Sử dụng JWT Bearer token.
    
    ## Rate Limiting
    - Authenticated: 1000 req/min
    - Unauthenticated: 100 req/min
    
    ## Error Codes
    - 400: Bad Request
    - 401: Unauthorized
    - 403: Forbidden
    - 404: Not Found
    - 429: Rate Limited
    - 500: Internal Server Error
    
  version: 1.0.0
  contact:
    email: dev@company.vn
    
servers:
  - url: https://erp.company.vn/api/v1
    description: Production
  - url: https://staging-erp.company.vn/api/v1
    description: Staging
  - url: http://localhost:8080/api/v1
    description: Development

tags:
  - name: Auth
    description: Authentication & Authorization
  - name: Users
    description: User Management
  - name: Materials
    description: Material Master Data
  - name: WMS
    description: Warehouse Management
  - name: Manufacturing
    description: Production Management
```

### 2. Auth Endpoints

```yaml
paths:
  /auth/login:
    post:
      tags: [Auth]
      summary: User Login
      description: Authenticate user and return JWT tokens
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [email, password]
              properties:
                email:
                  type: string
                  format: email
                  example: user@company.vn
                password:
                  type: string
                  format: password
                  example: Password@123
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LoginResponse'
        '401':
          description: Invalid credentials
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '423':
          description: Account locked
```

### 3. WMS Endpoints

```yaml
  /stock:
    get:
      tags: [WMS]
      summary: Get Stock List
      description: |
        Lấy danh sách tồn kho với pagination và filters.
        
        **Permissions required**: `wms:stock:read`
      security:
        - bearerAuth: []
      parameters:
        - name: warehouse_id
          in: query
          schema:
            type: string
            format: uuid
        - name: material_id
          in: query
          schema:
            type: string
            format: uuid
        - name: has_stock
          in: query
          description: Chỉ lấy vị trí có tồn kho
          schema:
            type: boolean
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StockListResponse'

  /goods-issue:
    post:
      tags: [WMS]
      summary: Issue Goods (FEFO)
      description: |
        Xuất kho theo nguyên tắc FEFO.
        Hệ thống tự động chọn lots sắp hết hạn trước.
        
        **Permissions required**: `wms:issue:create`
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/GoodsIssueRequest'
      responses:
        '200':
          description: Issue successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GoodsIssueResponse'
              example:
                movement_number: "MOV-OUT-2024-00001"
                issued_from_lots:
                  - lot_number: "LOT-2023-100"
                    quantity: 20
                    expiry_date: "2025-06-30"
                  - lot_number: "LOT-2024-001"
                    quantity: 5
                    expiry_date: "2026-01-14"
        '400':
          description: Insufficient stock
```

### 4. Manufacturing Endpoints

```yaml
  /boms/{id}:
    get:
      tags: [Manufacturing]
      summary: Get BOM Details
      description: |
        Lấy chi tiết BOM.
        
        **Response tùy thuộc vào permission:**
        - `manufacturing:bom:read`: Danh sách nguyên liệu (không số lượng)
        - `manufacturing:bom:quantity_view`: Có số lượng
        - `manufacturing:bom:formula_view`: Có công thức chi tiết
        
        **Note**: Mọi truy cập được ghi audit log.
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                oneOf:
                  - $ref: '#/components/schemas/BOMBasic'
                  - $ref: '#/components/schemas/BOMWithQuantity'
                  - $ref: '#/components/schemas/BOMFull'

  /traceability/backward/{lot_id}:
    get:
      tags: [Manufacturing]
      summary: Backward Traceability
      description: |
        Truy xuất ngược: Từ lot thành phẩm → các lot nguyên liệu đã sử dụng.
        
        Dùng để truy xuất nguồn gốc khi có vấn đề chất lượng.
      security:
        - bearerAuth: []
      parameters:
        - name: lot_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BackwardTraceResult'
```

### 5. Component Schemas

```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  
  schemas:
    Error:
      type: object
      properties:
        code:
          type: string
          example: "INVALID_CREDENTIALS"
        message:
          type: string
          example: "Email or password is incorrect"
    
    Pagination:
      type: object
      properties:
        page:
          type: integer
        limit:
          type: integer
        total:
          type: integer
        total_pages:
          type: integer
    
    Stock:
      type: object
      properties:
        location:
          $ref: '#/components/schemas/Location'
        material:
          $ref: '#/components/schemas/MaterialRef'
        lot:
          $ref: '#/components/schemas/LotRef'
        quantity:
          type: number
        reserved_quantity:
          type: number
        available_quantity:
          type: number
        uom:
          type: string
    
    GoodsIssueRequest:
      type: object
      required: [issue_type, items]
      properties:
        issue_date:
          type: string
          format: date
        issue_type:
          type: string
          enum: [PRODUCTION, SALES, SAMPLE, ADJUSTMENT]
        reference_id:
          type: string
          format: uuid
        reference_number:
          type: string
        items:
          type: array
          items:
            type: object
            required: [material_id, quantity]
            properties:
              material_id:
                type: string
                format: uuid
              quantity:
                type: number
                minimum: 0.0001
```

## OUTPUT
- Complete OpenAPI 3.0 YAML file
- Có thể import vào Swagger UI
- Có thể generate client code
```

---

## PROMPT 9.3: Training Materials & Slides

```markdown
# CONTEXT
Tạo tài liệu training cho từng department.

## YÊU CẦU

### 1. Training Schedule

| Session | Department | Duration | Topics |
|---------|------------|----------|--------|
| 1 | All | 1h | Overview, Login, Navigation |
| 2 | Warehouse | 2h | WMS, GRN, Stock, FEFO |
| 3 | Procurement | 1.5h | Suppliers, PR, PO |
| 4 | Production | 2h | BOM, Work Orders, QC |
| 5 | Sales | 1h | Customers, Orders |
| 6 | Admin | 1h | Users, Roles, Settings |

### 2. Slide Template

```markdown
# [MODULE NAME] Training

## Agenda
1. Mục đích module
2. Demo thực hành
3. Bài tập
4. Q&A

---

## Slide 1: Giới thiệu
- Module này làm gì?
- Ai sử dụng?
- Liên kết với modules khác

---

## Slide 2: Màn hình chính
[Screenshot]
- Giải thích từng phần

---

## Slide 3-N: Các chức năng
[Demo từng chức năng]

---

## Slide: Demo thực hành
Bài tập: [Mô tả bài tập]

---

## Slide: Lưu ý quan trọng
⚠️ Những điểm cần nhớ

---

## Slide: Q&A
Câu hỏi?
```

### 3. WMS Training Content

```markdown
# TRAINING: QUẢN LÝ KHO (WMS)

## Mục tiêu
Sau khóa học, học viên có thể:
- Tạo phiếu nhập kho (GRN)
- Hiểu nguyên tắc FEFO
- Kiểm tra tồn kho
- Xử lý hàng sắp hết hạn

## Nội dung

### Phần 1: Tổng quan (15 phút)
- Warehouse hierarchy: Kho → Zone → Location
- Khái niệm Lot/Batch
- Nguyên tắc FEFO vs FIFO

### Phần 2: Nhập kho (30 phút)
- Demo tạo GRN từ PO
- Nhập thông tin Lot
- QC và hoàn thành GRN
- **Bài tập**: Tạo 1 GRN

### Phần 3: Xuất kho (30 phút)
- Demo xuất kho cho sản xuất
- Giải thích FEFO tự động
- Xem lots đã xuất
- **Bài tập**: Xuất nguyên liệu

### Phần 4: Báo cáo (15 phút)
- Tồn kho theo material
- Hàng sắp hết hạn
- Lịch sử xuất nhập

### Phần 5: Q&A (15 phút)

## Bài tập tổng hợp
1. Nhập kho 100kg Vitamin C, lot VC-2024-001, HSD: 31/12/2025
2. Nhập kho 50kg Vitamin C, lot VC-2024-002, HSD: 30/06/2025
3. Xuất 80kg Vitamin C cho sản xuất
4. Kiểm tra: Lot nào được xuất trước? (Expected: VC-2024-002)
```

### 4. Hands-on Lab Guide

```markdown
# LAB: WMS FEFO Practice

## Scenario
Bạn là Warehouse Staff. Thực hiện các thao tác sau:

## Task 1: Nhập kho nguyên liệu

### Setup
- PO-2024-001 đã được confirmed
- Supplier giao hàng: 100kg Hyaluronic Acid

### Steps
1. Login với account warehouse@company.vn
2. Vào Kho > Nhập kho > Tạo mới
3. Chọn PO-2024-001
4. Nhập:
   - Số lượng nhận: 100
   - Lot NCC: SUP-HA-2024-001
   - NSX: 01/01/2024
   - HSD: 31/12/2025
   - Vị trí: A01-R01-S01
5. Lưu và Hoàn thành

### Verify
- [ ] GRN status = Completed
- [ ] Stock tăng 100kg
- [ ] Lot được tạo

## Task 2: Xuất kho FEFO

### Setup
Trong kho có 2 lots:
- LOT-001: 50kg, HSD 30/06/2025
- LOT-002: 100kg, HSD 31/12/2025

### Steps
1. Vào Kho > Xuất kho
2. Chọn Material: Hyaluronic Acid
3. Số lượng: 70kg
4. Click Xuất

### Verify
- [ ] Lot xuất: LOT-001 (50kg) + LOT-002 (20kg)
- [ ] LOT-001 hết (HSD sớm hơn)
- [ ] LOT-002 còn 80kg
```

## OUTPUT
- PowerPoint slides cho từng module
- Lab guide documents
- Quiz/Assessment questions
- Training attendance sheet
```

---

# PHASE 10: DEPLOYMENT & GO-LIVE

## PROMPT 10.1: Production Server Setup

```markdown
# CONTEXT
Setup production server cho ERP system.

## Server Specs
- OS: Ubuntu 22.04 LTS
- CPU: 16 cores
- RAM: 64GB
- Disk: 1TB SSD

## YÊU CẦU

### 1. Server Hardening Script

```bash
#!/bin/bash
# server-setup.sh

# Update system
apt update && apt upgrade -y

# Install required packages
apt install -y \
  docker.io \
  docker-compose-v2 \
  nginx \
  certbot \
  python3-certbot-nginx \
  fail2ban \
  ufw \
  htop \
  ncdu

# Configure firewall
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw enable

# Configure fail2ban
cat > /etc/fail2ban/jail.local << EOF
[sshd]
enabled = true
maxretry = 3
bantime = 3600

[nginx-http-auth]
enabled = true
EOF

systemctl enable fail2ban
systemctl start fail2ban

# Docker post-install
usermod -aG docker $USER

# Create app user
useradd -m -s /bin/bash erp
usermod -aG docker erp

# Create directories
mkdir -p /opt/erp
mkdir -p /var/log/erp
mkdir -p /backups

chown -R erp:erp /opt/erp
chown -R erp:erp /var/log/erp
chown -R erp:erp /backups

echo "Server setup completed!"
```

### 2. SSL Certificate Setup

```bash
#!/bin/bash
# setup-ssl.sh

DOMAIN="erp.company.vn"

# Get certificate
certbot --nginx -d $DOMAIN -d www.$DOMAIN \
  --non-interactive \
  --agree-tos \
  --email admin@company.vn

# Auto-renewal cron
echo "0 0 1 * * certbot renew --quiet" | crontab -
```

### 3. Nginx Production Config

```nginx
# /etc/nginx/sites-available/erp

upstream api_gateway {
    server 127.0.0.1:8080;
    keepalive 32;
}

upstream frontend {
    server 127.0.0.1:3000;
}

# Rate limiting
limit_req_zone $binary_remote_addr zone=api:10m rate=100r/s;
limit_req_zone $binary_remote_addr zone=login:10m rate=5r/s;

server {
    listen 80;
    server_name erp.company.vn;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name erp.company.vn;

    # SSL
    ssl_certificate /etc/letsencrypt/live/erp.company.vn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/erp.company.vn/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Gzip
    gzip on;
    gzip_types text/plain application/json application/javascript text/css;

    # Frontend
    location / {
        proxy_pass http://frontend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API
    location /api/ {
        limit_req zone=api burst=50 nodelay;
        
        proxy_pass http://api_gateway;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Login rate limit
    location /api/v1/auth/login {
        limit_req zone=login burst=5 nodelay;
        proxy_pass http://api_gateway;
    }

    # Health check
    location /health {
        proxy_pass http://api_gateway;
        access_log off;
    }
}
```

### 4. Docker Compose Production

```yaml
# docker-compose.prod.yml
version: '3.9'

services:
  api-gateway:
    image: erp/api-gateway:${VERSION:-latest}
    restart: always
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
    environment:
      - GIN_MODE=release
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "5"

  # Similar for other services...

  postgres:
    image: postgres:16-alpine
    restart: always
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 8G
    volumes:
      - /data/postgres:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

### 5. Backup Cron Jobs

```bash
# /etc/cron.d/erp-backup

# Database backup - daily at 2 AM
0 2 * * * erp /opt/erp/scripts/backup-db.sh >> /var/log/erp/backup.log 2>&1

# File backup - daily at 3 AM
0 3 * * * erp /opt/erp/scripts/backup-files.sh >> /var/log/erp/backup.log 2>&1

# Cleanup old backups - weekly
0 4 * * 0 erp /opt/erp/scripts/cleanup-backups.sh >> /var/log/erp/backup.log 2>&1
```

### 6. Deployment Script

```bash
#!/bin/bash
# deploy.sh

set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: ./deploy.sh <version>"
    exit 1
fi

echo "Deploying version $VERSION..."

# Pull new images
docker compose -f docker-compose.prod.yml pull

# Backup current state
./scripts/backup-db.sh

# Rolling update
docker compose -f docker-compose.prod.yml up -d --no-deps api-gateway
sleep 10
docker compose -f docker-compose.prod.yml up -d --no-deps auth-service
sleep 10
# ... other services

# Health check
./scripts/health-check.sh

echo "Deployment completed!"
```

## OUTPUT
- Server setup scripts
- Nginx production config
- Docker compose production file
- Backup scripts
- Deployment automation
```

---

## PROMPT 10.2: Go-Live Checklist & Runbook

```markdown
# CONTEXT
Tạo Go-Live checklist và runbook cho deployment.

## YÊU CẦU

### 1. Pre Go-Live Checklist

```markdown
# PRE GO-LIVE CHECKLIST

## 1 Week Before

### Infrastructure
- [ ] Production server provisioned
- [ ] SSL certificates installed
- [ ] Firewall configured
- [ ] Backup system tested
- [ ] Monitoring configured

### Application
- [ ] All services deployed to staging
- [ ] Staging tested by QA
- [ ] Performance test passed
- [ ] Security audit passed
- [ ] UAT sign-off received

### Data
- [ ] Data migration scripts tested
- [ ] Data validation passed
- [ ] Opening balances verified
- [ ] Master data imported

### People
- [ ] Training completed for all users
- [ ] Support team briefed
- [ ] Escalation path defined
- [ ] On-call schedule set

## 1 Day Before

### Final Checks
- [ ] All services healthy on staging
- [ ] Backup taken
- [ ] Rollback plan reviewed
- [ ] Communication sent to users
- [ ] Support channels ready

### Go/No-Go Decision
- [ ] All checklist items completed
- [ ] Sign-off from: IT Lead, Business Owner, QA Lead
```

### 2. Go-Live Runbook

```markdown
# GO-LIVE RUNBOOK

## Timeline: Saturday 6:00 AM - 12:00 PM

### Phase 1: Preparation (6:00 - 6:30)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 6:00 | Team standup call | Tech Lead | |
| 6:10 | Verify all team members online | Tech Lead | |
| 6:15 | Lock old system (read-only) | DBA | |
| 6:20 | Take final backup of old system | DBA | |
| 6:30 | Confirm ready to proceed | All | |

### Phase 2: Data Migration (6:30 - 8:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 6:30 | Export data from old system | DBA | |
| 7:00 | Run migration scripts | Dev Lead | |
| 7:30 | Validate data counts | QA | |
| 7:45 | Validate critical data | Business | |
| 8:00 | Data migration sign-off | All | |

### Phase 3: Application Deployment (8:00 - 9:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 8:00 | Deploy all services | DevOps | |
| 8:15 | Run database migrations | Dev Lead | |
| 8:30 | Configure production env | DevOps | |
| 8:45 | Health check all services | DevOps | |
| 9:00 | Deployment complete | Tech Lead | |

### Phase 4: Verification (9:00 - 10:30)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 9:00 | Smoke test - Login | QA | |
| 9:10 | Smoke test - Master Data | QA | |
| 9:20 | Smoke test - WMS | QA | |
| 9:30 | Smoke test - Manufacturing | QA | |
| 9:40 | Smoke test - Sales | QA | |
| 9:50 | End-to-end test flow | QA | |
| 10:15 | Performance spot check | DevOps | |
| 10:30 | Verification sign-off | QA Lead | |

### Phase 5: Go-Live (10:30 - 11:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 10:30 | Enable user access | Admin | |
| 10:35 | Monitor first user logins | DevOps | |
| 10:45 | Monitor system metrics | DevOps | |
| 11:00 | Official go-live announcement | PM | |

### Phase 6: Hypercare (11:00 - 12:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 11:00 | Support channel active | Support | |
| 11:30 | First hour report | Tech Lead | |
| 12:00 | Handover to support team | Tech Lead | |

## Rollback Procedure

### Trigger Conditions
- Critical bug affecting > 50% users
- Data corruption detected
- System performance < 50% target
- Security breach detected

### Rollback Steps
1. Announce system maintenance
2. Disable user access
3. Stop all services
4. Restore database backup
5. Revert to old system
6. Verify old system working
7. Communicate to users

### Rollback Decision
- Must be made by: Tech Lead + Business Owner
- Deadline: Before 11:00 AM
```

### 3. Incident Response Template

```markdown
# INCIDENT RESPONSE TEMPLATE

## Incident Details
- **ID**: INC-[YYYYMMDD]-[NNN]
- **Severity**: P1/P2/P3/P4
- **Status**: Open/Investigating/Resolved/Closed
- **Reported by**: 
- **Reported at**: 
- **Resolved at**: 

## Description
[What happened?]

## Impact
- Users affected: 
- Functions affected: 
- Business impact: 

## Timeline
| Time | Event |
|------|-------|
| HH:MM | Incident detected |
| HH:MM | Team notified |
| HH:MM | Investigation started |
| HH:MM | Root cause identified |
| HH:MM | Fix deployed |
| HH:MM | Incident resolved |

## Root Cause
[Technical explanation]

## Resolution
[What was done to fix?]

## Action Items
- [ ] Preventive measure 1
- [ ] Preventive measure 2

## Lessons Learned
[What can we do better?]
```

## OUTPUT
- Pre go-live checklist
- Go-live runbook with timeline
- Rollback procedure
- Incident response templates
- Communication templates
```

---

# PHASE 11: MONITORING & MAINTENANCE

## PROMPT 11.1: Monitoring & Alerting Setup

```markdown
# CONTEXT
Setup monitoring và alerting cho production.

## Tech Stack
- Prometheus (metrics)
- Grafana (visualization)
- Loki (logs)
- AlertManager (alerts)

## YÊU CẦU

### 1. Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - /etc/prometheus/rules/*.yml

scrape_configs:
  - job_name: 'api-gateway'
    static_configs:
      - targets: ['api-gateway:8080']
    metrics_path: /metrics

  - job_name: 'auth-service'
    static_configs:
      - targets: ['auth-service:8081']

  - job_name: 'wms-service'
    static_configs:
      - targets: ['wms-service:8086']

  # ... other services

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
```

### 2. Alert Rules

```yaml
# alerts/erp-alerts.yml
groups:
  - name: erp-critical
    rules:
      # Service down
      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.job }} is down"
          description: "{{ $labels.job }} has been down for more than 1 minute"

      # High error rate
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m])) by (service)
          /
          sum(rate(http_requests_total[5m])) by (service)
          > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate on {{ $labels.service }}"
          description: "Error rate is {{ $value | humanizePercentage }}"

      # Database connection issues
      - alert: DatabaseConnectionHigh
        expr: pg_stat_activity_count > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High database connections"
          description: "{{ $value }} active connections"

  - name: erp-business
    rules:
      # Low stock alert (business metric)
      - alert: LowStockAlert
        expr: wms_stock_available_quantity < wms_material_reorder_point
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "Low stock for {{ $labels.material_code }}"
          description: "Available: {{ $value }}, Reorder point: {{ $labels.reorder_point }}"

      # Lots expiring soon
      - alert: LotsExpiringSoon
        expr: wms_lots_expiring_count{days="30"} > 0
        for: 1d
        labels:
          severity: warning
        annotations:
          summary: "{{ $value }} lots expiring in 30 days"

      # Certificate expiring
      - alert: CertificateExpiring
        expr: supplier_certification_days_until_expiry < 30
        for: 1d
        labels:
          severity: warning
        annotations:
          summary: "Supplier certificate expiring"
          description: "{{ $labels.supplier_name }} - {{ $labels.cert_type }} expires in {{ $value }} days"
```

### 3. Grafana Dashboards

```json
// dashboards/erp-overview.json
{
  "title": "ERP Overview",
  "panels": [
    {
      "title": "Service Health",
      "type": "stat",
      "targets": [
        {
          "expr": "count(up == 1)",
          "legendFormat": "Healthy Services"
        }
      ]
    },
    {
      "title": "Request Rate",
      "type": "graph",
      "targets": [
        {
          "expr": "sum(rate(http_requests_total[5m])) by (service)",
          "legendFormat": "{{ service }}"
        }
      ]
    },
    {
      "title": "Response Time (p95)",
      "type": "graph",
      "targets": [
        {
          "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
          "legendFormat": "p95"
        }
      ]
    },
    {
      "title": "Error Rate",
      "type": "graph",
      "targets": [
        {
          "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))",
          "legendFormat": "Error %"
        }
      ]
    }
  ]
}
```

### 4. AlertManager Configuration

```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alerts@company.vn'
  smtp_auth_username: 'alerts@company.vn'
  smtp_auth_password: '<app-password>'

route:
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'default'
  
  routes:
    - match:
        severity: critical
      receiver: 'critical-alerts'
      repeat_interval: 1h
    
    - match:
        severity: warning
      receiver: 'warning-alerts'
      repeat_interval: 4h

receivers:
  - name: 'default'
    email_configs:
      - to: 'it-team@company.vn'

  - name: 'critical-alerts'
    email_configs:
      - to: 'it-lead@company.vn,cto@company.vn'
    # Telegram/Slack webhook có thể thêm ở đây

  - name: 'warning-alerts'
    email_configs:
      - to: 'it-team@company.vn'
```

### 5. Log Aggregation (Loki)

```yaml
# loki-config.yml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1

schema_config:
  configs:
    - from: 2024-01-01
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: filesystem
  filesystem:
    directory: /loki/chunks

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h

# Promtail config for log shipping
# promtail.yml
positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: erp-services
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        target_label: 'container'
```

## OUTPUT
- Prometheus configuration
- Alert rules (critical + business)
- Grafana dashboard JSON
- AlertManager configuration
- Loki + Promtail configuration
```

---

## PROMPT 11.2: Maintenance Procedures

```markdown
# CONTEXT
Tạo maintenance procedures cho operations team.

## YÊU CẦU

### 1. Daily Operations Checklist

```markdown
# DAILY OPERATIONS CHECKLIST

## Morning Check (8:00 AM)

### System Health
- [ ] All services UP (Grafana dashboard)
- [ ] No critical alerts overnight
- [ ] Response time normal (<200ms p95)
- [ ] Error rate < 0.1%
- [ ] Disk usage < 80%

### Application Health
- [ ] Login working
- [ ] API Gateway responding
- [ ] Database connections normal
- [ ] Redis cache healthy
- [ ] NATS queue healthy

### Business Alerts
- [ ] Check low stock alerts
- [ ] Check expiring lots (30 days)
- [ ] Check pending approvals backlog

## Evening Check (5:00 PM)
- [ ] Review day's incidents
- [ ] Verify backup completed
- [ ] Check alert queue
- [ ] Update operations log
```

### 2. Weekly Maintenance

```markdown
# WEEKLY MAINTENANCE

## Every Monday (9:00 AM)

### Performance Review
- [ ] Review weekly performance metrics
- [ ] Identify slow queries
- [ ] Check database growth
- [ ] Review error logs

### Security
- [ ] Review failed login attempts
- [ ] Check for suspicious activities
- [ ] Review audit logs
- [ ] Update blocked IPs if needed

### Housekeeping
- [ ] Clean old logs (>30 days)
- [ ] Vacuum PostgreSQL
- [ ] Clear Redis cache (if needed)
- [ ] Archive old data (if needed)
```

### 3. Monthly Maintenance

```markdown
# MONTHLY MAINTENANCE

## First Saturday (2:00 AM - 6:00 AM)

### Planned Downtime Tasks
- [ ] Apply security patches
- [ ] Update Docker images
- [ ] Database maintenance
- [ ] SSL certificate check
- [ ] Backup verification test

### Capacity Planning
- [ ] Review resource utilization
- [ ] Forecast growth
- [ ] Plan scaling if needed

### Documentation
- [ ] Update runbooks if needed
- [ ] Review and update procedures
- [ ] Update contact lists
```

### 4. Backup Verification Procedure

```markdown
# BACKUP VERIFICATION PROCEDURE

## Monthly Backup Test

### Purpose
Verify backups can be restored successfully.

### Procedure

1. **Prepare Test Environment**
   ```bash
   # Start test database
   docker run -d --name pg-restore-test \
     -e POSTGRES_PASSWORD=test \
     postgres:16-alpine
   ```

2. **Restore Latest Backup**
   ```bash
   # Get latest backup
   LATEST_BACKUP=$(ls -t /backups/db/*.sql.gz | head -1)
   
   # Restore
   gunzip -c $LATEST_BACKUP | docker exec -i pg-restore-test psql -U postgres
   ```

3. **Verify Data**
   ```sql
   -- Check record counts
   SELECT 'users', COUNT(*) FROM users
   UNION ALL
   SELECT 'materials', COUNT(*) FROM materials
   UNION ALL
   SELECT 'stock', COUNT(*) FROM stock;
   
   -- Compare with production
   ```

4. **Document Results**
   - Backup file: 
   - Backup date: 
   - Restore time: 
   - Data verification: PASS/FAIL
   - Issues found: 

5. **Cleanup**
   ```bash
   docker rm -f pg-restore-test
   ```

### Sign-off
- Tested by: 
- Date: 
- Result: 
```

### 5. Incident Management

```markdown
# INCIDENT MANAGEMENT PROCEDURE

## Severity Levels

| Level | Description | Response Time | Examples |
|-------|-------------|---------------|----------|
| P1 | System down | 15 min | All services down, data loss |
| P2 | Major impact | 30 min | Key function unavailable |
| P3 | Minor impact | 4 hours | Non-critical bug |
| P4 | Low impact | 24 hours | Cosmetic issues |

## Escalation Path

```
P1: On-call → Tech Lead → CTO → CEO
P2: On-call → Tech Lead → CTO
P3: On-call → Tech Lead
P4: Support Team
```

## Incident Response Steps

### 1. Detection
- Alert received OR user reported
- Log incident ticket

### 2. Triage
- Assess severity
- Notify relevant team
- Start incident channel

### 3. Investigation
- Gather logs
- Identify affected components
- Find root cause

### 4. Resolution
- Implement fix
- Test fix
- Deploy to production

### 5. Communication
- Update stakeholders
- Post-mortem if P1/P2

### 6. Documentation
- Complete incident report
- Update runbooks
- Create follow-up tasks
```

### 6. Common Troubleshooting

```markdown
# TROUBLESHOOTING GUIDE

## Service Won't Start

### Symptoms
- Container exits immediately
- Health check failing

### Steps
1. Check logs: `docker logs <service>`
2. Check config: Environment variables
3. Check dependencies: Database, Redis, NATS
4. Check ports: Port conflicts
5. Check resources: Memory, disk

## Database Connection Issues

### Symptoms
- "connection refused" errors
- Timeout errors

### Steps
1. Check PostgreSQL running: `docker ps | grep postgres`
2. Check connection count: `SELECT count(*) FROM pg_stat_activity;`
3. Check pg_hba.conf if auth issues
4. Restart connection pool

## High Memory Usage

### Symptoms
- OOM killer triggered
- Service restarts frequently

### Steps
1. Check memory: `docker stats`
2. Check for memory leaks
3. Review recent deployments
4. Increase memory limit temporarily
5. Fix root cause

## Slow Queries

### Symptoms
- High response time
- Database CPU high

### Steps
1. Find slow queries: `pg_stat_statements`
2. Check missing indexes
3. Analyze query plan: `EXPLAIN ANALYZE`
4. Add index or optimize query
```

## OUTPUT
- Daily/Weekly/Monthly checklists
- Backup verification procedure
- Incident management guide
- Troubleshooting guide
- Operations runbook
```

---

## 📋 SUMMARY: POST-IMPLEMENTATION PROMPTS

| Phase | Prompt | Purpose |
|-------|--------|---------|
| **7.1** | Unit Tests - Auth/User | Test authentication, authorization |
| **7.2** | Unit Tests - WMS | Test FEFO logic (CRITICAL) |
| **7.3** | Unit Tests - Manufacturing | Test BOM encryption, traceability |
| **7.4** | Integration Tests | API end-to-end tests |
| **7.5** | Performance Tests | Load testing với K6 |
| **7.6** | Security Tests | OWASP, penetration tests |
| **9.1** | User Manual | Vietnamese documentation |
| **9.2** | API Documentation | OpenAPI 3.0 spec |
| **9.3** | Training Materials | Slides, lab guides |
| **10.1** | Production Setup | Server, SSL, deployment |
| **10.2** | Go-Live Runbook | Checklist, procedures |
| **11.1** | Monitoring Setup | Prometheus, Grafana, alerts |
| **11.2** | Maintenance | Operations procedures |

---

**Created**: 2026-01-25  
**Purpose**: Post-implementation phases cho ERP Cosmetics
