package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/study-upc/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Window time.Duration // 时间窗口
	Limit  int           // 请求限制数
	Prefix string        // Redis key 前缀
}

// LoginRateLimit 登录限流中间件
// 防止暴力破解，限制 IP 和用户的登录尝试次数
func LoginRateLimit(redisClient *redis.Client, ipLimit, userLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 获取客户端 IP
		ip := c.ClientIP()

		// 获取用户名（如果提供了）
		var username string
		if jsonBody, exists := c.Get("json_body"); exists {
			if body, ok := jsonBody.(map[string]interface{}); ok {
				if u, ok := body["username"].(string); ok {
					username = u
				}
			}
		}

		// IP 限流：每个 IP 每小时最多尝试 ipLimit 次
		ipKey := fmt.Sprintf("login:limit:ip:%s", ip)
		ipCount, err := redisClient.Incr(ctx, ipKey).Result()
		if err != nil {
			// 记录错误但不阻塞请求
			c.Next()
			return
		}

		// 设置过期时间
		if ipCount == 1 {
			redisClient.Expire(ctx, ipKey, time.Hour)
		}

		// 检查 IP 限流
		if ipCount > int64(ipLimit) {
			ttl, _ := redisClient.TTL(ctx, ipKey).Result()
			response.Error(c, response.ErrForbidden, fmt.Sprintf("登录尝试过多，请在 %d 分钟后重试", int(ttl.Minutes())+1))
			c.Abort()
			return
		}

		// 用户名限流：每个用户名每 15 分钟最多尝试 userLimit 次
		if username != "" {
			userKey := fmt.Sprintf("login:limit:user:%s", username)
			userCount, err := redisClient.Incr(ctx, userKey).Result()
			if err != nil {
				c.Next()
				return
			}

			if userCount == 1 {
				redisClient.Expire(ctx, userKey, 15*time.Minute)
			}

			if userCount > int64(userLimit) {
				ttl, _ := redisClient.TTL(ctx, userKey).Result()
				response.Error(c, response.ErrForbidden, fmt.Sprintf("该账户登录尝试过多，请在 %d 分钟后重试", int(ttl.Minutes())+1))
				c.Abort()
				return
			}
		}

		// 添加响应头显示剩余尝试次数
		remaining := ipLimit - int(ipCount)
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(ipLimit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10))

		c.Next()
	}
}

// GeneralRateLimit 通用限流中间件
func GeneralRateLimit(redisClient *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		key := fmt.Sprintf("%s:%s", config.Prefix, c.ClientIP())

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, config.Window)
		}

		if count > int64(config.Limit) {
			response.Error(c, response.ErrForbidden, "请求过于频繁，请稍后重试")
			c.Abort()
			return
		}

		c.Next()
	}
}
