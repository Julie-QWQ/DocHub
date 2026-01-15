-- 添加 users 表缺失的字段
-- 版本: 002
-- 描述: 修复用户表缺失的字段

-- 添加缺失的字段(如果不存在)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS real_name VARCHAR(50),
    ADD COLUMN IF NOT EXISTS avatar VARCHAR(255),
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS major VARCHAR(100),
    ADD COLUMN IF NOT EXISTS class VARCHAR(50),
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 修改 status 默认值为 inactive (如果还是 active)
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'inactive';
