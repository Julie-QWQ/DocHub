-- 回滚: 移除 review_records 表的 updated_at 和 deleted_at 字段

DROP INDEX IF EXISTS idx_review_records_deleted_at;

ALTER TABLE review_records
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS updated_at;
