-- Drop order_items table
DROP TABLE IF EXISTS order_items;

-- Add back service_id column to orders table
ALTER TABLE orders ADD COLUMN service_id INTEGER REFERENCES services(id);
