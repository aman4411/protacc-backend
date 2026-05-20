-- Rollback order_status enum changes
-- Revert status updates
UPDATE orders SET status = 'pending_payment' WHERE status = 'pending_booking_payment';
UPDATE order_status_history SET status = 'pending_payment' WHERE status = 'pending_booking_payment';

UPDATE orders SET status = 'payment_received' WHERE status = 'booking_amount_received';
UPDATE order_status_history SET status = 'payment_received' WHERE status = 'booking_amount_received';

-- Revert default status
ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'pending_payment';

-- Note: PostgreSQL doesn't support removing enum values directly
-- This would require recreating the enum type, which is complex
-- For now, the new enum values will remain but won't be used

