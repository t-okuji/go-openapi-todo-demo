-- Migration rollback: Remove initial schema
-- Description: Drop all tables, triggers, and functions created in initial migration

-- Drop triggers
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes (they'll be dropped automatically with tables, but explicitly for clarity)
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_completed;
DROP INDEX IF EXISTS idx_todos_category_id;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS categories;

-- Drop extension (only if not used by other parts of the database)
-- DROP EXTENSION IF EXISTS "uuid-ossp";