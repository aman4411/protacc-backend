ALTER TABLE services ADD COLUMN IF NOT EXISTS estimated_delivery_days INTEGER;

-- Restore the single estimate from the upper bound of the range.
UPDATE services
SET estimated_delivery_days = COALESCE(max_delivery_days, min_delivery_days, 0)
WHERE estimated_delivery_days IS NULL;

ALTER TABLE services DROP COLUMN IF EXISTS min_delivery_days;
ALTER TABLE services DROP COLUMN IF EXISTS max_delivery_days;
