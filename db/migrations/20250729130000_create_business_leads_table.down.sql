-- Drop trigger
DROP TRIGGER IF EXISTS update_business_leads_updated_at_trigger ON business_leads;

-- Drop function
DROP FUNCTION IF EXISTS update_business_leads_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_business_leads_phone;
DROP INDEX IF EXISTS idx_business_leads_email;
DROP INDEX IF EXISTS idx_business_leads_created_at;
DROP INDEX IF EXISTS idx_business_leads_assigned_to;
DROP INDEX IF EXISTS idx_business_leads_priority;
DROP INDEX IF EXISTS idx_business_leads_status;

-- Drop table
DROP TABLE IF EXISTS business_leads; 