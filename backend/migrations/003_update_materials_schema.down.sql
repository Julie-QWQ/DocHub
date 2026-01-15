-- Study-UPC 资料表结构回滚
-- 版本: 003
-- 描述: 回滚资料表结构更新

-- 删除触发器
DROP TRIGGER IF EXISTS materials_search_vector_trigger ON materials;

-- 删除触发器函数
DROP FUNCTION IF EXISTS materials_search_vector_update();

-- 删除索引
DROP INDEX IF EXISTS idx_materials_search;
DROP INDEX IF EXISTS idx_materials_file_key;
DROP INDEX IF EXISTS idx_materials_course_name;
DROP INDEX IF EXISTS idx_materials_reviewer_id;
DROP INDEX IF EXISTS idx_reports_handler_id;
DROP INDEX IF EXISTS idx_reports_user_id;
DROP INDEX IF EXISTS idx_favorites_updated_at;
DROP INDEX IF EXISTS idx_download_records_updated_at;

-- 删除时间戳触发器
DROP TRIGGER IF EXISTS update_favorites_updated_at ON favorites;
DROP TRIGGER IF EXISTS update_download_records_updated_at ON download_records;

-- 回滚资料表结构
ALTER TABLE materials
    DROP COLUMN IF EXISTS course_name,
    DROP COLUMN IF EXISTS file_key,
    DROP COLUMN IF EXISTS mime_type,
    DROP COLUMN IF EXISTS favorite_count,
    DROP COLUMN IF EXISTS reviewer_id,
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS rejection_reason,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS search_vector;

-- 回滚举报表结构
ALTER TABLE reports
    RENAME COLUMN user_id TO reporter_id,
    DROP COLUMN IF EXISTS handler_id,
    DROP COLUMN IF EXISTS handled_at,
    DROP COLUMN IF EXISTS handle_note;

-- 回滚收藏表结构
ALTER TABLE favorites
    DROP COLUMN IF EXISTS deleted_at;

-- 回滚下载记录表结构
ALTER TABLE download_records
    DROP COLUMN IF EXISTS deleted_at;

-- 删除枚举类型
DROP TYPE IF EXISTS report_status CASCADE;
DROP TYPE IF EXISTS material_status CASCADE;
DROP TYPE IF EXISTS material_category CASCADE;
