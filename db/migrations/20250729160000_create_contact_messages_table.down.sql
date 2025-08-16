-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_contact_messages_updated_at ON contact_messages;
DROP FUNCTION IF EXISTS update_contact_messages_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_contact_messages_service_interest;
DROP INDEX IF EXISTS idx_contact_messages_email;
DROP INDEX IF EXISTS idx_contact_messages_created_at;
DROP INDEX IF EXISTS idx_contact_messages_status;

-- Drop table
DROP TABLE IF EXISTS contact_messages; 