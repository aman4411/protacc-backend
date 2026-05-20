-- Rollback existing order status updates
-- Revert status updates
UPDATE orders SET status = 'pending_payment' WHERE status = 'pending_booking_payment';
UPDATE order_status_history SET status = 'pending_payment' WHERE status = 'pending_booking_payment';

UPDATE orders SET status = 'payment_received' WHERE status = 'booking_amount_received';
UPDATE order_status_history SET status = 'payment_received' WHERE status = 'booking_amount_received';

-- Revert default status
ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'pending_payment';




