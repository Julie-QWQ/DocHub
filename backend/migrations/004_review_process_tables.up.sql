-- Study-UPC 审核流程表结构更新
-- 版本: 004
-- 描述: 更新审核流程相关表结构以支持完整的审核功能

-- 创建申请状态枚举类型
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_status') THEN
        CREATE TYPE application_status AS ENUM (
            'pending',    -- 待审核
            'approved',   -- 已通过
            'rejected',   -- 已拒绝
            'cancelled'   -- 已取消
        );
    END IF;
END $$;

-- 创建通知类型枚举类型
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
        CREATE TYPE notification_type AS ENUM (
            'system',      -- 系统通知
            'material',    -- 资料审核通知
            'committee',   -- 学委申请通知
            'report'       -- 举报处理通知
        );
    END IF;
END $$;

-- 创建通知状态枚举类型
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_status') THEN
        CREATE TYPE notification_status AS ENUM (
            'unread',     -- 未读
            'read'        -- 已读
        );
    END IF;
END $$;

-- 创建审核操作枚举类型
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'review_action') THEN
        CREATE TYPE review_action AS ENUM (
            'approve',    -- 通过
            'reject'      -- 拒绝
        );
    END IF;
END $$;

-- 创建审核目标类型枚举类型
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'review_target') THEN
        CREATE TYPE review_target AS ENUM (
            'material',   -- 资料
            'committee',  -- 学委申请
            'report'      -- 举报
        );
    END IF;
END $$;

-- 更新学委申请表
ALTER TABLE committee_applications
    -- 修改状态列类型
    ALTER COLUMN status TYPE application_status USING status::text::application_status,
    -- 添加缺失的列
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP;

-- 移除旧的唯一约束，允许用户多次申请（但只能有一个待审核申请）
DROP INDEX IF EXISTS committee_applications_user_id_key;
CREATE INDEX idx_committee_applications_user_id_status ON committee_applications(user_id, status);

-- 更新审核记录表
ALTER TABLE review_records
    -- 修改操作列类型
    ALTER COLUMN action TYPE review_action USING action::text::review_action,
    -- 添加新列
    ADD COLUMN IF NOT EXISTS target_type review_target,
    ADD COLUMN IF NOT EXISTS target_id BIGINT NOT NULL,
    ADD COLUMN IF NOT EXISTS original_data JSONB,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 为审核记录表创建新的索引
CREATE INDEX IF NOT EXISTS idx_review_records_target_type ON review_records(target_type);
CREATE INDEX IF NOT EXISTS idx_review_records_target_id ON review_records(target_id);
DROP INDEX IF EXISTS idx_review_records_material_id;

-- 更新通知表
ALTER TABLE notifications
    -- 修改类型列
    ALTER COLUMN type TYPE notification_type USING type::text::notification_type,
    -- 添加新列
    ADD COLUMN IF NOT EXISTS status notification_status NOT NULL DEFAULT 'unread',
    ADD COLUMN IF NOT EXISTS link VARCHAR(255),
    ADD COLUMN IF NOT EXISTS read_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- 为通知表创建新的索引
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
DROP INDEX IF EXISTS idx_notifications_is_read;

-- 更新举报表结构
ALTER TABLE reports
    -- 修改状态列类型
    ALTER COLUMN status TYPE report_status USING status::text::report_status,
    -- 重命名列以保持一致性
    RENAME COLUMN handle_comment TO handle_note;

-- 更新时间戳触发器
CREATE TRIGGER update_committee_applications_updated_at BEFORE UPDATE ON committee_applications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_review_records_updated_at BEFORE UPDATE ON review_records
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE committee_applications IS '学委申请表';
COMMENT ON COLUMN committee_applications.status IS '申请状态';
COMMENT ON COLUMN committee_applications.reason IS '申请理由';
COMMENT ON COLUMN committee_applications.reviewer_id IS '审核人ID';
COMMENT ON COLUMN committee_applications.review_comment IS '审核意见';
COMMENT ON COLUMN committee_applications.reviewed_at IS '审核时间';

COMMENT ON TABLE review_records IS '审核记录表';
COMMENT ON COLUMN review_records.target_type IS '审核目标类型';
COMMENT ON COLUMN review_records.target_id IS '目标ID';
COMMENT ON COLUMN review_records.action IS '审核操作';
COMMENT ON COLUMN review_records.comment IS '审核意见';
COMMENT ON COLUMN review_records.original_data IS '原始数据快照';

COMMENT ON TABLE notifications IS '通知表';
COMMENT ON COLUMN notifications.type IS '通知类型';
COMMENT ON COLUMN notifications.status IS '通知状态';
COMMENT ON COLUMN notifications.title IS '通知标题';
COMMENT ON COLUMN notifications.content IS '通知内容';
COMMENT ON COLUMN notifications.link IS '相关链接';
COMMENT ON COLUMN notifications.read_at IS '阅读时间';

COMMENT ON TABLE reports IS '举报表';
COMMENT ON COLUMN reports.status IS '处理状态';
COMMENT ON COLUMN reports.handler_id IS '处理人ID';
COMMENT ON COLUMN reports.handled_at IS '处理时间';
COMMENT ON COLUMN reports.handle_note IS '处理备注';
