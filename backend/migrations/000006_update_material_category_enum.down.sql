-- 回滚: 删除新增的 enum 值
-- 注意: PostgreSQL 不支持删除 enum 值,这是无法完全回滚的操作
-- 如果需要回滚,需要:
-- 1. 创建新的 enum 类型(不包含这些值)
-- 2. 将 materials 表的 category 列改为新类型
-- 3. 删除旧类型

-- 这里只记录无法回滚的说明
-- 无法直接删除 enum 值: 'textbook', 'note', 'exercise', 'thesis', 'exam_paper'
