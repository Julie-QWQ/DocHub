-- 删除系统配置表触发器
DROP TRIGGER IF EXISTS update_system_configs_updated_at ON system_configs;

-- 删除系统配置表
DROP TABLE IF EXISTS system_configs;

-- 删除访问日志表
DROP TABLE IF EXISTS access_logs;

-- 删除触发器函数 (如果没有其他表使用)
DROP FUNCTION IF EXISTS update_updated_at_column();
