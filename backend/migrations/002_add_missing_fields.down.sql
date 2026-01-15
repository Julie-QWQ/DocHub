-- 回滚: 删除添加的字段
ALTER TABLE users
    DROP COLUMN IF EXISTS real_name,
    DROP COLUMN IF EXISTS avatar,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS major,
    DROP COLUMN IF EXISTS class,
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS deleted_at;
