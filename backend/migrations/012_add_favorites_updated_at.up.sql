-- 添加 favorites 和 download_records 表的 updated_at 字段
-- 版本: 012
-- 描述: 修复 favorites 和 download_records 表缺少 updated_at 字段的问题
-- 注意: 此迁移针对 study_upc_dev 数据库

-- 为 favorites 表添加 updated_at 字段
ALTER TABLE favorites
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- 为 download_records 表添加 updated_at 字段
ALTER TABLE download_records
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- 注意: 触发器已经在初始化脚本中创建,无需重复创建
