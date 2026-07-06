ALTER TABLE services ADD COLUMN IF NOT EXISTS min_delivery_days INTEGER;
ALTER TABLE services ADD COLUMN IF NOT EXISTS max_delivery_days INTEGER;

-- Backfill the new range from the existing single estimate.
UPDATE services
SET min_delivery_days = COALESCE(estimated_delivery_days, 0),
    max_delivery_days = COALESCE(estimated_delivery_days, 0)
WHERE min_delivery_days IS NULL OR max_delivery_days IS NULL;

ALTER TABLE services DROP COLUMN IF EXISTS estimated_delivery_days;
