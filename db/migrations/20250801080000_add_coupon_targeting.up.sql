ALTER TABLE coupons ADD COLUMN IF NOT EXISTS applicable_category_ids INTEGER[];
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS applicable_service_ids INTEGER[];
