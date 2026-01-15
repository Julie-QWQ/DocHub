-- 回滚: 恢复 materials 表旧结构

-- 添加回 file_url 列
ALTER TABLE materials
    ADD COLUMN IF NOT EXISTS file_url VARCHAR(500);

-- 迁移数据:从 file_key 生成 file_url
UPDATE materials
SET file_url = 'https://study-upc.oss-cn-shanghai.aliyuncs.com/' || file_key
WHERE file_url IS NULL;

-- 删除新增的列
ALTER TABLE materials
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS course_name,
    DROP COLUMN IF EXISTS file_key,
    DROP COLUMN IF EXISTS mime_type,
    DROP COLUMN IF EXISTS favorite_count,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS reviewer_id,
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS rejection_reason,
    DROP COLUMN IF EXISTS deleted_at;

-- 删除枚举类型
DROP TYPE IF EXISTS material_category CASCADE;
DROP TYPE IF EXISTS material_status CASCADE;
