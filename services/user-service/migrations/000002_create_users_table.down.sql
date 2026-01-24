ALTER TABLE departments DROP CONSTRAINT IF EXISTS fk_departments_manager;
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP FUNCTION IF EXISTS update_users_updated_at();
DROP TABLE IF EXISTS users CASCADE;
