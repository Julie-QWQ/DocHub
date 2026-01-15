-- 修复 review_records 表结构,使其与 GORM 模型完全匹配

-- 删除旧的外键约束(因为我们需要改变表结构)
ALTER TABLE review_records DROP CONSTRAINT IF EXISTS review_records_material_id_fkey;

-- 删除旧的索引
DROP INDEX IF EXISTS idx_review_records_material_id;

-- 修改 material_id 为 target_id (更通用的审核目标)
ALTER TABLE review_records RENAME COLUMN material_id TO target_id;

-- 添加 target_type 字段(审核目标类型)
ALTER TABLE review_records
    ADD COLUMN IF NOT EXISTS target_type VARCHAR(20);

-- 添加 original_data 字段(原始数据快照)
ALTER TABLE review_records
    ADD COLUMN IF NOT EXISTS original_data JSONB;

-- 添加 updated_at 字段(如果不存在)
ALTER TABLE review_records
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- 添加 deleted_at 字段(如果不存在,用于软删除)
ALTER TABLE review_records
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 创建新的索引
CREATE INDEX IF NOT EXISTS idx_review_records_target_id ON review_records(target_id);
CREATE INDEX IF NOT EXISTS idx_review_records_target_type ON review_records(target_type);
CREATE INDEX IF NOT EXISTS idx_review_records_deleted_at ON review_records(deleted_at);

-- 添加注释
COMMENT ON COLUMN review_records.target_id IS '审核目标ID';
COMMENT ON COLUMN review_records.target_type IS '审核目标类型(material/committee/report)';
COMMENT ON COLUMN review_records.original_data IS '原始数据快照';
