-- 回滚: 移除 favorites 和 download_records 表的 updated_at 字段
-- 版本: 012

-- 删除 favorites 表触发器
DROP TRIGGER IF EXISTS update_favorites_updated_at ON favorites;

-- 删除 download_records 表触发器
DROP TRIGGER IF EXISTS update_download_records_updated_at ON download_records;

-- 删除 favorites 表的 updated_at 字段
ALTER TABLE favorites
DROP COLUMN IF EXISTS updated_at;

-- 删除 download_records 表的 updated_at 字段
ALTER TABLE download_records
DROP COLUMN IF EXISTS updated_at;
