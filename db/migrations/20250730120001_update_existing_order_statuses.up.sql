-- Update existing orders to use new status names
-- This must be done in a separate migration after adding enum values

-- pending_payment -> pending_booking_payment
UPDATE orders SET status = 'pending_booking_payment' WHERE status = 'pending_payment';
UPDATE order_status_history SET status = 'pending_booking_payment' WHERE status = 'pending_payment';

-- payment_received -> booking_amount_received (for existing orders)
UPDATE orders SET status = 'booking_amount_received' WHERE status = 'payment_received';
UPDATE order_status_history SET status = 'booking_amount_received' WHERE status = 'payment_received';

-- Update default status for new orders
ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'pending_booking_payment';
