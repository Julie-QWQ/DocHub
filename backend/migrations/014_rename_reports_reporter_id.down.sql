-- 回滚举报表字段重命名
-- 版本: 014
-- 描述: 将 user_id 改回 reporter_id

-- 删除 handled_at 字段
ALTER TABLE reports DROP COLUMN IF EXISTS handled_at;

-- 删除索引
DROP INDEX IF EXISTS idx_reports_user_id;

-- 恢复字段原名
ALTER TABLE reports RENAME COLUMN user_id TO reporter_id;

-- 恢复旧索引
CREATE INDEX IF NOT EXISTS idx_reports_reporter_id ON reports(reporter_id);
