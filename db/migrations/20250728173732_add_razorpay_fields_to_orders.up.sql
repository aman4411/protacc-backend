-- Add Razorpay related fields to orders table
ALTER TABLE orders 
ADD COLUMN razorpay_order_id VARCHAR(255),
ADD COLUMN razorpay_payment_id VARCHAR(255),
ADD COLUMN payment_method VARCHAR(100),
ADD COLUMN payment_gateway VARCHAR(50) DEFAULT 'razorpay';

-- Add index for faster lookup by Razorpay order ID
CREATE INDEX IF NOT EXISTS idx_orders_razorpay_order_id ON orders(razorpay_order_id);
