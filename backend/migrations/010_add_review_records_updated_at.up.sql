-- 为 review_records 表添加 updated_at 和 deleted_at 字段
-- 以匹配 GORM 模型定义

ALTER TABLE review_records
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 为 deleted_at 创建索引 (用于软删除)
CREATE INDEX IF NOT EXISTS idx_review_records_deleted_at ON review_records(deleted_at);
