package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/study-upc/backend/internal/pkg/config"
	"github.com/study-upc/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

var RDB *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.Config) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}

	logger.Info("Redis 连接成功", zap.String("addr", cfg.Redis.GetAddr()))
	return nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}
