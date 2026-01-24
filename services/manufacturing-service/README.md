# Manufacturing Service

Dịch vụ Sản xuất (Manufacturing Service) - Phần của Giai đoạn 3: Vận hành trong hệ thống ERP mỹ phẩm.

## 📋 Tổng Quan

Manufacturing Service quản lý:
- **BOM (Bill of Materials)**: Công thức sản phẩm với mã hóa AES-256-GCM
- **Work Orders**: Lệnh sản xuất với vòng đời đầy đủ
- **QC (Quality Control)**: Kiểm soát chất lượng IQC/IPQC/FQC
- **NCR**: Báo cáo không phù hợp (Non-Conformance Report)
- **Traceability**: Truy xuất nguồn gốc (ngược/xuôi)

## 🔧 Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin (HTTP), gRPC
- **Database**: PostgreSQL
- **Events**: NATS JetStream
- **Encryption**: AES-256-GCM cho formula_details

## 📊 Database Schema

| Table | Mô tả |
|-------|-------|
| `boms` | Bill of Materials với formula encrypted |
| `bom_line_items` | Thành phần trong BOM |
| `bom_versions` | Lịch sử phiên bản BOM |
| `work_orders` | Lệnh sản xuất |
| `wo_line_items` | Nguyên liệu dự kiến |
| `wo_material_issues` | Nguyên liệu đã xuất (liên kết WMS) |
| `qc_checkpoints` | Mẫu kiểm tra QC |
| `qc_inspections` | Các lần kiểm tra QC |
| `qc_inspection_items` | Chi tiết kết quả kiểm tra |
| `ncrs` | Báo cáo không phù hợp |
| `batch_traceability` | Truy xuất lô hàng |

## 🔐 BOM Security

```
Formula Details được mã hóa AES-256-GCM trước khi lưu vào DB.

Quyền truy cập:
- manufacturing:bom:formula_view - Xem công thức đầy đủ
- manufacturing:bom:quantity_view - Xem số lượng nguyên liệu

Chỉ RD Manager và Production Manager mới được xem formula đầy đủ.
```

## 🔄 Work Order Lifecycle

```
PLANNED → RELEASED → IN_PROGRESS → QC_PENDING → COMPLETED
                                        ↓
                                    CANCELLED
```

## 📡 API Endpoints

### BOM
- `POST /api/v1/boms` - Tạo BOM mới
- `GET /api/v1/boms` - Danh sách BOM
- `GET /api/v1/boms/:id` - Chi tiết BOM
- `POST /api/v1/boms/:id/approve` - Phê duyệt BOM

### Work Orders
- `POST /api/v1/work-orders` - Tạo WO
- `GET /api/v1/work-orders` - Danh sách WO
- `GET /api/v1/work-orders/:id` - Chi tiết WO
- `PATCH /api/v1/work-orders/:id/release` - Release WO
- `PATCH /api/v1/work-orders/:id/start` - Start WO
- `PATCH /api/v1/work-orders/:id/complete` - Complete WO

### QC
- `GET /api/v1/qc-checkpoints` - Danh sách checkpoint
- `POST /api/v1/qc-inspections` - Tạo inspection
- `GET /api/v1/qc-inspections/:id` - Chi tiết inspection
- `PATCH /api/v1/qc-inspections/:id/approve` - Phê duyệt

### NCR
- `POST /api/v1/ncrs` - Tạo NCR
- `GET /api/v1/ncrs` - Danh sách NCR
- `GET /api/v1/ncrs/:id` - Chi tiết NCR
- `PATCH /api/v1/ncrs/:id/close` - Đóng NCR

### Traceability
- `GET /api/v1/traceability/backward/:lot_id` - Truy xuất ngược
- `GET /api/v1/traceability/forward/:lot_id` - Truy xuất xuôi

## 📤 Events Published

| Event | Trigger |
|-------|---------|
| `manufacturing.bom.created` | BOM được tạo |
| `manufacturing.bom.approved` | BOM được duyệt |
| `manufacturing.wo.created` | WO được tạo |
| `manufacturing.wo.started` | WO bắt đầu → WMS reserve materials |
| `manufacturing.wo.completed` | WO hoàn thành → WMS nhận thành phẩm |
| `manufacturing.qc.failed` | QC thất bại |
| `manufacturing.ncr.created` | NCR được tạo |

## 🚀 Chạy Service

```bash
# Development
make run

# Build
make build

# Run migrations
make migrate-up

# Run tests
make test
```

## ⚙️ Environment Variables

```
PORT=8087
GRPC_PORT=9087
DB_HOST=localhost
DB_PORT=5438
DB_NAME=manufacturing_db
BOM_ENCRYPTION_KEY=<32-byte-hex-key>
NATS_URL=nats://localhost:4222
```

## 📁 Project Structure

```
manufacturing-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   │   ├── entity/
│   │   └── repository/
│   ├── infrastructure/
│   │   ├── event/
│   │   └── persistence/postgres/
│   ├── usecase/
│   │   ├── bom/
│   │   ├── workorder/
│   │   ├── qc/
│   │   ├── ncr/
│   │   └── traceability/
│   └── delivery/http/
│       ├── dto/
│       ├── handler/
│       └── router/
├── migrations/
├── Makefile
├── Dockerfile
└── go.mod
```
