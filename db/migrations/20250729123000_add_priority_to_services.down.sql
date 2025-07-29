-- Drop indexes
DROP INDEX IF EXISTS idx_services_category_priority;
DROP INDEX IF EXISTS idx_services_priority;
DROP INDEX IF EXISTS idx_service_categories_priority;

-- Remove priority columns
ALTER TABLE services DROP COLUMN IF EXISTS priority;
ALTER TABLE service_categories DROP COLUMN IF EXISTS priority; 