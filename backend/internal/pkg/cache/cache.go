package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Cache 缓存接口
type Cache interface {
	// Set 设置缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get 获取缓存
	Get(ctx context.Context, key string, dest interface{}) error
	// Del 删除缓存
	Del(ctx context.Context, keys ...string) error
	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		logger: nil, // 可选：注入logger
	}
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value failed: %w", err)
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("set cache failed: %w", err)
	}

	return nil
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("get cache failed: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("unmarshal cache value failed: %w", err)
	}

	return nil
}

// Del 删除缓存
func (c *RedisCache) Del(ctx context.Context, keys ...string) error {
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("delete cache failed: %w", err)
	}
	return nil
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check cache exists failed: %w", err)
	}
	return count > 0, nil
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := c.client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("set cache expiration failed: %w", err)
	}
	return nil
}

// CacheKey 缓存键生成器
type CacheKey struct {
	prefix string
}

// NewCacheKey 创建缓存键生成器
func NewCacheKey(prefix string) *CacheKey {
	return &CacheKey{prefix: prefix}
}

// Build 构建缓存键
func (ck *CacheKey) Build(parts ...string) string {
	key := ck.prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}

// 常用缓存键前缀
const (
	// 用户相关
	UserCachePrefix     = "user"
	UserInfoCacheKey    = "info"
	UserStatsCacheKey   = "stats"

	// 资料相关
	MaterialCachePrefix     = "material"
	MaterialDetailCacheKey  = "detail"
	MaterialListCacheKey    = "list"
	MaterialHotCacheKey     = "hot"

	// 通知相关
	NotificationCachePrefix = "notification"
	NotificationUnreadKey   = "unread"

	// 统计相关
	StatsCachePrefix    = "stats"
	DailyStatsKey       = "daily"
	PendingCountKey     = "pending"
)

// 缓存过期时间常量
const (
	CacheTTLShort  = 5 * time.Minute   // 短期缓存：5分钟
	CacheTTLMedium = 10 * time.Minute  // 中期缓存：10分钟
	CacheTTLLong   = 30 * time.Minute  // 长期缓存：30分钟
	CacheTTLDaily  = 24 * time.Hour    // 日缓存：24小时
)

// GetOrSet 获取或设置缓存（缓存穿透保护）
func GetOrSet(ctx context.Context, cache Cache, key string, dest interface{}, expiration time.Duration, fn func() (interface{}, error)) error {
	// 尝试从缓存获取
	err := cache.Get(ctx, key, dest)
	if err == nil {
		return nil // 缓存命中
	}

	// 缓存未命中，调用函数获取数据
	data, err := fn()
	if err != nil {
		return fmt.Errorf("get data failed: %w", err)
	}

	// 设置缓存
	if err := cache.Set(ctx, key, data, expiration); err != nil {
		// 记录错误但不影响业务
		// logger.Warn("set cache failed", zap.Error(err))
	}

	// 设置返回值
	if err := cache.Get(ctx, key, dest); err != nil {
		return fmt.Errorf("get cached data failed: %w", err)
	}

	return nil
}
