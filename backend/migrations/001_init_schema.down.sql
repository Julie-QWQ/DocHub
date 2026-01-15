-- Study-UPC 数据库回滚脚本
-- 版本: 001
-- 描述: 删除所有表和类型

-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_materials_updated_at ON materials;
DROP TRIGGER IF EXISTS update_committee_applications_updated_at ON committee_applications;
DROP TRIGGER IF EXISTS update_reports_updated_at ON reports;

-- 删除触发器函数
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 删除表（注意顺序，先删除有外键约束的表）
DROP TABLE IF EXISTS reports CASCADE;
DROP TABLE IF EXISTS download_records CASCADE;
DROP TABLE IF EXISTS favorites CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS review_records CASCADE;
DROP TABLE IF EXISTS committee_applications CASCADE;
DROP TABLE IF EXISTS materials CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- 删除枚举类型
DROP TYPE IF EXISTS material_status;
DROP TYPE IF EXISTS material_category;
