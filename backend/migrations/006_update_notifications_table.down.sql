-- 回滚 notifications 表结构更新
-- 版本: 006

-- 删除触发器
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;

-- 删除索引
DROP INDEX IF EXISTS idx_notifications_status;
DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_deleted_at;

-- 删除新增的字段
ALTER TABLE notifications
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS link,
    DROP COLUMN IF EXISTS read_at,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS deleted_at;
