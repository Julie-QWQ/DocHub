package handler

import (
	"strconv"

	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchService         service.SearchService
	recommendationService service.RecommendationService
}

// NewSearchHandler 创建搜索处理器实例
func NewSearchHandler(
	searchService service.SearchService,
	recommendationService service.RecommendationService,
) *SearchHandler {
	return &SearchHandler{
		searchService:         searchService,
		recommendationService: recommendationService,
	}
}

// Search 搜索资料
// @Summary 搜索资料
// @Description 根据关键词、分类、标签等条件搜索资料
// @Tags 搜索与推荐
// @Produce json
// @Security Bearer
// @Param keyword query string false "搜索关键词"
// @Param category query string false "分类"
// @Param course_name query string false "课程名称"
// @Param tags query []string false "标签"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方向"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=model.SearchResponse}
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	var req model.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 获取当前用户ID（可选）
	userID, _ := middleware.GetUserID(c)

	// 执行搜索
	result, err := h.searchService.Search(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, result)
}

// GetHotKeywords 获取热门搜索词
// @Summary 热门搜索词
// @Description 获取热门搜索词列表
// @Tags 搜索与推荐
// @Produce json
// @Param limit query int false "返回数量"
// @Success 200 {object} response.Response{data=[]model.HotKeyword}
// @Router /api/v1/search/hot-keywords [get]
func (h *SearchHandler) GetHotKeywords(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	keywords, err := h.searchService.GetHotKeywords(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, keywords)
}

// GetSearchHistory 获取搜索历史
// @Summary 搜索历史
// @Description 获取当前用户的搜索历史
// @Tags 搜索与推荐
// @Produce json
// @Security Bearer
// @Param limit query int false "返回数量"
// @Success 200 {object} response.Response{data=[]model.SearchHistory}
// @Router /api/v1/search/history [get]
func (h *SearchHandler) GetSearchHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	histories, err := h.searchService.GetUserSearchHistories(c.Request.Context(), userID, limit)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, histories)
}

// ClearSearchHistory 清空搜索历史
// @Summary 清空搜索历史
// @Description 清空当前用户的搜索历史
// @Tags 搜索与推荐
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/search/history [delete]
func (h *SearchHandler) ClearSearchHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.searchService.ClearSearchHistory(c.Request.Context(), userID); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetHotMaterials 获取热门资料
// @Summary 热门资料
// @Description 获取热门资料列表
// @Tags 搜索与推荐
// @Produce json
// @Param limit query int false "返回数量"
// @Success 200 {object} response.Response{data=[]model.Material}
// @Router /api/v1/materials/hot [get]
func (h *SearchHandler) GetHotMaterials(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	materials, err := h.recommendationService.GetHotMaterials(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, materials)
}

// GetRecommendations 获取推荐资料
// @Summary 推荐资料
// @Description 获取推荐资料列表
// @Tags 搜索与推荐
// @Produce json
// @Security Bearer
// @Param type query string false "推荐类型"
// @Param material_id query int false "资料ID"
// @Param limit query int false "返回数量"
// @Success 200 {object} response.Response{data=[]model.RecommendationResult}
// @Router /api/v1/materials/recommend [get]
func (h *SearchHandler) GetRecommendations(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	var req model.RecommendationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	results, err := h.recommendationService.GetRecommendations(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, results)
}
