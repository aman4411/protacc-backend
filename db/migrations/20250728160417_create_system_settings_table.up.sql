-- Create system_settings table for storing application configuration
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL,
    setting_key VARCHAR(100) NOT NULL,
    setting_value TEXT,
    data_type VARCHAR(20) NOT NULL DEFAULT 'string', -- 'string', 'number', 'boolean', 'json'
    description TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT FALSE, -- Whether this setting can be accessed by non-admin users
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(category, setting_key)
);

-- Create index for faster lookups
CREATE INDEX idx_system_settings_category ON system_settings(category);
CREATE INDEX idx_system_settings_public ON system_settings(is_public);

-- Insert default system settings
INSERT INTO system_settings (category, setting_key, setting_value, data_type, description, is_public) VALUES
-- General Settings
('general', 'site_name', 'ProtAcc', 'string', 'Website name displayed to users', true),
('general', 'site_description', 'Professional Chartered Accountant Services', 'string', 'Website description for SEO', true),
('general', 'admin_email', 'admin@protacc.com', 'string', 'Primary admin email address', false),
('general', 'support_email', 'support@protacc.com', 'string', 'Support email address', true),
('general', 'phone_number', '+1-234-567-8900', 'string', 'Primary contact phone number', true),
('general', 'address', '123 Business St, Finance City, FC 12345', 'string', 'Business address', true),
('general', 'timezone', 'UTC', 'string', 'Default timezone for the application', false),
('general', 'date_format', 'YYYY-MM-DD', 'string', 'Default date format', false),
('general', 'currency', 'INR', 'string', 'Default currency', true),
('general', 'currency_symbol', '₹', 'string', 'Currency symbol', true),

-- Email Settings
('email', 'smtp_host', 'smtp.gmail.com', 'string', 'SMTP server host', false),
('email', 'smtp_port', '587', 'number', 'SMTP server port', false),
('email', 'smtp_username', '', 'string', 'SMTP username', false),
('email', 'smtp_password', '', 'string', 'SMTP password', false),
('email', 'smtp_encryption', 'tls', 'string', 'SMTP encryption (tls/ssl/none)', false),
('email', 'from_email', 'noreply@protacc.com', 'string', 'Default from email address', false),
('email', 'from_name', 'ProtAcc', 'string', 'Default from name', false),

-- Payment Settings
('payment', 'default_booking_amount', '99', 'number', 'Default booking amount for services', false),
('payment', 'payment_gateway', 'razorpay', 'string', 'Primary payment gateway', false),
('payment', 'razorpay_key_id', '', 'string', 'Razorpay Key ID', false),
('payment', 'razorpay_key_secret', '', 'string', 'Razorpay Key Secret', false),
('payment', 'enable_cod', 'false', 'boolean', 'Enable Cash on Delivery', false),
('payment', 'tax_rate', '18', 'number', 'Default tax rate percentage', false),

-- Notification Settings
('notification', 'email_notifications', 'true', 'boolean', 'Enable email notifications', false),
('notification', 'sms_notifications', 'false', 'boolean', 'Enable SMS notifications', false),
('notification', 'push_notifications', 'true', 'boolean', 'Enable push notifications', false),
('notification', 'order_confirmation_email', 'true', 'boolean', 'Send order confirmation emails', false),
('notification', 'status_update_email', 'true', 'boolean', 'Send status update emails', false),

-- Security Settings
('security', 'jwt_expiry_hours', '24', 'number', 'JWT token expiry in hours', false),
('security', 'password_min_length', '8', 'number', 'Minimum password length', false),
('security', 'require_email_verification', 'false', 'boolean', 'Require email verification for new users', false),
('security', 'login_attempts_limit', '5', 'number', 'Maximum login attempts before lockout', false),
('security', 'lockout_duration_minutes', '30', 'number', 'Account lockout duration in minutes', false),

-- Business Settings
('business', 'business_hours', '{"monday":"9:00-18:00","tuesday":"9:00-18:00","wednesday":"9:00-18:00","thursday":"9:00-18:00","friday":"9:00-18:00","saturday":"9:00-13:00","sunday":"closed"}', 'json', 'Business operating hours', true),
('business', 'service_categories_limit', '10', 'number', 'Maximum number of service categories', false),
('business', 'services_per_category_limit', '50', 'number', 'Maximum services per category', false),

-- UI/UX Settings
('ui', 'theme_color', '#4f46e5', 'string', 'Primary theme color', true),
('ui', 'secondary_color', '#e5e7eb', 'string', 'Secondary theme color', true),
('ui', 'items_per_page', '10', 'number', 'Default items per page for listings', true),
('ui', 'enable_dark_mode', 'false', 'boolean', 'Enable dark mode option', true),

-- SEO Settings
('seo', 'meta_title', 'ProtAcc - Professional Chartered Accountant Services', 'string', 'Default meta title', true),
('seo', 'meta_description', 'Professional chartered accountant services including tax compliance, business registration, and financial consulting.', 'string', 'Default meta description', true),
('seo', 'meta_keywords', 'chartered accountant, tax services, business registration, financial consulting', 'string', 'Default meta keywords', true),
('seo', 'google_analytics_id', '', 'string', 'Google Analytics tracking ID', false),
('seo', 'google_tag_manager_id', '', 'string', 'Google Tag Manager ID', false);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_system_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_system_settings_updated_at_trigger
    BEFORE UPDATE ON system_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_system_settings_updated_at();
