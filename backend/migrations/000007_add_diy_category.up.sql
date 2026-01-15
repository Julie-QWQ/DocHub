-- 添加 diy (自定义类型) 到 material_category enum
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'diy';

-- 更新注释说明
COMMENT ON TYPE material_category IS '资料分类枚举: courseware, textbook, experiment, exam, exam_paper, note, exercise, reference, thesis, other, diy';
