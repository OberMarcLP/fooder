-- Rollback authentication changes

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS sessions;

-- Remove columns from existing tables
ALTER TABLE restaurants DROP COLUMN IF EXISTS updated_by;
ALTER TABLE restaurants DROP COLUMN IF EXISTS created_by;
ALTER TABLE suggestions DROP COLUMN IF EXISTS user_id;
ALTER TABLE ratings DROP COLUMN IF EXISTS user_id;

-- Drop users table
DROP TABLE IF EXISTS users;
