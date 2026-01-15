-- 添加 reports 表的 deleted_at 字段
-- 版本: 008
-- 描述: 修复 reports 表缺少 deleted_at 字段的问题

ALTER TABLE reports
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
