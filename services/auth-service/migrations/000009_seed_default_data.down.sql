-- Remove seeded data (roles and permissions will cascade delete)
DELETE FROM user_roles WHERE user_id IN (
    SELECT user_id FROM user_credentials WHERE email = 'admin@company.vn'
);

DELETE FROM user_credentials WHERE email = 'admin@company.vn';

-- Note: We don't delete roles and permissions as they might be in use
-- Only delete the default admin user
