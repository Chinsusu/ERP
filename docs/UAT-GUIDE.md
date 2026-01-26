# 📋 HƯỚNG DẪN UAT (User Acceptance Testing)
## Hệ thống ERP Mỹ phẩm Thiên nhiên

**Phiên bản**: 1.0  
**Ngày**: 2026-01-26  
**Thời gian UAT**: 2 tuần

---

## 📑 MỤC LỤC

1. [Tổng quan UAT](#1-tổng-quan-uat)
2. [Chuẩn bị UAT](#2-chuẩn-bị-uat)
3. [Kịch bản test theo phòng ban](#3-kịch-bản-test-theo-phòng-ban)
4. [Mẫu biên bản & báo lỗi](#4-mẫu-biên-bản--báo-lỗi)
5. [Tiêu chí Go-Live](#5-tiêu-chí-go-live)

---

## 1. TỔNG QUAN UAT

### 1.1 Mục đích UAT
- Xác nhận hệ thống hoạt động đúng nghiệp vụ thực tế
- Phát hiện lỗi trước khi go-live
- Đào tạo người dùng làm quen hệ thống
- Thu thập feedback cải thiện

### 1.2 Phạm vi test

| Module | Phòng ban | Độ ưu tiên |
|--------|-----------|------------|
| Quản lý Kho (WMS) | Kho | ⭐⭐⭐ Critical |
| Sản xuất (BOM, WO) | Sản xuất | ⭐⭐⭐ Critical |
| Mua hàng (PR, PO) | Mua hàng | ⭐⭐ High |
| Nhà cung cấp | Mua hàng | ⭐⭐ High |
| Bán hàng | Kinh doanh | ⭐⭐ High |
| Marketing | Marketing | ⭐ Medium |
| Báo cáo | Tất cả | ⭐⭐ High |

### 1.3 Timeline 2 tuần

```
TUẦN 1:
├── Ngày 1-2: Đào tạo + Setup tài khoản
├── Ngày 3-4: UAT Module Kho + Nhập xuất
└── Ngày 5: UAT Module Sản xuất (BOM, WO)

TUẦN 2:
├── Ngày 1-2: UAT Mua hàng + Nhà cung cấp
├── Ngày 3: UAT Bán hàng + Marketing
├── Ngày 4: UAT Báo cáo + Tổng hợp
└── Ngày 5: Fix bugs + Nghiệm thu
```

---

## 2. CHUẨN BỊ UAT

### 2.1 Môi trường
```
URL: https://uat.erp.company.vn
Database: Dữ liệu test (không phải production)
```

### 2.2 Tài khoản UAT

| Vai trò | Email | Password | Phòng ban |
|---------|-------|----------|-----------|
| Admin | admin@company.vn | Uat@2026! | IT |
| Kho trưởng | kho.tp@company.vn | Uat@2026! | Kho |
| NV Kho | kho.nv@company.vn | Uat@2026! | Kho |
| TP Sản xuất | sx.tp@company.vn | Uat@2026! | Sản xuất |
| NV Sản xuất | sx.nv@company.vn | Uat@2026! | Sản xuất |
| TP Mua hàng | mh.tp@company.vn | Uat@2026! | Mua hàng |
| NV Mua hàng | mh.nv@company.vn | Uat@2026! | Mua hàng |
| TP Kinh doanh | kd.tp@company.vn | Uat@2026! | Sales |
| Marketing | mkt@company.vn | Uat@2026! | Marketing |

### 2.3 Dữ liệu mẫu cần có

**Nguyên vật liệu (20+ items):**
- Tinh dầu Tràm trà, Tinh dầu Bưởi, Tinh dầu Oải hương
- Dầu dừa, Bơ hạt mỡ, Vitamin E
- Chiết xuất Nha đam, Chiết xuất Trà xanh
- Chai 30ml, Chai 50ml, Hũ 50g
- Nhãn, Hộp giấy, Túi giấy

**Nhà cung cấp (5+):**
- NCC001: Công ty Tinh dầu Việt
- NCC002: Công ty Dầu dừa Bến Tre
- NCC003: Công ty Bao bì Xanh
- ...

**Sản phẩm (10+):**
- Serum Vitamin C, Kem dưỡng ban đêm
- Sữa rửa mặt Trà xanh, Son dưỡng môi
- ...

---

## 3. KỊCH BẢN TEST THEO PHÒNG BAN

---

### 📦 3.1 MODULE KHO (WMS)
**Người test**: Phòng Kho  
**Thời gian**: 2 ngày

#### TC-KHO-001: Nhập kho theo PO
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào menu **Kho → Nhập kho** | Hiển thị danh sách GRN | |
| 2 | Click **"Tạo phiếu nhập"** | Mở form tạo mới | |
| 3 | Chọn PO cần nhập hàng | Hiển thị danh sách NVL từ PO | |
| 4 | Với mỗi NVL, nhập: | | |
| | - Số lượng thực nhận: **100** | OK | |
| | - Số Lot: **LOT-2026-001** | OK | |
| | - Ngày sản xuất: **01/01/2026** | OK | |
| | - Hạn sử dụng: **01/01/2028** | OK | |
| | - Vị trí kho: **A-01-01** | OK | |
| 5 | Click **"Hoàn thành nhập"** | Thông báo thành công | |
| 6 | Kiểm tra tồn kho | Số lượng tăng đúng | |

**Ghi chú**: _____________________

---

#### TC-KHO-002: ⭐ Xuất kho theo FEFO (QUAN TRỌNG)
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | **Chuẩn bị**: Có 2 Lot cùng 1 NVL: | | |
| | - Lot A: HSD **30/06/2026**, SL: 50 | | |
| | - Lot B: HSD **31/03/2026**, SL: 100 | | |
| 2 | Vào **Kho → Xuất kho** | | |
| 3 | Tạo phiếu xuất, chọn NVL, SL: **80** | | |
| 4 | Kiểm tra hệ thống gợi ý Lot | **Phải gợi ý Lot B trước** (HSD gần hơn) | |
| 5 | Xác nhận xuất | Tồn Lot B: 20, Lot A: 50 | |

**⚠️ Nếu hệ thống KHÔNG chọn Lot B trước → BÁO LỖI CRITICAL**

**Ghi chú**: _____________________

---

#### TC-KHO-003: Cảnh báo hàng sắp hết hạn
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Có NVL với HSD còn **25 ngày** | | |
| 2 | Vào **Dashboard** | Hiển thị cảnh báo màu đỏ | |
| 3 | Click vào cảnh báo | Hiển thị danh sách NVL sắp hết hạn | |

**Ghi chú**: _____________________

---

#### TC-KHO-004: Kiểm kê kho
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Kho → Kiểm kê** | | |
| 2 | Tạo phiên kiểm kê mới | Hiển thị form | |
| 3 | Chọn kho/vị trí cần kiểm | Hiển thị tồn hệ thống | |
| 4 | Nhập số lượng thực tế (khác hệ thống) | Hiển thị chênh lệch | |
| 5 | Lưu và gửi duyệt | Tạo phiếu điều chỉnh | |

**Ghi chú**: _____________________

---

### 🏭 3.2 MODULE SẢN XUẤT
**Người test**: Phòng Sản xuất  
**Thời gian**: 1 ngày

#### TC-SX-001: Tạo công thức BOM
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Sản xuất → Công thức (BOM)** | | |
| 2 | Click **"Tạo BOM mới"** | | |
| 3 | Nhập thông tin: | | |
| | - Mã BOM: **BOM-KEM-001** | | |
| | - Sản phẩm: **Kem dưỡng ban đêm** | | |
| | - SL đầu ra: **100 hũ** | | |
| 4 | Thêm nguyên liệu: | | |
| | - Dầu dừa: **5 kg** | | |
| | - Bơ hạt mỡ: **2 kg** | | |
| | - Vitamin E: **0.5 kg** | | |
| | - Hũ 50g: **100 cái** | | |
| 5 | Click **"Lưu nháp"** | Status = Draft | |
| 6 | Click **"Gửi duyệt"** | Status = Pending | |

**Ghi chú**: _____________________

---

#### TC-SX-002: ⭐ Kiểm tra mã hóa BOM (QUAN TRỌNG)
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Đăng nhập bằng TK **Nhân viên SX** | | |
| 2 | Mở BOM đã duyệt | | |
| 3 | Kiểm tra chi tiết công thức | **Chỉ thấy tên NVL, KHÔNG thấy tỷ lệ %** | |
| 4 | Đăng nhập TK **Trưởng SX** hoặc **Admin** | | |
| 5 | Mở cùng BOM đó | **Thấy đầy đủ tỷ lệ %** | |

**⚠️ Nếu NV thấy được tỷ lệ % → BÁO LỖI CRITICAL**

**Ghi chú**: _____________________

---

#### TC-SX-003: Tạo lệnh sản xuất (Work Order)
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Sản xuất → Lệnh sản xuất** | | |
| 2 | Click **"Tạo WO mới"** | | |
| 3 | Chọn BOM: **BOM-KEM-001** | Load danh sách NVL | |
| 4 | Nhập SL sản xuất: **200 hũ** (x2 BOM) | Tự tính NVL cần dùng | |
| 5 | Kiểm tra NVL cần: | | |
| | - Dầu dừa: **10 kg** (5x2) | Đúng | |
| | - Bơ hạt mỡ: **4 kg** (2x2) | Đúng | |
| 6 | Kiểm tra tồn kho | Hiển thị đủ/thiếu | |
| 7 | Click **"Tạo lệnh"** | WO Status = Planned | |

**Ghi chú**: _____________________

---

#### TC-SX-004: Hoàn thành sản xuất + QC
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Mở WO đã tạo, click **"Bắt đầu SX"** | Status = In Progress | |
| 2 | Click **"Xuất NVL"** | Tự động xuất kho theo FEFO | |
| 3 | Sau khi SX xong, click **"Nhập thành phẩm"** | | |
| 4 | Nhập: | | |
| | - SL thực tế: **198** (hao hụt 2) | | |
| | - Số Lot: **LOT-KEM-2026-001** | | |
| | - NSX: **01/02/2026** | | |
| | - HSD: **01/02/2028** | | |
| 5 | Thực hiện QC | Form kiểm tra chất lượng | |
| 6 | Nhập kết quả QC → **Pass** | | |
| 7 | Click **"Hoàn thành"** | WO = Completed, TP nhập kho | |

**Ghi chú**: _____________________

---

#### TC-SX-005: ⭐ Truy xuất nguồn gốc (QUAN TRỌNG)
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Sản xuất → Truy xuất** | | |
| 2 | Nhập Lot thành phẩm: **LOT-KEM-2026-001** | | |
| 3 | Click **"Truy xuất ngược"** | | |
| 4 | Kiểm tra kết quả: | | |
| | - Hiển thị danh sách Lot NVL đã dùng | ✅ | |
| | - Hiển thị NCC cung cấp từng Lot | ✅ | |
| | - Hiển thị ngày nhập, ngày SX | ✅ | |
| 5 | Click **"Truy xuất xuôi"** từ 1 Lot NVL | Hiển thị các Lot TP đã dùng NVL này | |

**⚠️ Nếu không truy xuất được → BÁO LỖI CRITICAL**

**Ghi chú**: _____________________

---

### 🛒 3.3 MODULE MUA HÀNG
**Người test**: Phòng Mua hàng  
**Thời gian**: 1.5 ngày

#### TC-MH-001: Tạo yêu cầu mua hàng (PR)
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Mua hàng → Yêu cầu mua hàng** | | |
| 2 | Click **"Tạo PR mới"** | | |
| 3 | Nhập: | | |
| | - Mô tả: "Mua NVL tháng 2/2026" | | |
| | - Ngày cần: **15/02/2026** | | |
| 4 | Thêm NVL: | | |
| | - Dầu dừa: **100 kg** | | |
| | - Vitamin E: **10 kg** | | |
| 5 | Click **"Lưu nháp"** | PR = Draft | |
| 6 | Click **"Gửi duyệt"** | PR = Pending Approval | |

**Ghi chú**: _____________________

---

#### TC-MH-002: Duyệt PR (Trưởng phòng)
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Đăng nhập TK **Trưởng phòng MH** | | |
| 2 | Vào **Mua hàng → PR chờ duyệt** | Thấy PR vừa tạo | |
| 3 | Mở chi tiết, kiểm tra | Đầy đủ thông tin | |
| 4 | Click **"Duyệt"** | PR = Approved | |
| 5 | Kiểm tra lịch sử duyệt | Có tên người duyệt + thời gian | |

**Ghi chú**: _____________________

---

#### TC-MH-003: Tạo PO từ PR
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Từ PR đã duyệt, click **"Tạo PO"** | Mở form PO, copy data từ PR | |
| 2 | Chọn NCC: **NCC001** | | |
| 3 | Nhập đơn giá: | | |
| | - Dầu dừa: **80,000 VND/kg** | | |
| | - Vitamin E: **500,000 VND/kg** | | |
| 4 | Kiểm tra tổng tiền | = 8,000,000 + 5,000,000 = **13,000,000** | |
| 5 | Nhập điều khoản thanh toán | | |
| 6 | Click **"Lưu & Gửi NCC"** | PO = Sent | |
| 7 | In PO | PDF đúng format, có logo | |

**Ghi chú**: _____________________

---

### 🏪 3.4 MODULE NHÀ CUNG CẤP
**Người test**: Phòng Mua hàng  
**Thời gian**: 0.5 ngày

#### TC-NCC-001: Thêm nhà cung cấp mới
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Mua hàng → Nhà cung cấp** | | |
| 2 | Click **"Thêm NCC"** | | |
| 3 | Nhập thông tin: | | |
| | - Mã: **NCC006** | | |
| | - Tên: "Công ty ABC" | | |
| | - MST: 0123456789 | | |
| | - Địa chỉ, SĐT, Email | | |
| 4 | Click **"Lưu"** | NCC được tạo | |

**Ghi chú**: _____________________

---

#### TC-NCC-002: Upload chứng chỉ GMP
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Mở chi tiết NCC | | |
| 2 | Tab **"Chứng chỉ"** → **"Thêm"** | | |
| 3 | Chọn loại: **GMP** | | |
| 4 | Nhập số chứng chỉ, ngày cấp, ngày hết hạn | | |
| 5 | Upload file PDF | | |
| 6 | Click **"Lưu"** | Chứng chỉ được lưu | |

**Ghi chú**: _____________________

---

#### TC-NCC-003: Cảnh báo chứng chỉ sắp hết hạn
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Tạo chứng chỉ với HSD còn **20 ngày** | | |
| 2 | Vào **Dashboard** | Hiển thị cảnh báo | |
| 3 | Click vào cảnh báo | Danh sách NCC có cert sắp hết hạn | |

**Ghi chú**: _____________________

---

### 💰 3.5 MODULE BÁN HÀNG
**Người test**: Phòng Kinh doanh  
**Thời gian**: 0.5 ngày

#### TC-BH-001: Tạo báo giá → Đơn hàng → Xuất hàng
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Bán hàng → Khách hàng** | | |
| 2 | Tạo khách hàng mới | KH được tạo | |
| 3 | Vào **Bán hàng → Báo giá** | | |
| 4 | Tạo báo giá cho KH: | | |
| | - Kem dưỡng: 50 hũ x 350,000 | | |
| | - Serum: 30 chai x 450,000 | | |
| 5 | Tính tổng | = 17,500,000 + 13,500,000 = **31,000,000** | |
| 6 | Gửi báo giá | Status = Sent | |
| 7 | KH đồng ý → **"Tạo đơn hàng"** | SO được tạo từ Quotation | |
| 8 | Xác nhận đơn hàng | SO = Confirmed | |
| 9 | Click **"Xuất hàng"** | Tạo phiếu xuất kho | |
| 10 | Kiểm tra Lot được chọn | **FEFO - Lot gần HSD trước** | |
| 11 | Xác nhận xuất | Trừ tồn kho | |

**Ghi chú**: _____________________

---

### 📊 3.6 MODULE BÁO CÁO
**Người test**: Tất cả phòng ban  
**Thời gian**: 0.5 ngày

#### TC-BC-001: Báo cáo tồn kho
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Vào **Báo cáo → Tồn kho** | | |
| 2 | Chọn kho, ngày | | |
| 3 | Click **"Xem báo cáo"** | Hiển thị tồn kho chi tiết | |
| 4 | Kiểm tra số liệu | Khớp với thực tế đã test | |
| 5 | Click **"Export Excel"** | File .xlsx download | |

**Ghi chú**: _____________________

---

#### TC-BC-002: Dashboard tổng quan
| Bước | Thao tác | Kết quả mong đợi | Pass/Fail |
|------|----------|------------------|-----------|
| 1 | Đăng nhập Admin | | |
| 2 | Xem Dashboard | Hiển thị các widget | |
| 3 | Kiểm tra widget Tồn kho | Số liệu đúng | |
| 4 | Kiểm tra widget Sản xuất | Số liệu đúng | |
| 5 | Kiểm tra cảnh báo | Hiển thị đúng các cảnh báo | |

**Ghi chú**: _____________________

---

## 4. MẪU BIÊN BẢN & BÁO LỖI

### 4.1 Bảng tổng hợp kết quả

| Module | Tổng TC | Pass | Fail | Tỷ lệ | Người test | Ngày |
|--------|---------|------|------|-------|------------|------|
| Kho (WMS) | | | | % | | |
| Sản xuất | | | | % | | |
| Mua hàng | | | | % | | |
| Nhà cung cấp | | | | % | | |
| Bán hàng | | | | % | | |
| Báo cáo | | | | % | | |
| **TỔNG** | | | | **%** | | |

### 4.2 Phân loại lỗi

| Mức độ | Mô tả | Thời gian fix | Ví dụ |
|--------|-------|---------------|-------|
| 🔴 **Critical** | Không thể sử dụng | 24h | FEFO sai, BOM lộ công thức |
| 🟠 **High** | Chức năng chính lỗi | 48h | Không tạo được WO |
| 🟡 **Medium** | Lỗi có workaround | 1 tuần | Filter sai, sort sai |
| 🟢 **Low** | Lỗi UI/UX | Sau go-live | Typo, căn lề |

### 4.3 Form báo lỗi

```
╔══════════════════════════════════════════════════════════╗
║                     BÁO CÁO LỖI                          ║
╠══════════════════════════════════════════════════════════╣
║ Mã lỗi: BUG-____        Ngày: ___/___/2026               ║
║ Người báo: _______________  Phòng: _______________       ║
║                                                          ║
║ Module: [ ] Kho  [ ] SX  [ ] MH  [ ] BH  [ ] Khác        ║
║ Mức độ: [ ] Critical  [ ] High  [ ] Medium  [ ] Low      ║
║                                                          ║
║ MÔ TẢ LỖI:                                               ║
║ ________________________________________________________ ║
║ ________________________________________________________ ║
║                                                          ║
║ CÁC BƯỚC TÁI HIỆN:                                       ║
║ 1. ____________________________________________________  ║
║ 2. ____________________________________________________  ║
║ 3. ____________________________________________________  ║
║                                                          ║
║ KẾT QUẢ MONG ĐỢI: _____________________________________  ║
║ KẾT QUẢ THỰC TẾ:  _____________________________________  ║
║                                                          ║
║ Screenshot: [ ] Có  [ ] Không                            ║
╚══════════════════════════════════════════════════════════╝
```

### 4.4 Biên bản nghiệm thu

```
╔══════════════════════════════════════════════════════════╗
║           BIÊN BẢN NGHIỆM THU UAT                        ║
║          HỆ THỐNG ERP MỸ PHẨM                            ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║ Ngày: ___/___/2026                                       ║
║                                                          ║
║ KẾT QUẢ:                                                 ║
║ - Tổng test case: _______                                ║
║ - Pass: _______ ( _____ %)                               ║
║ - Fail: _______ ( _____ %)                               ║
║                                                          ║
║ LỖI CÒN LẠI:                                             ║
║ - Critical: _______                                       ║
║ - High: _______                                          ║
║ - Medium: _______                                        ║
║ - Low: _______                                           ║
║                                                          ║
║ KẾT LUẬN:                                                ║
║ [ ] ĐẠT - Cho phép Go-Live                               ║
║ [ ] ĐẠT CÓ ĐIỀU KIỆN - Fix bugs trước Go-Live            ║
║ [ ] KHÔNG ĐẠT - Test lại                                 ║
║                                                          ║
║ CHỮ KÝ:                                                  ║
║                                                          ║
║ Đại diện Công ty:         Đại diện IT:                   ║
║                                                          ║
║ ___________________       ___________________            ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
```

---

## 5. TIÊU CHÍ GO-LIVE

### 5.1 Điều kiện BẮT BUỘC

| # | Tiêu chí | Yêu cầu | Đạt |
|---|----------|---------|-----|
| 1 | Lỗi Critical | = 0 | [ ] |
| 2 | Lỗi High | = 0 | [ ] |
| 3 | Tỷ lệ Pass | ≥ 95% | [ ] |
| 4 | **FEFO** hoạt động đúng | ✅ | [ ] |
| 5 | **BOM mã hóa** hoạt động | ✅ | [ ] |
| 6 | **Truy xuất nguồn gốc** hoạt động | ✅ | [ ] |
| 7 | Phân quyền đúng | ✅ | [ ] |
| 8 | Báo cáo số liệu chính xác | ✅ | [ ] |

### 5.2 Checklist Go-Live

```
UAT:
[ ] Biên bản nghiệm thu đã ký
[ ] Tất cả lỗi Critical/High đã fix
[ ] Re-test pass

DATA:
[ ] Master data đã nhập (NVL, SP, NCC, KH)
[ ] User accounts đã tạo
[ ] Tồn kho đầu kỳ đã nhập

TRAINING:
[ ] Tài liệu đã phát cho user
[ ] Training đã hoàn thành
[ ] Hotline support đã thông báo

APPROVAL:
[ ] IT Lead ký duyệt
[ ] Trưởng phòng nghiệp vụ ký
[ ] Ban Giám đốc phê duyệt
```

---

## LIÊN HỆ HỖ TRỢ

| Vai trò | Tên | SĐT | Email |
|---------|-----|-----|-------|
| IT Support | | | |
| Project Manager | | | |

**Hotline UAT**: 0xxx-xxx-xxx (8:00-17:00)

---

*Tài liệu chuẩn bị: 2026-01-26*
