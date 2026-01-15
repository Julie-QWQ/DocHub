package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 请求耗时
		latency := time.Since(start)
		requestID := GetRequestID(c)

		// 构建日志字段
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		// 根据状态码选择日志级别
		if c.Writer.Status() >= 500 {
			logger.Error("HTTP请求", fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("HTTP请求", fields...)
		} else {
			logger.Info("HTTP请求", fields...)
		}
	}
}
