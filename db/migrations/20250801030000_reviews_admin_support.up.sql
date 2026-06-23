-- Allow admin-authored reviews: no purchasing user, custom display name.
ALTER TABLE reviews ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS reviewer_name VARCHAR(150);
