-- 同步 materials 表结构到 Go 模型
-- 版本: 009
-- 描述: 修复 materials 表结构,添加缺失字段,处理 file_url 迁移
-- 注意: 此迁移只会执行一次,后续会被跳过

-- 只在枚举类型不存在时创建
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_category') THEN
        CREATE TYPE material_category AS ENUM (
            'courseware',    -- 课件
            'exam',          -- 试卷
            'experiment',    -- 实验
            'exercise',      -- 习题
            'reference',     -- 参考资料
            'other'          -- 其他
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_status') THEN
        CREATE TYPE material_status AS ENUM (
            'pending',       -- 待审核
            'approved',      -- 已通过
            'rejected',      -- 已拒绝
            'deleted'        -- 已删除
        );
    END IF;
END $$;

-- 添加缺失的列(使用 IF NOT EXISTS 保证幂等性)
ALTER TABLE materials
    ADD COLUMN IF NOT EXISTS category material_category NOT NULL DEFAULT 'other',
    ADD COLUMN IF NOT EXISTS course_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS file_key VARCHAR(500) UNIQUE,
    ADD COLUMN IF NOT EXISTS mime_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS favorite_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reviewer_id BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 如果 file_url 列存在且 file_key 不存在,迁移数据
DO $$
BEGIN
    -- 检查是否有 file_url 列
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'materials' AND column_name = 'file_url'
    ) THEN
        -- 将 file_url 复制到 file_key (去掉域名部分)
        UPDATE materials
        SET file_key = SUBSTRING(file_url FROM '.*/([^/]+)$')
        WHERE file_key IS NULL OR file_key = '';

        -- 删除 file_url 列
        ALTER TABLE materials DROP COLUMN IF EXISTS file_url;
    END IF;
END $$;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_materials_category ON materials(category);
CREATE INDEX IF NOT EXISTS idx_materials_status ON materials(status);
CREATE INDEX IF NOT EXISTS idx_materials_course_name ON materials(course_name);
CREATE INDEX IF NOT EXISTS idx_materials_file_key ON materials(file_key);
CREATE INDEX IF NOT EXISTS idx_materials_reviewer_id ON materials(reviewer_id);
DROP INDEX IF EXISTS idx_materials_title;

-- 创建全文搜索索引
DROP INDEX IF EXISTS idx_materials_search;
CREATE INDEX IF NOT EXISTS idx_materials_search ON materials USING gin(search_vector);

-- 创建触发器自动更新搜索向量
CREATE OR REPLACE FUNCTION materials_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('simple', COALESCE(NEW.course_name, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为 materials 表添加搜索向量更新触发器
DROP TRIGGER IF EXISTS materials_search_vector_trigger ON materials;
CREATE TRIGGER materials_search_vector_trigger
    BEFORE INSERT OR UPDATE ON materials
    FOR EACH ROW
    EXECUTE FUNCTION materials_search_vector_update();

-- 更新现有数据的搜索向量
UPDATE materials SET search_vector =
    setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(description, '')), 'B') ||
    setweight(to_tsvector('simple', COALESCE(course_name, '')), 'C')
WHERE search_vector IS NULL;

-- 添加注释
COMMENT ON COLUMN materials.category IS '资料分类';
COMMENT ON COLUMN materials.course_name IS '课程名称';
COMMENT ON COLUMN materials.status IS '审核状态';
COMMENT ON COLUMN materials.file_key IS 'OSS 存储键';
COMMENT ON COLUMN materials.mime_type IS 'MIME 类型';
COMMENT ON COLUMN materials.favorite_count IS '收藏次数';
COMMENT ON COLUMN materials.reviewer_id IS '审核人ID';
COMMENT ON COLUMN materials.reviewed_at IS '审核时间';
COMMENT ON COLUMN materials.rejection_reason IS '拒绝原因';
