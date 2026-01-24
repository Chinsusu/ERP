DROP TRIGGER IF EXISTS departments_updated_at ON departments;
DROP FUNCTION IF EXISTS update_departments_updated_at();
DROP TABLE IF EXISTS departments CASCADE;
