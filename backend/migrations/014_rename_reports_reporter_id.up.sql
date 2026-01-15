-- 重命名举报表的 reporter_id 字段为 user_id
-- 版本: 014
-- 描述: 统一字段命名，将 reporter_id 改为 user_id

-- 重命名字段
ALTER TABLE reports RENAME COLUMN reporter_id TO user_id;

-- 重命名索引以保持一致性
DROP INDEX IF EXISTS idx_reports_reporter_id;
CREATE INDEX idx_reports_user_id ON reports(user_id);

-- 添加 handled_at 字段（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'reports' AND column_name = 'handled_at'
    ) THEN
        ALTER TABLE reports ADD COLUMN handled_at TIMESTAMP;
    END IF;
END $$;
