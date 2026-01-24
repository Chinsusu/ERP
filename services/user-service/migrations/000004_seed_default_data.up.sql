-- Seed default departments
INSERT INTO departments (code, name, parent_id, level, path, status) VALUES
    ('EXEC', 'Executive', NULL, 0, '/EXEC/', 'active'),
    ('IT', 'Information Technology', NULL, 0, '/IT/', 'active'),
    ('HR', 'Human Resources', NULL, 0, '/HR/', 'active'),
    ('FIN', 'Finance', NULL, 0, '/FIN/', 'active'),
    ('OPS', 'Operations', NULL, 0, '/OPS/', 'active'),
    ('MFG', 'Manufacturing', (SELECT id FROM departments WHERE code = 'OPS'), 1, '/OPS/MFG/', 'active'),
    ('WH', 'Warehouse', (SELECT id FROM departments WHERE code = 'OPS'), 1, '/OPS/WH/', 'active'),
    ('QC', 'Quality Control', (SELECT id FROM departments WHERE code = 'MFG'), 2, '/OPS/MFG/QC/', 'active')
ON CONFLICT (code) DO NOTHING;

-- Seed sample users (matching with auth service admin)
INSERT INTO users (id, email, employee_code, first_name, last_name, phone, department_id, status) VALUES
    ('09c95223-32b6-4b50-87e7-4ea1333ae072', 'admin@company.vn', 'EMP20260124001', 'System', 'Administrator', '+84123456789', (SELECT id FROM departments WHERE code = 'IT'), 'active')
ON CONFLICT (email) DO NOTHING;

-- Seed user profile for admin
INSERT INTO user_profiles (user_id, join_date, address, emergency_contact) VALUES
    ('09c95223-32b6-4b50-87e7-4ea1333ae072', '2026-01-01', 'Ho Chi Minh City, Vietnam', '+84987654321')
ON CONFLICT (user_id) DO NOTHING;
