# TÀI LIỆU ĐÀO TẠO HỆ THỐNG ERP

**Phiên bản**: 1.0  
**Cập nhật**: 2026-01-25

---

## 1. LỊCH TRÌNH ĐÀO TẠO

| Buổi | Bộ phận | Thời lượng | Nội dung chính |
|------|---------|------------|----------------|
| 1 | Tất cả | 1h | Tổng quan, Đăng nhập, Điều hướng |
| 2 | Kho | 2h | GRN, Xuất kho, FEFO, Kiểm kê |
| 3 | Mua hàng | 1.5h | Nhà cung cấp, PR, PO |
| 4 | Sản xuất | 2h | BOM, Lệnh SX, QC, Truy xuất |
| 5 | Bán hàng | 1.5h | Khách hàng, Báo giá, Đơn hàng |
| 6 | Admin | 1h | Users, Roles, Permissions |

---

## 2. NỘI DUNG ĐÀO TẠO

### 2.1 Buổi 1: Tổng quan (Tất cả)

**Mục tiêu**: Hiểu hệ thống và thao tác cơ bản

**Nội dung**:
- Giới thiệu hệ thống ERP mỹ phẩm
- Đăng nhập, đổi mật khẩu
- Điều hướng menu
- Dark mode, tùy chỉnh giao diện

---

### 2.2 Buổi 2: Quản lý Kho (WMS)

**Mục tiêu**: Thành thạo nhập/xuất kho, hiểu FEFO

#### Slide 1: Tổng quan Kho
- Cấu trúc: Kho → Zone → Vị trí
- Khái niệm Lot/Batch
- **FEFO vs FIFO**: Ngành mỹ phẩm dùng FEFO

#### Slide 2: Nguyên tắc FEFO
> **"First Expired First Out"** - Hàng sắp hết hạn xuất trước

- Hệ thống **tự động chọn** lot có HSD gần nhất
- Nhân viên không cần chọn lot thủ công
- Giảm thiểu hàng hết hạn trong kho

#### Slide 3: Quy trình Nhập kho (GRN)
1. Nhận hàng từ nhà cung cấp
2. Đối chiếu với PO
3. Nhập thông tin: Lot, NSX, **HSD** (bắt buộc)
4. QC lấy mẫu kiểm tra
5. Hoàn thành → Hàng vào tồn kho

#### Slide 4: Quy trình Xuất kho
1. Tạo yêu cầu xuất (Sản xuất/Bán hàng/Mẫu)
2. Hệ thống tự động chọn lot FEFO
3. Xác nhận xuất
4. In phiếu xuất kho

---

### 2.3 Buổi 3: Mua hàng

**Mục tiêu**: Tạo và theo dõi PR/PO

#### Slide 1: Quản lý Nhà cung cấp
- Thêm nhà cung cấp mới
- Quản lý chứng nhận (GMP, ISO, Organic)
- Cảnh báo chứng nhận hết hạn
- Danh sách ASL (Approved Supplier List)

#### Slide 2: Quy trình PR → PO
```
Tạo PR → Duyệt → Tạo PO → Xác nhận → Nhận hàng
```

#### Slide 3: Các cấp duyệt
| Giá trị | Người duyệt |
|---------|-------------|
| < 10M | Tự động |
| 10-50M | Trưởng phòng |
| 50-200M | TP Mua hàng |
| > 200M | CFO |

---

### 2.4 Buổi 4: Sản xuất

**Mục tiêu**: Tạo lệnh SX, hiểu bảo mật BOM

#### Slide 1: Bảo mật công thức
> ⚠️ **BOM là TÀI SẢN MẬT** - Không chụp, in, sao chép

- Mã hóa AES-256
- Phân quyền 3 cấp: Xem NL / Xem SL / Xem công thức
- Mọi truy cập được ghi log

#### Slide 2: Quy trình Lệnh sản xuất
```
Tạo → Phát hành → Bắt đầu → Xuất NL → Hoàn thành
```

#### Slide 3: Truy xuất nguồn gốc
- **Ngược**: Thành phẩm → Nguyên liệu → Nhà cung cấp
- **Xuôi**: Nguyên liệu → Thành phẩm nào đã dùng

---

### 2.5 Buổi 5: Bán hàng

**Mục tiêu**: Quản lý khách hàng và đơn hàng

#### Slide 1: Quản lý khách hàng
- Nhóm khách: VIP / Gold / Silver
- Hạn mức tín dụng
- Lịch sử mua hàng

#### Slide 2: Quy trình bán hàng
```
Báo giá → Đơn hàng → Xác nhận → Giao hàng
```

---

### 2.6 Buổi 6: Quản trị (Admin)

**Mục tiêu**: Quản lý users, roles, permissions

#### Slide 1: Quản lý người dùng
- Tạo user, gán role
- Reset mật khẩu
- Mở khóa tài khoản

#### Slide 2: Phân quyền RBAC
- Format: `service:resource:action`
- Ví dụ: `wms:stock:read`
- Wildcard: `wms:*:*`

---

## 3. BÀI TẬP THỰC HÀNH

### Lab 1: Nhập xuất kho FEFO

**Scenario**: Nhập và xuất Vitamin C

**Task 1: Nhập kho**
1. Đăng nhập tài khoản Nhân viên kho
2. Tạo GRN cho **100kg Vitamin C**, HSD: **01/01/2027**
3. Hoàn thành phiếu nhập

**Task 2: Xuất kho**
1. Giả sử trong kho có sẵn **50kg** (HSD: 01/06/2026)
2. Tạo yêu cầu xuất **70kg**
3. Quan sát:
   - Hệ thống trừ **50kg** từ lot cũ (HSD 2026)
   - Hệ thống trừ **20kg** từ lot mới (HSD 2027)

**Kết quả mong đợi**: Lot sắp hết hạn được xuất trước

---

### Lab 2: Tạo lệnh sản xuất

**Scenario**: Sản xuất 1000 lọ Serum

1. Tạo lệnh SX cho sản phẩm "Serum HA"
2. Kiểm tra định mức nguyên liệu
3. Phát hành lệnh
4. Bắt đầu → Quan sát nguyên liệu được giữ kho
5. (Mô phỏng) Hoàn thành sản xuất

---

### Lab 3: Mua hàng

**Scenario**: Mua Hyaluronic Acid

1. Tạo PR cho 50kg Hyaluronic Acid
2. Gửi duyệt
3. Sau khi duyệt → Tạo PO
4. Chọn nhà cung cấp (kiểm tra GMP)
5. Xác nhận PO

---

## 4. CÂU HỎI ĐÁNH GIÁ

1. **Tại sao Hạn sử dụng là trường bắt buộc khi nhập kho?**
   - *TL*: Để hệ thống vận hành FEFO tự động

2. **Nếu QC đánh giá "Không đạt", hàng có vào tồn kho?**
   - *TL*: Không, hàng nằm ở Quarantine

3. **Ai được xem công thức chi tiết trong BOM?**
   - *TL*: Chỉ người có quyền `formula_view` (thường là R&D)

4. **Tại sao không tạo được PO cho một nhà cung cấp?**
   - *TL*: NCC chưa duyệt hoặc chứng nhận GMP hết hạn

---

## 5. TÀI LIỆU THAM KHẢO

- [Hướng dẫn sử dụng chi tiết](file:///opt/ERP/docs/USER-MANUAL-VN.md)
- [API Documentation](file:///opt/ERP/docs/OPENAPI-SPEC.yaml)
- [Deployment Guide](file:///opt/ERP/docs/16-DEPLOYMENT.md)

---

**Liên hệ đào tạo**: training@company.vn
