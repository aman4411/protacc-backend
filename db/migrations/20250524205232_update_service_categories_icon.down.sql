-- Revert icon column to original size and update values to simple identifiers
UPDATE service_categories
SET icon = CASE 
    WHEN name = 'Business Registration' THEN 'building'
    WHEN name = 'Tax & Compliance' THEN 'file-invoice'
    WHEN name = 'Trademark & IP' THEN 'trademark'
    WHEN name = 'Digital Services' THEN 'laptop'
    ELSE NULL
END;

ALTER TABLE service_categories ALTER COLUMN icon TYPE VARCHAR(50); 