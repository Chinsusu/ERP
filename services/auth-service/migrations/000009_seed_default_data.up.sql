-- Seed default system roles
INSERT INTO roles (name, description, is_system) VALUES
    ('Super Admin', 'Full system access - all permissions', true),
    ('Admin', 'Administrative access to most features', true),
    ('Manager', 'Department manager with approval rights', true),
    ('Staff', 'Standard user with basic permissions', true),
    ('Viewer', 'Read-only access to system', true)
ON CONFLICT (name) DO NOTHING;

-- Seed default permissions
INSERT INTO permissions (code, name, service, resource, action, description) VALUES
    -- Wildcard - Super Admin only
    ('*:*:*', 'Full Access', '*', '*', '*', 'Super admin - all permissions'),
    
    -- Auth Service
    ('auth:user:read', 'View Users', 'auth', 'user', 'read', 'View user credentials'),
    ('auth:user:create', 'Create User', 'auth', 'user', 'create', 'Create new user credentials'),
    ('auth:user:update', 'Update User', 'auth', 'user', 'update', 'Update user credentials'),
    ('auth:user:delete', 'Delete User', 'auth', 'user', 'delete', 'Delete user credentials'),
    ('auth:role:read', 'View Roles', 'auth', 'role', 'read', 'View roles'),
    ('auth:role:create', 'Create Role', 'auth', 'role', 'create', 'Create new roles'),
    ('auth:role:update', 'Update Role', 'auth', 'role', 'update', 'Update roles'),
    ('auth:role:delete', 'Delete Role', 'auth', 'role', 'delete', 'Delete roles'),
    ('auth:permission:read', 'View Permissions', 'auth', 'permission', 'read', 'View permissions'),
    
    -- User Service
    ('user:user:read', 'View User Profiles', 'user', 'user', 'read', 'View user profiles and info'),
    ('user:user:create', 'Create User Profile', 'user', 'user', 'create', 'Create user profiles'),
    ('user:user:update', 'Update User Profile', 'user', 'user', 'update', 'Update user profiles'),
    ('user:user:delete', 'Delete User Profile', 'user', 'user', 'delete', 'Delete user profiles'),
    ('user:department:read', 'View Departments', 'user', 'department', 'read', 'View department hierarchy'),
    ('user:department:create', 'Create Department', 'user', 'department', 'create', 'Create departments'),
    
    -- Procurement Service
    ('procurement:pr:read', 'View PR', 'procurement', 'pr', 'read', 'View purchase requisitions'),
    ('procurement:pr:create', 'Create PR', 'procurement', 'pr', 'create', 'Create purchase requisitions'),
    ('procurement:pr:update', 'Update PR', 'procurement', 'pr', 'update', 'Update purchase requisitions'),
    ('procurement:pr:approve', 'Approve PR', 'procurement', 'pr', 'approve', 'Approve purchase requisitions'),
    ('procurement:po:read', 'View PO', 'procurement', 'po', 'read', 'View purchase orders'),
    ('procurement:po:create', 'Create PO', 'procurement', 'po', 'create', 'Create purchase orders'),
    ('procurement:po:approve', 'Approve PO', 'procurement', 'po', 'approve', 'Approve purchase orders'),
    
    -- WMS Service
    ('wms:stock:read', 'View Stock', 'wms', 'stock', 'read', 'View inventory and stock levels'),
    ('wms:stock:adjust', 'Adjust Stock', 'wms', 'stock', 'adjust', 'Perform stock adjustments'),
    ('wms:grn:read', 'View GRN', 'wms', 'grn', 'read', 'View goods receipt notes'),
    ('wms:grn:create', 'Create GRN', 'wms', 'grn', 'create', 'Create goods receipt notes'),
    ('wms:lot:read', 'View Lots', 'wms', 'lot', 'read', 'View lot/batch information'),
    
    -- Manufacturing Service
    ('manufacturing:bom:read', 'View BOM', 'manufacturing', 'bom', 'read', 'View bill of materials'),
    ('manufacturing:bom:create', 'Create BOM', 'manufacturing', 'bom', 'create', 'Create bill of materials'),
    ('manufacturing:bom:approve', 'Approve BOM', 'manufacturing', 'bom', 'approve', 'Approve BOM (sensitive!)'),
    ('manufacturing:wo:read', 'View Work Orders', 'manufacturing', 'wo', 'read', 'View work orders'),
    ('manufacturing:wo:create', 'Create Work Orders', 'manufacturing', 'wo', 'create', 'Create work orders'),
    
    -- Sales Service
    ('sales:customer:read', 'View Customers', 'sales', 'customer', 'read', 'View customer information'),
   ('sales:customer:create', 'Create Customer', 'sales', 'customer', 'create', 'Create new customers'),
    ('sales:order:read', 'View Sales Orders', 'sales', 'order', 'read', 'View sales orders'),
    ('sales:order:create', 'Create Sales Order', 'sales', 'order', 'create', 'Create sales orders'),
    ('sales:order:approve', 'Approve Sales Order', 'sales', 'order', 'approve', 'Approve sales orders'),
    
    -- Marketing Service
    ('marketing:campaign:read', 'View Campaigns', 'marketing', 'campaign', 'read', 'View marketing campaigns'),
    ('marketing:campaign:create', 'Create Campaign', 'marketing', 'campaign', 'create', 'Create campaigns'),
    ('marketing:kol:read', 'View KOLs', 'marketing', 'kol', 'read', 'View KOL database'),
    ('marketing:kol:create', 'Create KOL', 'marketing', 'kol', 'create', 'Add new KOLs')
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to Super Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Super Admin'),
    id
FROM permissions
WHERE code = '*:*:*'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign admin permissions (all except wildcard)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Admin'),
    id
FROM permissions
WHERE service IN ('auth', 'user', 'procurement', 'wms', 'manufacturing', 'sales', 'marketing')
  AND code != '*:*:*'
  AND action IN ('read', 'create', 'update')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign manager permissions (read + create, no delete/approve)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Manager'),
    id
FROM permissions
WHERE code != '*:*:*'
  AND action IN ('read', 'create')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign staff permissions (mostly read)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Staff'),
    id
FROM permissions
WHERE code != '*:*:*'
  AND action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign viewer permissions (only read)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'Viewer'),
    id
FROM permissions
WHERE action = 'read'
  AND code != '*:*:*'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Create default admin user
DO $$
DECLARE
    v_user_id UUID := gen_random_uuid();
    v_super_admin_role_id UUID;
BEGIN
    -- Get Super Admin role ID
    SELECT id INTO v_super_admin_role_id FROM roles WHERE name = 'Super Admin';
    
    -- Insert admin user credentials
    INSERT INTO user_credentials (
        id,
        user_id,
        email,
        password_hash,
        is_active,
        email_verified
    ) VALUES (
        gen_random_uuid(),
        v_user_id,
        'admin@company.vn',
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLHJ7.W6', -- Admin@123
        true,
        true
    )
    ON CONFLICT (email) DO NOTHING;
    
    -- Assign Super Admin role to admin user
    INSERT INTO user_roles (user_id, role_id)
    VALUES (v_user_id, v_super_admin_role_id)
    ON CONFLICT (user_id, role_id) DO NOTHING;
    
END $$;
