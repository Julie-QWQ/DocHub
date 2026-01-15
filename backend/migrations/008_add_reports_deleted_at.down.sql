-- 回滚: 移除 reports 表的 deleted_at 字段

ALTER TABLE reports
    DROP COLUMN IF EXISTS deleted_at;
