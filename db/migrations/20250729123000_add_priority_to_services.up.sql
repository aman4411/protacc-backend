-- Add priority column to service_categories table
ALTER TABLE service_categories ADD COLUMN priority INTEGER DEFAULT 0;

-- Add priority column to services table  
ALTER TABLE services ADD COLUMN priority INTEGER DEFAULT 0;

-- Update existing categories with default priorities (lower number = higher priority)
UPDATE service_categories SET priority = id * 10;

-- Update existing services with default priorities  
UPDATE services SET priority = id * 10;

-- Create indexes for better performance on priority ordering
CREATE INDEX idx_service_categories_priority ON service_categories(priority ASC, created_at DESC);
CREATE INDEX idx_services_priority ON services(priority ASC, created_at DESC);

-- Create index for services ordered by category priority and service priority
CREATE INDEX idx_services_category_priority ON services(category_id, priority ASC, created_at DESC); 