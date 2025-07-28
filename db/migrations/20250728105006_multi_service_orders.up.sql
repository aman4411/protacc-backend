-- Remove service_id from orders table as we'll use order_items instead
ALTER TABLE orders DROP COLUMN service_id;

-- Create order_items table to store multiple services per order
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    service_id INTEGER REFERENCES services(id),
    quantity INTEGER DEFAULT 1,
    price DECIMAL(10,2) NOT NULL,
    booking_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster queries
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_service_id ON order_items(service_id);

-- Update existing orders to use order_items structure
-- Note: This handles any existing orders by creating corresponding order_items
-- Since we don't have service_id anymore, we'll need to handle this in the application
