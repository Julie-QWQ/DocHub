-- 将 material_category 从 ENUM 转换为 VARCHAR(50)
-- 这样管理员就可以动态创建任意类型的资料分类

-- 1. 将枚举类型转换为 VARCHAR(50)
ALTER TABLE materials ALTER COLUMN category TYPE VARCHAR(50) USING category::text;

-- 2. 添加外键约束,确保 category 必须存在于 material_categories 表中
-- 这样既保留了数据完整性,又允许动态添加新类型
ALTER TABLE materials
  ADD CONSTRAINT materials_category_fkey
  FOREIGN KEY (category)
  REFERENCES material_categories(code)
  ON DELETE RESTRICT
  ON UPDATE CASCADE;

-- 3. 添加检查约束,确保 category 不为空
ALTER TABLE materials
  ADD CONSTRAINT materials_category_check
  CHECK (category IS NOT NULL AND category != '');

-- 4. 添加注释
COMMENT ON COLUMN materials.category IS '资料分类(外键关联到 material_category_configs.code,支持动态扩展)';

-- 5. (可选) 删除旧的枚举类型 - 谨慎操作!
-- DROP TYPE IF EXISTS material_category;
