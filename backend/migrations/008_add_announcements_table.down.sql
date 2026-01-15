-- 删除触发器
DROP TRIGGER IF EXISTS trigger_update_announcements_updated_at ON announcements;

-- 删除函数
DROP FUNCTION IF EXISTS update_announcements_updated_at();

-- 删除表
DROP TABLE IF EXISTS announcements;
