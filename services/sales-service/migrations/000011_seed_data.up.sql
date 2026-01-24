-- Seed data for Sales Service

-- Customer Groups
INSERT INTO customer_groups (id, code, name, description, discount_percent, is_active) VALUES
    ('a0000001-0001-0001-0001-000000000001', 'VIP', 'VIP Customers', 'Premium customers with highest discounts', 15.00, true),
    ('a0000001-0001-0001-0001-000000000002', 'GOLD', 'Gold Customers', 'Loyal customers with good discounts', 10.00, true),
    ('a0000001-0001-0001-0001-000000000003', 'SILVER', 'Silver Customers', 'Regular customers with standard discounts', 5.00, true),
    ('a0000001-0001-0001-0001-000000000004', 'REGULAR', 'Regular Customers', 'New and regular customers', 0.00, true),
    ('a0000001-0001-0001-0001-000000000005', 'WHOLESALE', 'Wholesale Customers', 'Bulk buyers and distributors', 20.00, true);

-- Sample Customers
INSERT INTO customers (id, customer_code, name, tax_code, customer_type, customer_group_id, email, phone, payment_terms, credit_limit, currency, status) VALUES
    ('b0000001-0001-0001-0001-000000000001', 'CUST-0001', 'Beauty Shop Sài Gòn', '0312345678', 'RETAIL', 'a0000001-0001-0001-0001-000000000002', 'contact@beautyshop.vn', '028-1234-5678', 'Net 30', 50000000, 'VND', 'ACTIVE'),
    ('b0000001-0001-0001-0001-000000000002', 'CUST-0002', 'Hasaki Beauty & Clinic', '0312345679', 'WHOLESALE', 'a0000001-0001-0001-0001-000000000001', 'purchasing@hasaki.vn', '028-9876-5432', 'Net 45', 500000000, 'VND', 'ACTIVE'),
    ('b0000001-0001-0001-0001-000000000003', 'CUST-0003', 'Guardian Vietnam', '0312345680', 'DISTRIBUTOR', 'a0000001-0001-0001-0001-000000000005', 'supply@guardian.vn', '028-5555-1234', 'Net 60', 1000000000, 'VND', 'ACTIVE'),
    ('b0000001-0001-0001-0001-000000000004', 'CUST-0004', 'Spa Boutique Hà Nội', '0112345678', 'RETAIL', 'a0000001-0001-0001-0001-000000000003', 'info@spaboutique.vn', '024-1234-5678', 'Net 15', 20000000, 'VND', 'ACTIVE');

-- Sample Customer Addresses
INSERT INTO customer_addresses (customer_id, address_type, address_line1, city, country, is_default) VALUES
    ('b0000001-0001-0001-0001-000000000001', 'BOTH', '123 Nguyễn Huệ, Quận 1', 'Hồ Chí Minh', 'Vietnam', true),
    ('b0000001-0001-0001-0001-000000000002', 'BILLING', '456 Lê Lợi, Quận 1', 'Hồ Chí Minh', 'Vietnam', true),
    ('b0000001-0001-0001-0001-000000000002', 'SHIPPING', '789 Điện Biên Phủ, Quận 3', 'Hồ Chí Minh', 'Vietnam', false),
    ('b0000001-0001-0001-0001-000000000003', 'BOTH', '100 Trần Hưng Đạo, Quận 5', 'Hồ Chí Minh', 'Vietnam', true),
    ('b0000001-0001-0001-0001-000000000004', 'BOTH', '50 Hoàn Kiếm, Quận Hoàn Kiếm', 'Hà Nội', 'Vietnam', true);

-- Sample Customer Contacts  
INSERT INTO customer_contacts (customer_id, contact_name, position, email, phone, is_primary) VALUES
    ('b0000001-0001-0001-0001-000000000001', 'Nguyễn Văn A', 'Quản lý cửa hàng', 'nguyenvana@beautyshop.vn', '0901234567', true),
    ('b0000001-0001-0001-0001-000000000002', 'Trần Thị B', 'Trưởng phòng mua hàng', 'tranthib@hasaki.vn', '0912345678', true),
    ('b0000001-0001-0001-0001-000000000002', 'Lê Văn C', 'Nhân viên kế toán', 'levanc@hasaki.vn', '0923456789', false),
    ('b0000001-0001-0001-0001-000000000003', 'Phạm Thị D', 'Giám đốc cung ứng', 'phamthid@guardian.vn', '0934567890', true),
    ('b0000001-0001-0001-0001-000000000004', 'Hoàng Văn E', 'Chủ spa', 'hoangvane@spaboutique.vn', '0945678901', true);
