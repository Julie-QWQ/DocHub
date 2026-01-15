-- 创建资料类型配置表
CREATE TABLE IF NOT EXISTS material_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,         -- 类型代码(如 'courseware', 'exam_paper')
    name VARCHAR(100) NOT NULL,               -- 类型名称(如 '课件', '试卷')
    description TEXT,                         -- 描述
    icon VARCHAR(100),                        -- 图标
    sort_order INTEGER NOT NULL DEFAULT 0,    -- 排序
    is_active BOOLEAN NOT NULL DEFAULT true,  -- 是否启用
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL         -- 软删除标记
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_material_categories_code ON material_categories(code);
CREATE INDEX IF NOT EXISTS idx_material_categories_sort_order ON material_categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_material_categories_is_active ON material_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_material_categories_deleted_at ON material_categories(deleted_at);

-- 插入默认资料类型
INSERT INTO material_categories (code, name, description, icon, sort_order) VALUES
    ('courseware', '课件', '教学课件、PPT等', 'document', 1),
    ('textbook', '教材', '课程教材、教科书', 'book', 2),
    ('reference', '参考资料', '参考书、辅助资料', 'reference', 3),
    ('exam_paper', '试卷', '考试试卷、往年试题', 'exam', 4),
    ('exercise', '习题', '课后习题、练习题', 'exercise', 5),
    ('experiment', '实验指导', '实验报告、实验指导', 'experiment', 6),
    ('note', '笔记', '课堂笔记、复习笔记', 'note', 7),
    ('thesis', '论文', '学术论文、毕业论文', 'thesis', 8),
    ('other', '其他', '其他类型资料', 'other', 9)
ON CONFLICT (code) DO NOTHING;

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_material_categories_updated_at
    BEFORE UPDATE ON material_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE material_categories IS '资料类型配置表';
COMMENT ON COLUMN material_categories.code IS '类型代码,用于系统标识';
COMMENT ON COLUMN material_categories.name IS '类型名称,用于显示';
COMMENT ON COLUMN material_categories.description IS '类型描述';
COMMENT ON COLUMN material_categories.icon IS '图标名称';
COMMENT ON COLUMN material_categories.sort_order IS '排序,数字越小越靠前';
COMMENT ON COLUMN material_categories.is_active IS '是否启用';
