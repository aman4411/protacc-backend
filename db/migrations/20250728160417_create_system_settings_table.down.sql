-- Drop trigger and function
DROP TRIGGER IF EXISTS update_system_settings_updated_at_trigger ON system_settings;
DROP FUNCTION IF EXISTS update_system_settings_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_system_settings_category;
DROP INDEX IF EXISTS idx_system_settings_public;

-- Drop table
DROP TABLE IF EXISTS system_settings;
