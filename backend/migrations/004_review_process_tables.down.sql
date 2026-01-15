-- Study-UPC 审核流程表结构回滚
-- 版本: 004
-- 描述: 回滚审核流程相关表结构更改

-- 删除触发器
DROP TRIGGER IF EXISTS update_committee_applications_updated_at ON committee_applications;
DROP TRIGGER IF EXISTS update_review_records_updated_at ON review_records;
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;

-- 回滚通知表更改
ALTER TABLE notifications
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS read_at,
    DROP COLUMN IF EXISTS link,
    DROP COLUMN IF EXISTS status,
    ALTER COLUMN type TYPE VARCHAR(50) USING type::text;

-- 删除通知表索引
DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_status;

-- 回滚审核记录表更改
ALTER TABLE review_records
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS original_data,
    DROP COLUMN IF EXISTS target_id,
    DROP COLUMN IF EXISTS target_type,
    ALTER COLUMN action TYPE VARCHAR(20) USING action::text;

-- 删除审核记录表索引
DROP INDEX IF EXISTS idx_review_records_target_type;
DROP INDEX IF EXISTS idx_review_records_target_id;

-- 回滚学委申请表更改
ALTER TABLE committee_applications
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS deleted_at,
    ALTER COLUMN status TYPE VARCHAR(20) USING status::text;

-- 恢复学委申请表唯一约束
DROP INDEX IF EXISTS idx_committee_applications_user_id_status;
CREATE UNIQUE INDEX committee_applications_user_id_key ON committee_applications(user_id);

-- 回滚举报表更改
ALTER TABLE reports
    RENAME COLUMN handle_note TO handle_comment,
    ALTER COLUMN status TYPE VARCHAR(20) USING status::text;

-- 删除枚举类型
DROP TYPE IF EXISTS review_target CASCADE;
DROP TYPE IF EXISTS review_action CASCADE;
DROP TYPE IF EXISTS notification_status CASCADE;
DROP TYPE IF EXISTS notification_type CASCADE;
DROP TYPE IF EXISTS application_status CASCADE;
