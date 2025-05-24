-- Drop indexes
DROP INDEX IF EXISTS idx_services_category_id;
DROP INDEX IF EXISTS idx_cart_items_user_id;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_service_id;
DROP INDEX IF EXISTS idx_order_status_history_order_id;
DROP INDEX IF EXISTS idx_order_status_history_created_by;

-- Drop tables
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS service_categories;

-- Drop enum types
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS service_status;