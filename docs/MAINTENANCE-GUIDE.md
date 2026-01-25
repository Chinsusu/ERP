# HƯỚNG DẪN BẢO TRÌ HỆ THỐNG (MAINTENANCE GUIDE)

## 1. Kiểm tra hàng ngày (Daily Operations)

### Buổi sáng (8:00 AM)
- [ ] Kiểm tra trạng thái các service (Grafana: Service Health).
- [ ] Kiểm tra các cảnh báo (Alerts) trong đêm qua.
- [ ] Kiểm tra thông lượng (Request rate) và tỷ lệ lỗi (Error rate < 0.1%).
- [ ] Kiểm tra dung lượng đĩa cứng (Disk usage < 80%).

### Buổi chiều (5:00 PM)
- [ ] Xem xét các sự cố trong ngày và cập nhật tài liệu Incident Report.
- [ ] Xác nhận bản backup cơ sở dữ liệu gần nhất đã được tạo thành công.

---

## 2. Bảo trì định kỳ

### Hàng tuần (Sáng thứ Hai)
- **Review Hiệu năng**: Kiểm tra các query chậm (slow queries) và tăng trưởng dữ liệu.
- **Bảo mật**: Xem các lần login thất bại liên tiếp (brute force attempts).
- **Dọn dẹp**: Xóa logs cũ hơn 30 ngày (`scripts/cleanup-backups.sh`).

### Hàng tháng (Thứ Bảy đầu tiên)
- **Kiểm tra Backup**: Thực hiện diễn tập khôi phục (Restore Test) trên môi trường Staging/Test.
- **Cập nhật**: Áp dụng các bản vá bảo mật (Security patches) và cập nhật Docker images.
- **Kế hoạch năng lực**: Dự báo mức tăng trưởng dữ liệu để nâng cấp tài nguyên (CPU/RAM/Disk).

---

## 3. Quy trình Khôi phục dữ liệu (Recovery)

1. **Chuẩn bị môi trường**: Khởi tạo container Postgres mới.
2. **Khôi phục**:
   ```bash
   gunzip -c /backups/db/erp_prod_LATEST.sql.gz | docker exec -i new-postgres psql -U user -d erp_prod
   ```
3. **Xác minh**: Kiểm tra số lượng Record giữa bản backup và bản thực tế.

---

## 4. Quản lý sự cố (Incident Management)

- **P1 (Nghiêm trọng)**: Hệ thống ngưng hoạt động toàn bộ. Phản hồi trong 15 phút.
- **P2 (Lớn)**: Một chức năng chính (VD: Nhập kho) không khả dụng. Phản hồi trong 30 phút.
- **P3 (Nhỏ)**: Lỗi không ảnh hưởng nghiêm trọng đến vận hành. Phản hồi trong 4 giờ.
