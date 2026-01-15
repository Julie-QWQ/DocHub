-- 公告表
CREATE TABLE IF NOT EXISTS announcements (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(200) NOT NULL,           -- 公告标题
    content TEXT NOT NULL,                 -- 公告内容
    priority VARCHAR(20) NOT NULL DEFAULT 'normal' CHECK (priority IN ('normal', 'high')), -- 优先级
    author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- 发布者ID
    is_active BOOLEAN NOT NULL DEFAULT true,  -- 是否启用
    published_at TIMESTAMP,                -- 发布时间
    expires_at TIMESTAMP                   -- 过期时间
);

-- 索引
CREATE INDEX idx_announcements_author_id ON announcements(author_id);
CREATE INDEX idx_announcements_is_active ON announcements(is_active);
CREATE INDEX idx_announcements_published_at ON announcements(published_at DESC);
CREATE INDEX idx_announcements_priority ON announcements(priority);

-- 更新时间触发器
CREATE OR REPLACE FUNCTION update_announcements_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_announcements_updated_at
    BEFORE UPDATE ON announcements
    FOR EACH ROW
    EXECUTE FUNCTION update_announcements_updated_at();

-- 插入示例数据
INSERT INTO announcements (title, content, priority, author_id, is_active, published_at) VALUES
    ('欢迎使用 Study-UPC', '学院学习资料托管平台，支持资料上传、下载、收藏等功能。', 'high', 1, true, CURRENT_TIMESTAMP),
    ('资料上传须知', '上传资料前请确保文件内容合法合规，支持 PDF、DOC、PPT 等格式。', 'normal', 1, true, CURRENT_TIMESTAMP),
    ('学委申请开放', '现在可以申请成为学委，申请后需要管理员审核通过。', 'normal', 1, true, CURRENT_TIMESTAMP);
