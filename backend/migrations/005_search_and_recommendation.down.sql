-- 回滚搜索和推荐功能表结构

-- 删除触发器
DROP TRIGGER IF EXISTS update_search_histories_updated_at ON search_histories;
DROP TRIGGER IF EXISTS update_hot_keywords_updated_at ON hot_keywords;

-- 删除表
DROP TABLE IF EXISTS search_histories;
DROP TABLE IF EXISTS hot_keywords;
