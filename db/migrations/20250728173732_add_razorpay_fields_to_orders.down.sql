-- Remove Razorpay related fields from orders table
DROP INDEX IF EXISTS idx_orders_razorpay_order_id;

ALTER TABLE orders 
DROP COLUMN IF EXISTS razorpay_order_id,
DROP COLUMN IF EXISTS razorpay_payment_id,
DROP COLUMN IF EXISTS payment_method,
DROP COLUMN IF EXISTS payment_gateway;
