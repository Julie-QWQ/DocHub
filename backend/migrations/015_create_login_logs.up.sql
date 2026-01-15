-- 创建登录日志表
CREATE TABLE IF NOT EXISTS login_logs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    ip_address VARCHAR(50),
    user_agent TEXT,
    success BOOLEAN NOT NULL DEFAULT true
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_created_at ON login_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_login_logs_success ON login_logs(success);

-- 创建复合索引用于查询今日活跃用户
CREATE INDEX IF NOT EXISTS idx_login_user_date ON login_logs(user_id, DATE(created_at));
