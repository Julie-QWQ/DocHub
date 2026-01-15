-- 移除学委申请表的 user_id 唯一约束
-- 版本: 013
-- 描述: 允许用户有多个申请记录（取消/拒绝的记录可以保留）

-- 移除唯一约束（如果存在）
ALTER TABLE committee_applications DROP CONSTRAINT IF EXISTS committee_applications_user_id_key;

-- 确保有复合索引用于查询优化
CREATE INDEX IF NOT EXISTS idx_committee_applications_user_id_status ON committee_applications(user_id, status);
