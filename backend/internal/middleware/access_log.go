package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/service"
)

// AccessLog 访问日志中间件
// 只记录页面访问，不记录 API 请求（路径以 /api/ 开头）
func AccessLog(statsService service.StatisticsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求
		c.Next()

		// 只记录成功的请求(状态码 < 400)
		if c.Writer.Status() >= 400 {
			return
		}

		path := c.Request.URL.Path

		// 过滤条件：
		// 1. 健康检查接口
		// 2. Swagger 文档
		// 3. API 请求（前端路由监听会主动记录页面浏览）
		if path == "/health" || path == "/liveness" || strings.HasPrefix(path, "/swagger/") || strings.HasPrefix(path, "/api/") {
			return
		}

		// 获取用户ID
		var userID *uint
		if uid, exists := c.Get("user_id"); exists && uid != nil {
			if id, ok := uid.(uint); ok {
				userID = &id
			}
		}

		// 只记录 GET 请求（页面浏览），忽略 POST/PUT/DELETE
		if c.Request.Method != "GET" {
			return
		}

		// 异步记录页面访问日志,不影响请求性能
		go func() {
			_ = statsService.RecordAccess(
				userID,
				c.ClientIP(),
				path,
				c.Request.Method,
				c.Request.UserAgent(),
				c.Request.Referer(),
			)
		}()
	}
}
