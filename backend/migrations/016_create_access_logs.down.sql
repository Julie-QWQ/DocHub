-- 删除索引
DROP INDEX IF EXISTS idx_access_logs_user_date;
DROP INDEX IF EXISTS idx_access_logs_ip_date;
DROP INDEX IF EXISTS idx_access_logs_ip_address;
DROP INDEX IF EXISTS idx_access_logs_path;
DROP INDEX IF EXISTS idx_access_logs_created_at;
DROP INDEX IF EXISTS idx_access_logs_user_id;

-- 删除访问日志表
DROP TABLE IF EXISTS access_logs;
