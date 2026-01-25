# HƯỚNG DẪN SỬ DỤNG HỆ THỐNG ERP
## Dành cho công ty mỹ phẩm thiên nhiên

**Phiên bản**: 1.0  
**Cập nhật**: 2026-01-25

---

## MỤC LỤC

1. [Giới thiệu](#1-giới-thiệu)
2. [Đăng nhập & Tài khoản](#2-đăng-nhập--tài-khoản)
3. [Quản lý Nguyên vật liệu](#3-quản-lý-nguyên-vật-liệu)
4. [Quản lý Nhà cung cấp](#4-quản-lý-nhà-cung-cấp)
5. [Mua hàng (PR/PO)](#5-mua-hàng-prpo)
6. [Quản lý Kho (WMS)](#6-quản-lý-kho-wms)
7. [Sản xuất](#7-sản-xuất)
8. [Bán hàng](#8-bán-hàng)
9. [Marketing](#9-marketing)
10. [Báo cáo](#10-báo-cáo)
11. [Quản trị hệ thống](#11-quản-trị-hệ-thống)
12. [Câu hỏi thường gặp](#12-câu-hỏi-thường-gặp)

---

## 1. GIỚI THIỆU

Hệ thống ERP được thiết kế đặc thù cho ngành mỹ phẩm với các tính năng:

- **FEFO** (First Expired First Out): Ưu tiên xuất hàng sắp hết hạn trước
- **Truy xuất nguồn gốc**: Từ thành phẩm → nguyên liệu → nhà cung cấp
- **Bảo mật công thức**: Mã hóa AES-256 cho BOM
- **Chứng nhận GMP**: Theo dõi chứng chỉ nhà cung cấp

---

## 2. ĐĂNG NHẬP & TÀI KHOẢN

### Đăng nhập
1. Truy cập địa chỉ hệ thống
2. Nhập **Email** và **Mật khẩu**
3. Click **"Đăng nhập"**

### Đổi mật khẩu
1. Click avatar góc phải → **Đổi mật khẩu**
2. Nhập mật khẩu cũ và mật khẩu mới
3. Click **Lưu**

> ⚠️ **LƯU Ý**: Tài khoản bị khóa sau **5 lần** nhập sai. Liên hệ Admin để mở khóa.

---

## 3. QUẢN LÝ NGUYÊN VẬT LIỆU

### Thêm nguyên liệu mới
1. **Master Data → Nguyên liệu → Thêm mới**
2. Điền thông tin:
   - **Mã nguyên liệu**: Tự động sinh (RM-0001)
   - **Tên INCI** (International Nomenclature): Tên quốc tế
   - **CAS Number**: Mã đăng ký hóa chất
   - **Chất gây dị ứng**: Đánh dấu nếu có
   - **Điều kiện bảo quản**: Thường / Lạnh (2-8°C) / Đông lạnh
3. Click **Lưu**

### Tra cứu nguyên liệu
- Tìm theo: Tên, Mã, INCI Name, CAS Number
- Lọc theo: Danh mục, Trạng thái, Điều kiện bảo quản

---

## 4. QUẢN LÝ NHÀ CUNG CẤP

### Thêm nhà cung cấp
1. **Mua hàng → Nhà cung cấp → Thêm mới**
2. Điền thông tin công ty, địa chỉ, liên hệ
3. Thêm **Chứng nhận** (GMP, ISO, Organic, HALAL...)
4. Click **Lưu** → Trạng thái: *Chờ duyệt*

### Duyệt nhà cung cấp
1. Kiểm tra chứng nhận GMP còn hiệu lực
2. Click **Duyệt** → Nhà cung cấp vào danh sách ASL

### Cảnh báo chứng nhận hết hạn
- Hệ thống tự động cảnh báo trước **90 ngày / 30 ngày**
- Xem: **Mua hàng → Chứng nhận sắp hết hạn**

> ⚠️ **QUAN TRỌNG**: Không thể tạo PO cho nhà cung cấp có chứng nhận GMP hết hạn.

---

## 5. MUA HÀNG (PR/PO)

### Tạo Yêu cầu mua hàng (PR)
1. **Mua hàng → PR → Tạo mới**
2. Chọn nguyên liệu, số lượng, ngày cần
3. Click **Gửi duyệt**

### Quy trình duyệt
| Giá trị | Người duyệt |
|---------|-------------|
| < 10 triệu | Tự động duyệt |
| 10-50 triệu | Trưởng phòng |
| 50-200 triệu | Trưởng bộ phận Mua hàng |
| > 200 triệu | CFO |

### Chuyển PR → PO
1. PR được duyệt → Click **Tạo PO**
2. Chọn nhà cung cấp, giá, điều khoản
3. Click **Xác nhận PO** → Gửi email cho NCC

---

## 6. QUẢN LÝ KHO (WMS)

> ⚠️ **NGUYÊN TẮC FEFO**: Hệ thống tự động ưu tiên xuất lot **sắp hết hạn trước**.

### Nhập kho (GRN)
1. **Kho → Nhập kho → Tạo mới**
2. Chọn **PO** tương ứng
3. Nhập thông tin từng dòng:
   - Số lượng thực nhận
   - **Lot nhà cung cấp**
   - NSX / **HSD** *(Bắt buộc)*
   - Vị trí lưu kho
4. Click **Lưu** → Trạng thái: *Chờ QC*
5. Sau QC pass → Click **Hoàn thành**

### Xuất kho
1. **Kho → Xuất kho → Tạo mới**
2. Chọn loại xuất: Sản xuất / Bán hàng / Mẫu
3. Chọn nguyên liệu, số lượng
4. Hệ thống **tự động chọn lot** theo FEFO
5. Click **Xác nhận xuất**

### Kiểm tra tồn kho
- **Kho → Tồn kho**: Xem theo vị trí, nguyên liệu
- **Kho → Hàng sắp hết hạn**: Cảnh báo 90/30/7 ngày

---

## 7. SẢN XUẤT

### ⚠️ Bảo mật công thức (BOM)
- Công thức sản phẩm được **mã hóa AES-256**
- Quyền xem theo cấp:
  | Quyền | Xem nguyên liệu | Xem số lượng | Xem công thức |
  |-------|-----------------|--------------|---------------|
  | Staff | ✓ | ✗ | ✗ |
  | Prod. Manager | ✓ | ✓ | ✗ |
  | R&D | ✓ | ✓ | ✓ |
- **Mọi truy cập được ghi log**

### Tạo lệnh sản xuất
1. **Sản xuất → Lệnh SX → Tạo mới**
2. Chọn sản phẩm, số lượng kế hoạch
3. Hệ thống tính định mức nguyên liệu
4. Click **Tạo** → *Đã lên kế hoạch*
5. Click **Phát hành** → *Sẵn sàng SX*
6. Click **Bắt đầu** → Nguyên liệu được giữ kho (FEFO)

### Hoàn thành sản xuất
1. Nhập **số lượng thành phẩm** thực tế
2. Nhập **Lot thành phẩm**, HSD
3. QC kiểm tra
4. Click **Hoàn thành** → Thành phẩm nhập kho

### Truy xuất nguồn gốc
- **Truy xuất ngược**: Lot thành phẩm → Lot nguyên liệu đã dùng
- **Truy xuất xuôi**: Lot nguyên liệu → Thành phẩm nào đã sử dụng

---

## 8. BÁN HÀNG

### Quản lý khách hàng
1. **Bán hàng → Khách hàng → Thêm mới**
2. Điền thông tin, nhóm khách hàng (VIP/Gold/Silver)
3. Đặt **hạn mức tín dụng**

### Tạo đơn hàng
1. **Bán hàng → Đơn hàng → Tạo mới**
2. Chọn khách hàng, sản phẩm, số lượng
3. Hệ thống kiểm tra tồn kho và hạn mức
4. Click **Xác nhận** → Hàng được giữ kho

### Giao hàng
1. Chọn đơn hàng → **Tạo phiếu giao**
2. Chọn đơn vị vận chuyển
3. Xuất kho (FEFO) → In phiếu giao

---

## 9. MARKETING

### Quản lý KOL
1. **Marketing → KOL → Thêm mới**
2. Phân loại: Mega / Macro / Micro / Nano
3. Theo dõi bài đăng, engagement, ROI

### Gửi mẫu sản phẩm
1. **Marketing → Gửi mẫu → Tạo yêu cầu**
2. Chọn KOL, sản phẩm, số lượng
3. Duyệt → Xuất kho (tích hợp WMS)
4. Theo dõi trạng thái giao hàng

---

## 10. BÁO CÁO

### Dashboard
- Tổng quan: Tồn kho, Đơn hàng, Sản xuất
- KPI theo thời gian thực

### Báo cáo có sẵn
| Báo cáo | Mô tả |
|---------|-------|
| Tồn kho tổng hợp | Theo nguyên liệu, kho, vị trí |
| Hàng sắp hết hạn | Theo số ngày (90/30/7) |
| Biến động kho | Nhập/Xuất theo thời gian |
| Hiệu suất NCC | Giao đúng hạn, chất lượng |
| Sản lượng SX | Theo sản phẩm, ca, ngày |
| Doanh số | Theo khách hàng, sản phẩm |

### Xuất báo cáo
- Định dạng: **Excel**, **CSV**
- Click biểu tượng **Tải xuống**

---

## 11. QUẢN TRỊ HỆ THỐNG

### Quản lý người dùng
1. **Admin → Người dùng → Thêm mới**
2. Nhập thông tin, gán vai trò
3. Gán vào phòng ban

### Phân quyền
- Quyền theo format: `service:resource:action`
- Ví dụ: `wms:stock:read`, `manufacturing:bom:formula_view`
- Có thể gán wildcard: `wms:*:*`

### Xem Audit Log
- **Admin → Audit Log**
- Theo dõi: Ai làm gì, lúc nào, dữ liệu thay đổi

---

## 12. CÂU HỎI THƯỜNG GẶP

**Q: Quên mật khẩu?**  
A: Click "Quên mật khẩu" → Nhập email → Check email để reset.

**Q: Tài khoản bị khóa?**  
A: Liên hệ Admin để mở khóa (bị khóa sau 5 lần sai).

**Q: Tại sao không chọn được Lot khi xuất kho?**  
A: Hệ thống tự động chọn theo FEFO. Lot sắp hết hạn được ưu tiên.

**Q: Hàng đã nhập nhưng không thấy tồn kho?**  
A: Kiểm tra GRN đã được QC duyệt và "Hoàn thành" chưa.

**Q: Không tạo được lệnh sản xuất?**  
A: Kiểm tra sản phẩm có BOM đã được duyệt chưa.

**Q: Không xem được công thức?**  
A: Bạn không có quyền `formula_view`. Liên hệ R&D Manager.

**Q: Không tạo được PO cho nhà cung cấp?**  
A: Kiểm tra nhà cung cấp đã được duyệt và có chứng nhận GMP còn hiệu lực.

---

**Hỗ trợ kỹ thuật**: it-support@company.vn  
**Hotline**: 1900-xxxx
