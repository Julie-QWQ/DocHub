package middleware

import (
	"strings"

	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// JWTAuth JWT 认证中间件
func JWTAuth(jwtManager *utils.JWTManager, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, response.ErrUnauthorized, "未提供认证 Token")
			c.Abort()
			return
		}

		// Bearer Token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, response.ErrUnauthorized, "Token 格式错误")
			c.Abort()
			return
		}

		accessToken := parts[1]

		// 验证 Token
		claims, err := jwtManager.ValidateAccessToken(accessToken)
		if err != nil {
			switch err {
			case utils.ErrTokenExpired:
				response.Error(c, response.ErrInvalidToken, "Token 已过期")
			case utils.ErrTokenTypeWrong:
				response.Error(c, response.ErrInvalidToken, "Token 类型错误")
			default:
				response.Error(c, response.ErrInvalidToken, "Token 无效")
			}
			c.Abort()
			return
		}

		// 检查 Token 是否在黑名单中
		key := "auth:blacklist:" + accessToken
		exists, err := redisClient.Exists(c.Request.Context(), key).Result()
		if err != nil {
			response.Error(c, response.ErrInternal, "检查 Token 黑名单失败")
			c.Abort()
			return
		}
		if exists > 0 {
			response.Error(c, response.ErrInvalidToken, "Token 已失效")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalJWTAuth 可选的 JWT 认证中间件
// 如果提供了 Token 则验证，没有提供则不验证
func OptionalJWTAuth(jwtManager *utils.JWTManager, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		accessToken := parts[1]

		claims, err := jwtManager.ValidateAccessToken(accessToken)
		if err != nil {
			c.Next()
			return
		}

		// 检查 Token 是否在黑名单中
		key := "auth:blacklist:" + accessToken
		exists, err := redisClient.Exists(c.Request.Context(), key).Result()
		if err != nil || exists > 0 {
			c.Next()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
