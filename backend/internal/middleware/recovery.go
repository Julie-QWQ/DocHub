package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/pkg/logger"
	"github.com/study-upc/backend/internal/pkg/response"
	"go.uber.org/zap"
)

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				logger.Error("服务器内部错误",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回统一错误响应
				response.ServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ErrorHandler 错误处理中间件（处理 c.Error）
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Error("请求处理错误",
				zap.Error(err.Err),
				zap.String("path", c.Request.URL.Path),
			)

			// 如果还没有响应，则返回错误
			if !c.Writer.Written() {
				response.ServerError(c, err.Error())
			}
		}
	}
}
