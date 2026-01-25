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

---

## 2. ĐĂNG NHẬP & TÀI KHOẢN

### Giới thiệu
Hệ thống sử dụng xác thực tập trung để đảm bảo an toàn thông tin và phân quyền chính xác cho từng nhân viên.

### Các thao tác cơ bản
1. Truy cập vào địa chỉ hệ thống.
2. Nhập Email và Mật khẩu.
3. Click "Đăng nhập".

⚠️ **LƯU Ý**: Tài khoản sẽ bị khóa tự động sau 5 lần nhập sai mật khẩu liên tiếp để bảo vệ khỏi các cuộc tấn công brute-force.

---

## 6. QUẢN LÝ KHO (WMS)

### Giới thiệu
Module WMS quản lý toàn bộ hoạt động kho, từ lúc nhập hàng cho đến khi xuất kho sản xuất hoặc bán hàng.

⚠️ **QUAN TRỌNG**: Hệ thống áp dụng nguyên tắc **FEFO (First Expired First Out)**. Hàng sắp hết hạn sẽ được xuất trước một cách tự động.

### Nhập kho (GRN)
1. Vào menu: **Kho > Nhập kho > Tạo mới**.
2. Chọn PO (Purchase Order) tương ứng.
3. Hệ thống sẽ tự động tải danh sách nguyên liệu từ PO.
4. Nhập thông tin chi tiết:
   - Số lượng thực nhận.
   - Số Lot nhà cung cấp.
   - Ngày sản xuất & **Hạn sử dụng (Bắt buộc)**.
   - Vị trí lưu kho.
5. Click "Lưu". Trạng thái sẽ là "Chờ QC".
6. Sau khi QC duyệt, click "Hoàn thành" để nhập hàng vào kho chính.

---

## 7. SẢN XUẤT

### Bảo mật công thức (BOM)
⚠️ **CẢNH BÁO BẢO MẬT**:
- Công thức sản phẩm là tài sản mật của công ty.
- Hệ thống mã hóa thông tin công thức bằng chuẩn AES-256.
- Chỉ người có quyền `formula_view` mới có thể xem chi tiết công thức.

### Lệnh sản xuất (Work Order)
1. Vào menu: **Sản xuất > Lệnh sản xuất > Tạo mới**.
2. Chọn sản phẩm.
3. Nhập số lượng kế hoạch.
4. Hệ thống tự động tính toán định mức nguyên liệu cần thiết dựa trên BOM.
5. Click "Tạo" → Trạng thái: "Đã lên kế hoạch".
6. Click "Bắt đầu" để hệ thống thực hiện giữ kho (reservation) nguyên liệu theo nguyên tắc FEFO.

### Truy xuất nguồn gốc
Hệ thống cho phép truy xuất 2 chiều:
- **Truy xuất ngược**: Từ Lot thành phẩm tra cứu ra toàn bộ Lot nguyên liệu đã sử dụng.
- **Truy xuất xuôi**: Từ Lot nguyên liệu tra cứu ra các lô thành phẩm nào đã sử dụng nó.

---

## 12. CÂU HỎI THƯỜNG GẶP (FAQ)

**Q: Tại sao tôi không thấy nút "Xem công thức" trong BOM?**
A: Bạn không có quyền truy cập thông tin nhạy cảm này. Vui lòng liên hệ quản lý bộ phận R&D.

**Q: Làm sao để mở khóa tài khoản?**
A: Vui lòng liên hệ bộ phận IT hoặc Admin hệ thống để thực hiện Reset Failed Attempts.

**Q: Tại sao hệ thống lại chọn Lot khác với Lot tôi muốn xuất?**
A: Hệ thống luôn ưu tiên Lot có hạn sử dụng gần nhất (FEFO) để giảm thiểu rủi ro hàng hết hạn trong kho.
