# TÀI LIỆU ĐÀO TẠO HỆ THỐNG ERP

## 1. Lịch trình đào tạo (Training Schedule)

| Buổi | Bộ phận | Thời lượng | Nội dung chính |
|------|---------|------------|----------------|
| 1 | Tất cả | 1h | Tổng quan hệ thống, Đăng nhập, Điều hướng |
| 2 | Kho (Warehouse) | 2h | Quy trình GRN, Tồn kho, Nguyên tắc FEFO |
| 3 | Mua hàng | 1.5h | Nhà cung cấp, PR, PO |
| 4 | Sản xuất (Production) | 2h | BOM, Lệnh sản xuất (Work Order), QC |
| 5 | Quản trị (Admin) | 1h | Phân quyền, User, Audit Log |

---

## 2. Tài liệu Slide - Nội dung đào tạo Kho (WMS)

### Slide 1: Giới thiệu Module Kho
- **Mục đích**: Quản lý dòng chảy nguyên liệu và thành phẩm.
- **Tính năng chính**: Nhập kho từ PO, Xuất kho sản xuất, Kiểm kê tồn kho, Cảnh báo HSD.
- **Slogan**: "FEFO - Hàng cũ ra trước, hàng mới ra sau".

### Slide 2: Quy trình Nhập kho (GRN)
- **Bước 1**: Nhận hàng và đối soát với PO.
- **Bước 2**: Nhập thông tin Lot/Batch và Hạn sử dụng (Bắt buộc).
- **Bước 3**: QC lấy mẫu kiểm tra.
- **Bước 4**: Lưu kho chính thức.

### Slide 3: Nguyên tắc FEFO tự động
- Hệ thống tự động khóa (reserve) các lô hàng có hạn sử dụng gần nhất.
- Nhân viên kho không cần tự chọn Lot, chỉ cần làm theo chỉ dẫn trên màn hình Xuất kho.

---

## 3. Bài tập thực hành (Hands-on Lab)

### Scenario: Nhập kho và Xuất kho Vitamin C

**Task 1: Nhập kho**
1. Đăng nhập tài khoản Nhân viên kho.
2. Tạo phiếu nhập kho (GRN) cho 100kg Vitamin C.
3. Nhập HSD là: 01/01/2026.
4. Hoàn thành phiếu nhập.

**Task 2: Xuất kho**
1. Giả định trong kho đã có sẵn 50kg Vitamin C (HSD: 01/06/2025).
2. Tạo yêu cầu xuất kho 70kg Vitamin C.
3. Quan sát hệ thống:
   - Hệ thống sẽ tự động trừ 50kg từ lô cũ (HSD 2025).
   - Hệ thống sẽ trừ 20kg tiếp theo từ lô bạn vừa nhập (HSD 2026).

---

## 4. Hướng dẫn đánh giá (Assessment)

1. **Câu hỏi**: Tại sao Hạn sử dụng là trường bắt buộc khi nhập kho?
   - *Trả lời*: Để hệ thống có thể vận hành nguyên tắc FEFO tự động.
2. **Câu hỏi**: Nếu QC đánh giá lô hàng "Không đạt", hàng có được nhập vào tồn kho khả dụng không?
   - *Trả lời*: Không, hàng sẽ nằm ở khu vực cách ly (Quarantine) và không thể xuất sản xuất.
