-- add email verification table and user flag

ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS email_verification_codes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    email VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    purpose VARCHAR(20) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_email_verification_codes_email ON email_verification_codes(email);
CREATE INDEX IF NOT EXISTS idx_email_verification_codes_expires_at ON email_verification_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_email_verification_codes_purpose ON email_verification_codes(purpose);

COMMENT ON COLUMN users.email_verified IS 'email verified';
COMMENT ON COLUMN email_verification_codes.email IS 'email address';
COMMENT ON COLUMN email_verification_codes.code IS 'verification code';
COMMENT ON COLUMN email_verification_codes.expires_at IS 'expires at';
COMMENT ON COLUMN email_verification_codes.is_used IS 'is used';
COMMENT ON COLUMN email_verification_codes.used_at IS 'used at';
COMMENT ON COLUMN email_verification_codes.purpose IS 'purpose: register/login/reset_password';
