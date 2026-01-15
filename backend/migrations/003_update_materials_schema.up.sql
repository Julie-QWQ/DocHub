-- Study-UPC 资料表结构更新
-- 版本: 003
-- 描述: 更新资料表结构以支持完整的资料管理功能
-- 注意: 此迁移只会执行一次,后续会被跳过

-- 只在枚举类型不存在时创建(避免删除已有数据)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_category') THEN
        CREATE TYPE material_category AS ENUM (
            'textbook',      -- 教材
            'reference',     -- 参考书
            'exam_paper',    -- 试卷
            'note',          -- 笔记
            'exercise',      -- 习题
            'experiment',    -- 实验指导
            'thesis',        -- 论文
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

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'report_status') THEN
        CREATE TYPE report_status AS ENUM (
            'pending',       -- 待处理
            'approved',      -- 已通过
            'rejected'       -- 已驳回
        );
    END IF;
END $$;

-- 更新资料表结构
ALTER TABLE materials
    -- 修改分类列类型
    ALTER COLUMN category TYPE material_category USING category::text::material_category,
    -- 修改状态列类型
    ALTER COLUMN status TYPE material_status USING status::text::material_status,
    -- 添加新列
    ADD COLUMN IF NOT EXISTS course_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS file_key VARCHAR(500) UNIQUE,
    ADD COLUMN IF NOT EXISTS mime_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS favorite_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reviewer_id BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 为 file_key 创建索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_materials_file_key ON materials(file_key);

-- 为 course_name 创建索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_materials_course_name ON materials(course_name);

-- 为 reviewer_id 创建索引（如果不存在）
CREATE INDEX IF NOT EXISTS idx_materials_reviewer_id ON materials(reviewer_id);

-- 添加全文搜索向量列（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'materials' AND column_name = 'search_vector'
    ) THEN
        ALTER TABLE materials ADD COLUMN search_vector tsvector;
    END IF;
END $$;

-- 创建全文搜索索引
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

-- 更新举报表结构
ALTER TABLE reports
    -- 重命名列以匹配新模型
    RENAME COLUMN reporter_id TO user_id,
    -- 修改状态列类型
    ALTER COLUMN status TYPE report_status USING status::text::report_status,
    -- 添加新列
    ADD COLUMN IF NOT EXISTS handler_id BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS handled_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS handle_note TEXT;

-- 更新索引
DROP INDEX IF EXISTS idx_reports_reporter_id;
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);
CREATE INDEX IF NOT EXISTS idx_reports_handler_id ON reports(handler_id);

-- 更新收藏表结构（添加 deleted_at 以支持软删除）
ALTER TABLE favorites
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 更新下载记录表结构（添加 deleted_at 以支持软删除）
ALTER TABLE download_records
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 更新时间戳触发器
CREATE TRIGGER update_favorites_updated_at BEFORE UPDATE ON favorites
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_download_records_updated_at BEFORE UPDATE ON download_records
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE materials IS '学习资料表';
COMMENT ON COLUMN materials.title IS '资料标题';
COMMENT ON COLUMN materials.description IS '资料描述';
COMMENT ON COLUMN materials.category IS '资料分类';
COMMENT ON COLUMN materials.course_name IS '课程名称';
COMMENT ON COLUMN materials.status IS '审核状态';
COMMENT ON COLUMN materials.file_name IS '原始文件名';
COMMENT ON COLUMN materials.file_size IS '文件大小（字节）';
COMMENT ON COLUMN materials.file_key IS 'OSS 存储键';
COMMENT ON COLUMN materials.mime_type IS 'MIME 类型';
COMMENT ON COLUMN materials.download_count IS '下载次数';
COMMENT ON COLUMN materials.favorite_count IS '收藏次数';
COMMENT ON COLUMN materials.view_count IS '浏览次数';
COMMENT ON COLUMN materials.reviewer_id IS '审核人ID';
COMMENT ON COLUMN materials.reviewed_at IS '审核时间';
COMMENT ON COLUMN materials.rejection_reason IS '拒绝原因';

COMMENT ON TABLE favorites IS '收藏表';
COMMENT ON TABLE download_records IS '下载记录表';
COMMENT ON TABLE reports IS '举报表';
