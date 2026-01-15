-- 创建访问日志表
CREATE TABLE IF NOT EXISTS access_logs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id BIGINT,                      -- 用户ID (可为空,游客访问)
    ip_address VARCHAR(50),              -- IP地址
    path VARCHAR(255) NOT NULL,          -- 访问路径
    method VARCHAR(10) NOT NULL,         -- 请求方法
    user_agent TEXT,                     -- 用户代理
    referer VARCHAR(500)                 -- 来源页面
);

-- 创建索引以提升查询性能
CREATE INDEX IF NOT EXISTS idx_access_logs_user_id ON access_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_access_logs_created_at ON access_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_access_logs_path ON access_logs(path);
CREATE INDEX IF NOT EXISTS idx_access_logs_ip_address ON access_logs(ip_address);

-- 创建复合索引用于独立访客统计
CREATE INDEX IF NOT EXISTS idx_access_logs_ip_date ON access_logs(ip_address, DATE(created_at));

-- 创建复合索引用于用户访问统计
CREATE INDEX IF NOT EXISTS idx_access_logs_user_date ON access_logs(user_id, DATE(created_at));
