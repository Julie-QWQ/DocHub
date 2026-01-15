-- 恢复学委申请表的 user_id 唯一约束
-- 版本: 013
-- 描述: 回滚到只允许一个用户有一个申请记录

-- 删除复合索引
DROP INDEX IF EXISTS idx_committee_applications_user_id_status;

-- 重新创建唯一约束
CREATE UNIQUE INDEX committee_applications_user_id_key ON committee_applications(user_id);
