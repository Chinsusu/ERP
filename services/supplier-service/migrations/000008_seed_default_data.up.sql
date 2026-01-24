-- Seed default data for Supplier Service

-- Sample suppliers
INSERT INTO suppliers (id, code, name, legal_name, tax_code, supplier_type, business_type, email, phone, currency, payment_terms, status, overall_rating, created_at, updated_at) VALUES
('b0000001-0000-0000-0000-000000000001', 'SUP-0001', 'ABC Chemicals Vietnam', 'Công ty TNHH ABC Chemicals Vietnam', '0123456789', 'MANUFACTURER', 'DOMESTIC', 'sales@abc-chemicals.vn', '+84 28 1234 5678', 'VND', 'Net 30', 'APPROVED', 4.5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('b0000001-0000-0000-0000-000000000002', 'SUP-0002', 'Natural Ingredients Ltd', 'Natural Ingredients Company Limited', '9876543210', 'MANUFACTURER', 'INTERNATIONAL', 'contact@naturalingredients.com', '+66 2 1234 5678', 'USD', 'Net 45', 'APPROVED', 4.2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('b0000001-0000-0000-0000-000000000003', 'SUP-0003', 'VN Packaging Co', 'Công ty Bao Bì Việt Nam', '1122334455', 'MANUFACTURER', 'DOMESTIC', 'info@vnpackaging.vn', '+84 28 9876 5432', 'VND', 'Net 30', 'APPROVED', 4.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('b0000001-0000-0000-0000-000000000004', 'SUP-0004', 'Organic Essentials Thailand', 'Organic Essentials (Thailand) Co., Ltd', 'TH987654321', 'TRADER', 'INTERNATIONAL', 'sales@organicth.com', '+66 2 9999 8888', 'USD', 'Net 60', 'PENDING', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Sample addresses
INSERT INTO supplier_addresses (id, supplier_id, address_type, address_line1, district, city, country, postal_code, is_primary) VALUES
('c0000001-0000-0000-0000-000000000001', 'b0000001-0000-0000-0000-000000000001', 'FACTORY', '123 Industrial Park, Lot A5', 'Binh Tan', 'Ho Chi Minh City', 'Vietnam', '70000', TRUE),
('c0000001-0000-0000-0000-000000000002', 'b0000001-0000-0000-0000-000000000001', 'BILLING', '456 Nguyen Van Linh, Floor 5', 'District 7', 'Ho Chi Minh City', 'Vietnam', '70000', FALSE),
('c0000001-0000-0000-0000-000000000003', 'b0000001-0000-0000-0000-000000000002', 'FACTORY', '789 Bangkok Industrial Estate', 'Bang Pu', 'Samut Prakan', 'Thailand', '10280', TRUE);

-- Sample contacts
INSERT INTO supplier_contacts (id, supplier_id, contact_type, full_name, position, department, email, phone, mobile, is_primary) VALUES
('d0000001-0000-0000-0000-000000000001', 'b0000001-0000-0000-0000-000000000001', 'PRIMARY', 'Nguyễn Văn A', 'Sales Manager', 'Sales', 'nguyen.a@abc-chemicals.vn', '+84 28 1234 5678', '+84 901 234 567', TRUE),
('d0000001-0000-0000-0000-000000000002', 'b0000001-0000-0000-0000-000000000001', 'TECHNICAL', 'Trần Thị B', 'Technical Support', 'R&D', 'tran.b@abc-chemicals.vn', '+84 28 1234 5679', '+84 902 345 678', FALSE),
('d0000001-0000-0000-0000-000000000003', 'b0000001-0000-0000-0000-000000000001', 'QUALITY', 'Lê Văn C', 'QC Manager', 'Quality Control', 'le.c@abc-chemicals.vn', '+84 28 1234 5680', '+84 903 456 789', FALSE),
('d0000001-0000-0000-0000-000000000004', 'b0000001-0000-0000-0000-000000000002', 'PRIMARY', 'Somchai T.', 'Export Manager', 'Sales', 'somchai@naturalingredients.com', '+66 2 1234 5679', '+66 81 234 5678', TRUE);

-- Sample certifications
INSERT INTO supplier_certifications (id, supplier_id, certification_type, certificate_number, issuing_body, issue_date, expiry_date, status) VALUES
('e0000001-0000-0000-0000-000000000001', 'b0000001-0000-0000-0000-000000000001', 'GMP', 'GMP-VN-2024-001', 'Vietnam Food Administration', '2024-01-01', '2027-01-01', 'VALID'),
('e0000001-0000-0000-0000-000000000002', 'b0000001-0000-0000-0000-000000000001', 'ISO9001', 'ISO-2024-ABC-001', 'TUV Nord', '2024-01-01', '2026-12-31', 'VALID'),
('e0000001-0000-0000-0000-000000000003', 'b0000001-0000-0000-0000-000000000001', 'ISO22716', 'ISO22716-VN-001', 'SGS', '2023-06-01', '2026-05-31', 'VALID'),
('e0000001-0000-0000-0000-000000000004', 'b0000001-0000-0000-0000-000000000002', 'GMP', 'GMP-TH-2023-456', 'Thai FDA', '2023-01-01', '2026-01-01', 'EXPIRING_SOON'),
('e0000001-0000-0000-0000-000000000005', 'b0000001-0000-0000-0000-000000000002', 'ORGANIC', 'ORG-EU-2024-789', 'Ecocert', '2024-01-01', '2025-12-31', 'VALID'),
('e0000001-0000-0000-0000-000000000006', 'b0000001-0000-0000-0000-000000000003', 'ISO9001', 'ISO-2023-PKG-001', 'Bureau Veritas', '2023-03-01', '2026-02-28', 'VALID');

-- Sample evaluations
INSERT INTO supplier_evaluations (id, supplier_id, evaluation_date, evaluation_period, quality_score, delivery_score, price_score, service_score, documentation_score, overall_score, on_time_delivery_rate, quality_acceptance_rate, strengths, weaknesses, evaluated_by, status) VALUES
('f0000001-0000-0000-0000-000000000001', 'b0000001-0000-0000-0000-000000000001', '2024-01-15', '2023-Q4', 5.0, 4.5, 4.0, 5.0, 5.0, 4.7, 95.5, 98.2, 'Excellent quality, responsive technical support', 'Occasional delivery delays', '09c95223-32b6-4b50-87e7-4ea1333ae072', 'APPROVED'),
('f0000001-0000-0000-0000-000000000002', 'b0000001-0000-0000-0000-000000000002', '2024-01-15', '2023-Q4', 4.5, 4.0, 3.5, 4.0, 4.5, 4.1, 88.0, 96.5, 'Good organic certifications, wide product range', 'Higher prices, longer lead times', '09c95223-32b6-4b50-87e7-4ea1333ae072', 'APPROVED');
