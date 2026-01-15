-- 更新 notifications 表结构
-- 版本: 006
-- 描述: 添加缺失的字段以匹配 Notification 模型

-- 添加缺失的字段
ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'unread',
    ADD COLUMN IF NOT EXISTS link VARCHAR(255),
    ADD COLUMN IF NOT EXISTS read_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);

-- 更新时间戳触发器
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 迁移现有数据：将 is_read 转换为 status
UPDATE notifications
SET status = CASE
    WHEN is_read = TRUE THEN 'read'
    ELSE 'unread'
END;

-- 删除旧的 is_read 字段（可选，如果确定不再需要）
-- ALTER TABLE notifications DROP COLUMN IF EXISTS is_read;
