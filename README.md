# HỆ THỐNG ERP MỸ PHẨM - MICROSERVICES ARCHITECTURE

## TỔNG QUAN DỰ ÁN

Hệ thống ERP toàn diện cho công ty sản xuất mỹ phẩm thiên nhiên tại Việt Nam, được xây dựng theo kiến trúc Microservices với khả năng hoạt động offline, phù hợp cho môi trường on-premise.

### Đặc điểm chính

- **Quy mô**: ~100 nhân viên concurrent
- **Sản phẩm**: Serum, kem dưỡng, dầu gội, sữa tắm từ nguyên liệu thiên nhiên
- **SKU**: 10,000 - 50,000 sản phẩm
- **Transactions**: 100,000 - 500,000 giao dịch/tháng
- **Deployment**: On-premise với khả năng offline
- **AI/ML**: Tích hợp Ollama cho dự báo và phân tích

### Đặc thù ngành mỹ phẩm

✅ **Truy xuất nguồn gốc (Traceability)**: Theo dõi từ nguyên liệu → thành phẩm → khách hàng  
✅ **BOM bảo mật**: Công thức sản phẩm với phân quyền chặt chẽ  
✅ **FEFO Logic**: First Expired First Out - ưu tiên hàng sắp hết hạn  
✅ **GMP Compliance**: Quy trình sản xuất theo chuẩn GMP  
✅ **Quản lý chứng nhận**: Theo dõi GMP, ISO, Organic, Ecocert  
✅ **Cold Storage**: Giám sát nhiệt độ 2-8°C cho nguyên liệu đặc biệt  
✅ **KOL Marketing**: Quản lý sample gửi cho KOL/Influencer

## KIẾN TRÚC HỆ THỐNG

### Tech Stack

**Backend (Microservices)**
- Go 1.22+ với Gin framework
- GORM + PostgreSQL (mỗi service có database riêng)
- Redis (cache & session)
- NATS/RabbitMQ (message queue)
- gRPC (internal communication)

**Frontend**
- Vue 3 + TypeScript
- PrimeVue UI Framework
- Pinia (state management)
- Axios (HTTP client)

**Infrastructure**
- Docker + Docker Compose
- Nginx (reverse proxy)
- MinIO (object storage)
- Ollama (AI/LLM local)

**Monitoring**
- Prometheus + Grafana
- Loki (logs)
- Jaeger (distributed tracing)

### Danh sách Services

| # | Service | Mô tả |
|---|---------|-------|
| 1 | API Gateway | Routing, authentication, rate limiting |
| 2 | Auth Service | JWT authentication & RBAC authorization |
| 3 | User Service | Quản lý user, department, employee |
| 4 | Master Data Service | Materials, products, categories, UoM |
| 5 | Supplier Service | Quản lý nhà cung cấp, chứng nhận, ASL |
| 6 | Procurement Service | PR, PO, RFQ workflow |
| 7 | WMS Service | Warehouse, stock, lot, FEFO, GRN/GI |
| 8 | Manufacturing Service | BOM, work order, QC, traceability |
| 9 | Sales Service | Customer, quotation, sales order |
| 10 | Marketing Service | Campaign, KOL, sample tracking |
| 11 | Finance Service | AP, AR, basic accounting |
| 12 | Reporting Service | Reports, dashboards, analytics |
| 13 | Notification Service | Email, push, alerts |
| 14 | AI Service | Forecasting, demand planning |
| 15 | File Service | Document & image storage (MinIO) |

## MỤC LỤC TÀI LIỆU

### Core Architecture
- [01 - Kiến trúc hệ thống](./docs/01-ARCHITECTURE.md)
- [02 - Chi tiết các Services](./docs/02-SERVICE-SPECIFICATIONS.md)
- [13 - API Gateway](./docs/13-API-GATEWAY.md)
- [14 - Event Catalog](./docs/14-EVENT-CATALOG.md)

### Core Services
- [03 - Auth Service](./docs/03-AUTH-SERVICE.md)
- [04 - User Service](./docs/04-USER-SERVICE.md)
- [05 - Master Data Service](./docs/05-MASTER-DATA-SERVICE.md)

### Supply Chain Services
- [06 - Supplier Service](./docs/06-SUPPLIER-SERVICE.md)
- [07 - Procurement Service](./docs/07-PROCUREMENT-SERVICE.md)
- [08 - WMS Service](./docs/08-WMS-SERVICE.md)

### Production & Sales
- [09 - Manufacturing Service](./docs/09-MANUFACTURING-SERVICE.md)
- [10 - Sales Service](./docs/10-SALES-SERVICE.md)
- [11 - Marketing Service](./docs/11-MARKETING-SERVICE.md)

### Supporting Services
- [12 - Notification Service](./docs/12-NOTIFICATION-SERVICE.md)

### Database & Deployment
- [15 - Database Schemas](./docs/15-DATABASE-SCHEMAS.md)
- [16 - Deployment Guide](./docs/16-DEPLOYMENT.md)

