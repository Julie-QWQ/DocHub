package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/pkg/database"
	"github.com/study-upc/backend/internal/pkg/response"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check 健康检查接口
func (h *HealthHandler) Check(c *gin.Context) {
	// 检查数据库
	ctx := context.Background()
	sqlDB, err := database.DB.DB()
	if err != nil || sqlDB == nil {
		response.Fail(c, response.CodeServerError, "数据库连接失败")
		return
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		response.Fail(c, response.CodeServerError, "数据库连接失败")
		return
	}

	// 检查 Redis
	if err := database.RDB.Ping(ctx).Err(); err != nil {
		response.Fail(c, response.CodeServerError, "Redis连接失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Liveness 存活检查（轻量级）
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
