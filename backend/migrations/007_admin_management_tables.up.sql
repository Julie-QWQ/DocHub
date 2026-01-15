-- 访问日志表 (用于统计网站访问次数)
CREATE TABLE IF NOT EXISTS access_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    ip_address VARCHAR(50),
    path VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    user_agent TEXT,
    referer VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_access_logs_user_id ON access_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_access_logs_created_at ON access_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_access_logs_path ON access_logs(path);

-- 系统配置表 (如果不存在)
CREATE TABLE IF NOT EXISTS system_configs (
    id SERIAL PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value TEXT,
    description VARCHAR(255),
    category VARCHAR(50) DEFAULT 'general',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 预置配置
INSERT INTO system_configs (config_key, config_value, description, category) VALUES
('max_upload_size', '52428800', '最大上传文件大小（字节）', 'upload'),
('allowed_file_types', 'pdf,docx,doc,pptx,ppt,zip,rar', '允许的文件类型', 'upload'),
('download_expire_time', '3600', '下载链接有效期（秒）', 'download'),
('registration_enabled', 'true', '是否开放注册', 'auth'),
('site_name', 'Study-UPC', '网站名称', 'general'),
('site_description', '学院学习资料托管平台', '网站描述', 'general'),
('maintenance_mode', 'false', '维护模式', 'general')
ON CONFLICT (config_key) DO NOTHING;

CREATE UNIQUE INDEX IF NOT EXISTS uk_system_configs_key ON system_configs(config_key);

-- 更新 updated_at 触发器 (如果系统配置表还没有)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_system_configs_updated_at ON system_configs;
CREATE TRIGGER update_system_configs_updated_at BEFORE UPDATE ON system_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
