# 18 - GLOSSARY

## TỔNG QUAN

Thuật ngữ và từ viết tắt được sử dụng trong hệ thống ERP mỹ phẩm.

---

## GENERAL TERMS

**ERP (Enterprise Resource Planning)**: Hệ thống quản lý tài nguyên doanh nghiệp tích hợp tất cả quy trình kinh doanh.

**Microservices**: Kiến trúc phần mềm chia ứng dụng thành các services nhỏ, độc lập.

**API Gateway**: Điểm vào duy nhất cho tất cả requests từ client đến backend services.

**gRPC**: High-performance RPC framework dùng cho communication giữa các services.

**NATS**: Message queue system cho event-driven communication.

**Docker Compose**: Tool để define và run multi-container Docker applications.

---

## COSMETICS-SPECIFIC TERMS

### INCI (International Nomenclature of Cosmetic Ingredients)
Hệ thống đặt tên chuẩn quốc tế cho nguyên liệu mỹ phẩm. Ví dụ: "Ascorbic Acid" là INCI name của Vitamin C.

### CAS Number (Chemical Abstracts Service)
Số đăng ký duy nhất cho mỗi hóa chất. Ví dụ: CAS 50-81-7 là Vitamin C.

### GMP (Good Manufacturing Practice)
Quy chuẩn sản xuất tốt cho ngành dược phẩm và mỹ phẩm.

### ISO 22716
Tiêu chuẩn quốc tế về GMP cho ngành mỹ phẩm.

### FEFO (First Expired First Out)
Nguyên tắc xuất hàng sắp hết hạn trước, quan trọng cho mỹ phẩm.

### COA (Certificate of Analysis)
Giấy chứng nhận phân tích chất lượng từ nhà cung cấp.

### SDS (Safety Data Sheet)
Bảng dữ liệu an toàn hóa chất.

### Allergen
Chất gây dị ứng, phải được ghi rõ trên nhãn mỹ phẩm.

### Organic Certification
Chứng nhận nguyên liệu/sản phẩm hữu cơ (Ecocert, USDA Organic, Cosmos).

---

## BUSINESS PROCESS TERMS

### BOM (Bill of Materials)
Danh sách nguyên liệu cần thiết để sản xuất 1 sản phẩm, bao gồm công thức và quy trình.

### SKU (Stock Keeping Unit)
Mã định danh sản phẩm dùng để track inventory.

### UoM (Unit of Measure)
Đơn vị đo lường (kg, liter, bottle, piece, etc.).

### PR (Purchase Requisition)
Yêu cầu mua hàng từ department, cần approval trước khi tạo PO.

### RFQ (Request for Quotation)
Yêu cầu báo giá từ supplier.

### PO (Purchase Order)
Đơn đặt hàng gửi cho nhà cung cấp.

### GRN (Goods Receipt Note)
Phiếu nhập kho khi nhận hàng từ supplier.

### WO (Work Order)
Lệnh sản xuất.

### SO (Sales Order)
Đơn hàng bán.

### NCR (Non-Conformance Report)
Báo cáo không phù hợp (khi có vấn đề chất lượng).

---

## QUALITY CONTROL TERMS

### IQC (Incoming Quality Control)
Kiểm tra chất lượng nguyên liệu đầu vào.

### IPQC (In-Process Quality Control)
Kiểm tra chất lượng trong quá trình sản xuất.

### FQC (Final Quality Control)
Kiểm tra chất lượng cuối cùng trước khi xuất kho.

### OQC (Outgoing Quality Control)
Kiểm tra chất lượng hàng xuất trước khi giao cho khách.

### Batch/Lot
Lô hàng có cùng điều kiện sản xuất, dùng để traceability.

### Traceability
Khả năng theo dõi nguyên liệu/sản phẩm từ nguồn gốc đến người tiêu dùng.

---

## WAREHOUSE TERMS

### Warehouse
Kho hàng.

### Zone
Khu vực trong kho (receiving, storage, picking, shipping, quarantine).

### Location
Vị trí cụ thể trong kho (aisle, rack, shelf, bin).

### Stock on Hand
Số lượng tồn kho thực tế.

### Reserved Stock
Số lượng đã được giữ chỗ cho đơn hàng/sản xuất.

### Available Stock
Stock on hand - Reserved stock.

### Cycle Count
Kiểm kho định kỳ (không cần đóng cửa kho).

### Physical Inventory
Kiểm kê tồn kho toàn bộ.

### Cold Storage
Kho lạnh (2-8°C) cho nguyên liệu đặc biệt.

### Quarantine Zone
Khu vực cách ly hàng chờ QC.

---

## TECHNICAL TERMS

### JWT (JSON Web Token)
Token format dùng cho authentication.

### RBAC (Role-Based Access Control)
Kiểm soát truy cập dựa trên role.

### Circuit Breaker
Pattern để handle service failures gracefully.

### Rate Limiting
Giới hạn số requests per user/time để tránh abuse.

### Saga Pattern
Pattern để handle distributed transactions.

### Eventual Consistency
Dữ liệu sẽ nhất quán sau một khoảng thời gian.

### MinIO
S3-compatible object storage.

### Prometheus
Monitoring & alerting system.

### Grafana
Visualization & dashboards platform.

### Loki
Log aggregation system.

### Jaeger
Distributed tracing system.

---

## MARKETING TERMS

### KOL (Key Opinion Leader)
Người có ảnh hưởng (influencer, blogger, celebrity).

### Engagement Rate
Tỷ lệ tương tác (likes + comments / followers * 100%).

### ROI (Return on Investment)
Lợi nhuận trên đầu tư.

### Campaign
Chiến dịch marketing.

### Sample
Sản phẩm mẫu gửi miễn phí để test/review.

### Tier
Phân cấp KOL (MEGA, MACRO, MICRO, NANO dựa trên số followers).

---

## ROLE ABBREVIATIONS

**CFO**: Chief Financial Officer  
**COO**: Chief Operating Officer  
**QC Manager**: Quality Control Manager  
**R&D Manager**: Research & Development Manager  
**HR**: Human Resources

---

## STATUS VALUES

### General
- **DRAFT**: Nháp, chưa submit
- **PENDING**: Đang chờ xử lý
- **APPROVED**: Đã duyệt
- **REJECTED**: Bị từ chối
- **CANCELLED**: Đã hủy
- **COMPLETED**: Hoàn thành
- **ACTIVE**: Đang hoạt động
- **INACTIVE**: Ngừng hoạt động

### Work Order
- **PLANNED**: Đã lên kế hoạch
- **RELEASED**: Đã phát hành
- **IN_PROGRESS**: Đang thực hiện
- **QC_PENDING**: Chờ QC
- **COMPLETED**: Hoàn thành

### Stock Status
- **AVAILABLE**: Có sẵn
- **RESERVED**: Đã giữ chỗ
- **QUARANTINE**: Đang cách ly chờ QC
- **EXPIRED**: Hết hạn
- **DAMAGED**: Hư hỏng

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