### Planning
- [17 - Implementation Roadmap](./docs/17-IMPLEMENTATION-ROADMAP.md)
- [18 - Glossary](./docs/18-GLOSSARY.md)

## QUICK START

### Prerequisites

```bash
# Cài đặt Docker & Docker Compose
docker --version  # >= 24.0
docker compose version  # >= 2.20

# Cài đặt Go (cho development)
go version  # >= 1.22

# Cài đặt Node.js (cho frontend)
node --version  # >= 18.0
npm --version   # >= 9.0
```

### Development Setup

```bash
# Clone repository
git clone <repository-url>
cd erp-cosmetics

# Copy environment variables
cp .env.example .env

# Start all services với Docker Compose
docker compose up -d

# Check services status
docker compose ps

# View logs
docker compose logs -f [service-name]
```

### Accessing Services

| Service | URL | Credentials |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | admin / admin123 |
| API Gateway | http://localhost:8080 | - |
| PgAdmin | http://localhost:5050 | admin@erp.com / admin |
| MinIO Console | http://localhost:9001 | minioadmin / minioadmin |
| Grafana | http://localhost:3001 | admin / admin |

### Running Individual Service (Development)

```bash
# Navigate to service directory
cd services/auth-service

# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start service
go run cmd/main.go

# Run tests
go test ./...
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build
```

## PROJECT STRUCTURE

```
erp-cosmetics/
├── services/                    # Microservices
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── master-data-service/
│   ├── supplier-service/
│   ├── procurement-service/
│   ├── wms-service/
│   ├── manufacturing-service/
│   ├── sales-service/
│   ├── marketing-service/
│   ├── finance-service/
│   ├── reporting-service/
│   ├── notification-service/
│   ├── ai-service/
│   └── file-service/
├── frontend/                    # Vue 3 + PrimeVue
├── docs/                        # Documentation
├── deploy/                      # Deployment configs
│   ├── docker/
│   └── nginx/
├── scripts/                     # Utility scripts
├── docker-compose.yml
├── docker-compose.dev.yml
├── .env.example
└── README.md
```

## DEVELOPMENT WORKFLOW

### Branch Strategy

- `main`: Production-ready code
- `develop`: Development branch
- `feature/*`: New features
- `bugfix/*`: Bug fixes
- `hotfix/*`: Production hotfixes

### Commit Convention

```
feat(service-name): Add new feature
fix(service-name): Fix bug
docs: Update documentation
refactor(service-name): Code refactoring
test(service-name): Add tests
chore: Update dependencies
```

### Code Review Process

1. Create feature branch từ `develop`
2. Implement feature với tests
3. Create Pull Request
4. Code review (min 1 approval)
5. Merge to `develop`
6. Deploy to staging
7. QA testing
8. Merge to `main` và deploy production

## TESTING

### Unit Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific service tests
cd services/auth-service
go test -v ./...
```

### Integration Tests

```bash
# Start test environment
docker compose -f docker-compose.test.yml up -d

# Run integration tests
make integration-test
```

### E2E Tests

```bash
cd frontend
npm run test:e2e
```

## MONITORING & LOGGING

### Logs

```bash
# View all logs
docker compose logs -f

# View specific service
docker compose logs -f auth-service

# View last 100 lines
docker compose logs --tail=100 auth-service
```

### Metrics (Prometheus + Grafana)

- Truy cập Grafana: http://localhost:3001
- Default dashboards đã được import tự động
- Custom metrics endpoint: `/metrics` trên mỗi service

### Tracing (Jaeger)

- Truy cập Jaeger UI: http://localhost:16686
- Theo dõi distributed traces giữa services

## DEPLOYMENT

### Development

```bash
docker compose -f docker-compose.dev.yml up -d
```

### Staging

```bash
docker compose -f docker-compose.staging.yml up -d
```

### Production

Xem chi tiết tại [16-DEPLOYMENT.md](./docs/16-DEPLOYMENT.md)

## SECURITY

### Environment Variables

- **KHÔNG BAO GIỜ** commit file `.env` vào git
- Sử dụng `.env.example` làm template
- Production secrets quản lý qua Docker secrets

### JWT Tokens

- Access token: 15 phút
- Refresh token: 7 ngày
- Rotation tự động khi refresh

### Database

- Mỗi service có database riêng
- Connection pooling được configure
- Automated backups mỗi ngày

## CONTRIBUTING

1. Fork the repository
2. Create feature branch
3. Make your changes
4. Write/update tests
5. Update documentation
6. Submit pull request

## LICENSE

Proprietary - All rights reserved

## SUPPORT

- **Documentation**: Xem thư mục `/docs`
- **Issues**: Tạo issue trên GitHub
- **Email**: dev-team@company.com
- **Slack**: #erp-dev channel

## CHANGELOG

Xem [CHANGELOG.md](./CHANGELOG.md) để biết chi tiết các thay đổi theo version.

---

**Version**: 1.0.0  
**Last Updated**: 2026-01-23  
**Maintainers**: ERP Development Team
