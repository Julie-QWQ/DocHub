-- 回滚 review_records 表结构修改

-- 删除索引
DROP INDEX IF EXISTS idx_review_records_deleted_at;
DROP INDEX IF EXISTS idx_review_records_target_type;
DROP INDEX IF EXISTS idx_review_records_target_id;

-- 删除添加的字段
ALTER TABLE review_records DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE review_records DROP COLUMN IF EXISTS updated_at;
ALTER TABLE review_records DROP COLUMN IF EXISTS original_data;
ALTER TABLE review_records DROP COLUMN IF EXISTS target_type;

-- 恢复 material_id 列
ALTER TABLE review_records RENAME COLUMN target_id TO material_id;

-- 重新创建外键约束
ALTER TABLE review_records
    ADD CONSTRAINT review_records_material_id_fkey
    FOREIGN KEY (material_id) REFERENCES materials(id);

-- 重新创建索引
CREATE INDEX IF NOT EXISTS idx_review_records_material_id ON review_records(material_id);
