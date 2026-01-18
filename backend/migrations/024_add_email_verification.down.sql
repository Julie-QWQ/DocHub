-- rollback email verification table and user flag

DROP INDEX IF EXISTS idx_email_verification_codes_email;
DROP INDEX IF EXISTS idx_email_verification_codes_expires_at;
DROP INDEX IF EXISTS idx_email_verification_codes_purpose;
DROP TABLE IF EXISTS email_verification_codes;

ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
