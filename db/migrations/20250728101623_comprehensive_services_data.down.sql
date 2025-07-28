-- Clear all data added by the comprehensive services migration
DELETE FROM order_status_history;
DELETE FROM orders;
DELETE FROM cart_items;
DELETE FROM services;
DELETE FROM service_categories;
