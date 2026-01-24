DROP TRIGGER IF EXISTS user_profiles_updated_at ON user_profiles;
DROP FUNCTION IF EXISTS update_user_profiles_updated_at();
DROP TABLE IF EXISTS user_profiles CASCADE;
