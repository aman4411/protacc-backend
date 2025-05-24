CREATE TABLE IF NOT EXISTS email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    otp VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT unique_active_verification UNIQUE (email, otp, expires_at)
);

CREATE INDEX IF NOT EXISTS idx_email_verifications_email_otp ON email_verifications(email, otp);
CREATE INDEX IF NOT EXISTS idx_email_verifications_user_id ON email_verifications(user_id);