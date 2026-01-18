-- 添加邮箱验证相关字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE;

-- 创建邮箱验证码表
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

-- 创建索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_email_verification_codes_email ON email_verification_codes(email);
CREATE INDEX IF NOT EXISTS idx_email_verification_codes_expires_at ON email_verification_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_email_verification_codes_purpose ON email_verification_codes(purpose);

-- 添加注释
COMMENT ON COLUMN users.email_verified IS '邮箱是否已验证';
COMMENT ON COLUMN email_verification_codes.email IS '邮箱地址';
COMMENT ON COLUMN email_verification_codes.code IS '验证码';
COMMENT ON COLUMN email_verification_codes.expires_at IS '过期时间';
COMMENT ON COLUMN email_verification_codes.is_used IS '是否已使用';
COMMENT ON COLUMN email_verification_codes.used_at IS '使用时间';
COMMENT ON COLUMN email_verification_codes.purpose IS '用途: register/login/reset_password';
