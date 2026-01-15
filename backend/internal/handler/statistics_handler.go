package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	statsService service.StatisticsService
}

// NewStatisticsHandler 创建统计处理器
func NewStatisticsHandler(statsService service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statsService: statsService,
	}
}

// GetOverviewStatistics 获取概览统计
// @Summary 获取概览统计
// @Description 获取系统的概览统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.OverviewStatistics}
// @Router /api/v1/statistics/overview [get]
func (h *StatisticsHandler) GetOverviewStatistics(c *gin.Context) {
	stats, err := h.statsService.GetOverviewStatistics()
	if err != nil {
		response.Error(c, response.CodeServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetUserStatistics 获取用户统计
// @Summary 获取用户统计
// @Description 获取用户相关的统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.UserStatistics}
// @Router /api/v1/statistics/users [get]
func (h *StatisticsHandler) GetUserStatistics(c *gin.Context) {
	stats, err := h.statsService.GetUserStatistics()
	if err != nil {
		response.Error(c, response.CodeServerError, "获取用户统计失败")
		return
	}

	response.Success(c, stats)
}

// GetUserTrend 获取用户趋势
// @Summary 获取用户趋势
// @Description 获取用户增长趋势数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.Response{data=[]model.TrendData}
// @Router /api/v1/statistics/users/trend [get]
func (h *StatisticsHandler) GetUserTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	trend, err := h.statsService.GetUserTrend(days)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取用户趋势失败")
		return
	}

	response.Success(c, trend)
}

// GetMaterialStatistics 获取资料统计
// @Summary 获取资料统计
// @Description 获取资料相关的统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.MaterialStatistics}
// @Router /api/v1/statistics/materials [get]
func (h *StatisticsHandler) GetMaterialStatistics(c *gin.Context) {
	stats, err := h.statsService.GetMaterialStatistics()
	if err != nil {
		response.Error(c, response.CodeServerError, "获取资料统计失败")
		return
	}

	response.Success(c, stats)
}

// GetMaterialTrend 获取资料趋势
// @Summary 获取资料趋势
// @Description 获取资料上传趋势数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.Response{data=[]model.TrendData}
// @Router /api/v1/statistics/materials/trend [get]
func (h *StatisticsHandler) GetMaterialTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	trend, err := h.statsService.GetMaterialTrend(days)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取资料趋势失败")
		return
	}

	response.Success(c, trend)
}

// GetDownloadStatistics 获取下载统计
// @Summary 获取下载统计
// @Description 获取下载相关的统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.DownloadStatistics}
// @Router /api/v1/statistics/downloads [get]
func (h *StatisticsHandler) GetDownloadStatistics(c *gin.Context) {
	stats, err := h.statsService.GetDownloadStatistics()
	if err != nil {
		response.Error(c, response.CodeServerError, "获取下载统计失败")
		return
	}

	response.Success(c, stats)
}

// GetDownloadTrend 获取下载趋势
// @Summary 获取下载趋势
// @Description 获取下载趋势数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.Response{data=[]model.TrendData}
// @Router /api/v1/statistics/downloads/trend [get]
func (h *StatisticsHandler) GetDownloadTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	trend, err := h.statsService.GetDownloadTrend(days)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取下载趋势失败")
		return
	}

	response.Success(c, trend)
}

// GetApplicationStatistics 获取学委申请统计
// @Summary 获取学委申请统计
// @Description 获取学委申请相关的统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.ApplicationStatistics}
// @Router /api/v1/statistics/applications [get]
func (h *StatisticsHandler) GetApplicationStatistics(c *gin.Context) {
	stats, err := h.statsService.GetApplicationStatistics()
	if err != nil {
		response.Error(c, response.CodeServerError, "获取申请统计失败")
		return
	}

	response.Success(c, stats)
}

// GetVisitStatistics 获取访问统计
// @Summary 获取访问统计
// @Description 获取网站访问相关的统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.VisitStatistics}
// @Router /api/v1/statistics/visits [get]
func (h *StatisticsHandler) GetVisitStatistics(c *gin.Context) {
	stats, err := h.statsService.GetVisitStatistics()
	if err != nil {
		response.Error(c, response.CodeServerError, "获取访问统计失败")
		return
	}

	response.Success(c, stats)
}

// GetVisitTrend 获取访问趋势
// @Summary 获取访问趋势
// @Description 获取网站访问趋势数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "天数" default(30)
// @Success 200 {object} response.Response{data=[]model.TrendData}
// @Router /api/v1/statistics/visits/trend [get]
func (h *StatisticsHandler) GetVisitTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	trend, err := h.statsService.GetVisitTrend(days)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取访问趋势失败")
		return
	}

	response.Success(c, trend)
}

// RecordPageView 记录页面浏览
// @Summary 记录页面浏览
// @Description 记录用户访问的页面（前端路由变化时调用）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.PageViewRequest true "页面浏览信息"
// @Success 200 {object} response.Response
// @Router /api/v1/statistics/page-view [post]
func (h *StatisticsHandler) RecordPageView(c *gin.Context) {
	var req model.PageViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized, "未授权")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.Error(c, response.ErrInternal, "用户ID类型错误")
		return
	}

	// 异步记录，不阻塞请求
	go func() {
		_ = h.statsService.RecordAccess(
			&userIDUint,
			c.ClientIP(),
			req.Path,
			"GET",
			c.Request.UserAgent(),
			req.Referer,
		)
	}()

	response.Success(c, nil)
}
