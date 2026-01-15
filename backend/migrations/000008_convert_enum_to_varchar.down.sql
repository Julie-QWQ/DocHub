-- 回滚: 将 VARCHAR(50) 转回 ENUM
-- 注意: 这只能回滚到已有的枚举值,如果有新值会失败

-- 1. 删除外键约束
ALTER TABLE materials DROP CONSTRAINT IF EXISTS materials_category_fkey;

-- 2. 删除检查约束
ALTER TABLE materials DROP CONSTRAINT IF EXISTS materials_category_check;

-- 3. 转回枚举类型(如果值都在枚举中)
ALTER TABLE materials ALTER COLUMN category TYPE material_category USING category::material_category;
