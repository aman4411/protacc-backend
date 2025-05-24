-- Modify icon column to support full URLs
ALTER TABLE service_categories ALTER COLUMN icon TYPE VARCHAR(255);

-- Update existing icons to use image URLs
UPDATE service_categories
SET icon = CASE 
    WHEN name = 'Business Registration' THEN '/images/categories/business-registration.svg'
    WHEN name = 'Tax & Compliance' THEN '/images/categories/tax-compliance.svg'
    WHEN name = 'Trademark & IP' THEN '/images/categories/trademark-ip.svg'
    WHEN name = 'Digital Services' THEN '/images/categories/digital-services.svg'
    ELSE '/images/categories/default.svg'
END; 