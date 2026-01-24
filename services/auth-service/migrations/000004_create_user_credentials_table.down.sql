-- Drop user_credentials table
DROP TRIGGER IF EXISTS update_user_credentials_updated_at ON user_credentials;
DROP TABLE IF EXISTS user_credentials CASCADE;
