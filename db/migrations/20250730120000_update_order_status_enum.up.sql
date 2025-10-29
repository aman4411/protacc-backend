-- Update order_status enum to support two-stage payment flow
-- Add new enum values (must be done separately from using them)
ALTER TYPE order_status ADD VALUE 'pending_booking_payment';
ALTER TYPE order_status ADD VALUE 'booking_amount_received';
ALTER TYPE order_status ADD VALUE 'pending_final_payment';
ALTER TYPE order_status ADD VALUE 'full_payment_received';
