-- Study-UPC 搜索和推荐功能表结构
-- 版本: 005
-- 描述: 创建搜索历史和热门搜索词表

-- 创建搜索历史表
CREATE TABLE IF NOT EXISTS search_histories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    keyword VARCHAR(200) NOT NULL,
    result_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 为搜索历史表创建索引
CREATE INDEX IF NOT EXISTS idx_search_histories_user_id ON search_histories(user_id);
CREATE INDEX IF NOT EXISTS idx_search_histories_created_at ON search_histories(created_at DESC);

-- 创建热门搜索词表
CREATE TABLE IF NOT EXISTS hot_keywords (
    id BIGSERIAL PRIMARY KEY,
    keyword VARCHAR(200) NOT NULL UNIQUE,
    search_count INT NOT NULL DEFAULT 0,
    last_searched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 为热门搜索词表创建索引
CREATE INDEX IF NOT EXISTS idx_hot_keywords_search_count ON hot_keywords(search_count DESC);
CREATE INDEX IF NOT EXISTS idx_hot_keywords_last_searched ON hot_keywords(last_searched_at DESC);

-- 添加注释
COMMENT ON TABLE search_histories IS '搜索历史表';
COMMENT ON COLUMN search_histories.user_id IS '用户ID';
COMMENT ON COLUMN search_histories.keyword IS '搜索关键词';
COMMENT ON COLUMN search_histories.result_count IS '结果数量';

COMMENT ON TABLE hot_keywords IS '热门搜索词表';
COMMENT ON COLUMN hot_keywords.keyword IS '搜索关键词';
COMMENT ON COLUMN hot_keywords.search_count IS '搜索次数';
COMMENT ON COLUMN hot_keywords.last_searched_at IS '最后搜索时间';

-- 更新时间戳触发器
CREATE TRIGGER update_search_histories_updated_at BEFORE UPDATE ON search_histories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_hot_keywords_updated_at BEFORE UPDATE ON hot_keywords
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
