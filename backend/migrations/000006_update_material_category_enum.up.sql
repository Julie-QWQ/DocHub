-- 更新 material_category enum 类型,添加新的分类
-- 注意: PostgreSQL 不支持在 enum 中间插入值,新值会被添加到末尾
-- 另外,PostgreSQL 不支持重命名 enum 值,所以 'exam' 保留,我们添加 'exam_paper'

-- 添加新的 enum 值 (按字母顺序添加,PostgreSQL 会将它们添加到列表末尾)
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'textbook';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'note';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'exercise';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'thesis';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'exam_paper';

-- 添加注释说明
COMMENT ON TYPE material_category IS '资料分类枚举: courseware, textbook, experiment, exam, exam_paper, note, exercise, reference, thesis, other';
