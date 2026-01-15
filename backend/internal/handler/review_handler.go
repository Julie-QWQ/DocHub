package handler

import (
	"strconv"

	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 审核处理器
type ReviewHandler struct {
	reviewService service.ReviewService
}

// NewReviewHandler 创建审核处理器实例
func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// HandleReportRequest 处理举报请求
type HandleReportRequest struct {
	Approved bool   `json:"approved" binding:"required"` // 是否通过
	Note     string `json:"note" binding:"omitempty,max=500"` // 处理备注
}

// ListReviewHistoryRequest 审核历史请求参数
type ListReviewHistoryRequest struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
	TargetType string `form:"target_type"` // 目标类型
	TargetID  string `form:"target_id"`    // 目标ID
}

// ReviewMaterial 审核资料
// @Summary 审核资料
// @Description 管理员审核资料
// @Tags 管理员-资料审核
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "资料ID"
// @Param request body model.ReviewMaterialRequest true "审核信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/materials/{id}/review [post]
func (h *ReviewHandler) ReviewMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的资料ID")
		return
	}

	var req model.ReviewMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 获取当前用户ID（审核人）
	reviewerID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	// 将 Status 转换为 approved 布尔值
	approved := req.Status == model.StatusApproved

	if err := h.reviewService.ReviewMaterial(c.Request.Context(), uint(id), reviewerID, approved, req.RejectionReason); err != nil {
		switch err {
		case service.ErrMaterialAlreadyReviewed:
			response.Error(c, response.ErrForbidden, "该资料已审核")
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// HandleReport 处理举报
// @Summary 处理举报
// @Description 管理员处理举报
// @Tags 管理员-举报
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "举报ID"
// @Param request body HandleReportRequest true "处理信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/reports/{id}/handle [post]
func (h *ReviewHandler) HandleReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的举报ID")
		return
	}

	var req HandleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 获取当前用户ID（处理人）
	handlerID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.reviewService.HandleReport(c.Request.Context(), uint(id), handlerID, req.Approved, req.Note); err != nil {
		switch err {
		case service.ErrAlreadyReviewed:
			response.Error(c, response.ErrForbidden, "该举报已处理")
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// GetReviewHistory 获取审核历史
// @Summary 审核历史
// @Description 获取审核历史记录
// @Tags 管理员-审核
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param target_type query string false "目标类型"
// @Param target_id query string false "目标ID"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/v1/admin/review/history [get]
func (h *ReviewHandler) GetReviewHistory(c *gin.Context) {
	var req ListReviewHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	page := req.Page
	pageSize := req.PageSize

	var targetType *model.ReviewTarget
	var targetID *uint

	if req.TargetType != "" {
		t := model.ReviewTarget(req.TargetType)
		targetType = &t
	}

	if req.TargetID != "" {
		id, err := strconv.ParseUint(req.TargetID, 10, 32)
		if err != nil {
			response.Error(c, response.ErrInvalidParams, "无效的目标ID")
			return
		}
		uid := uint(id)
		targetID = &uid
	}

	records, total, err := h.reviewService.GetReviewHistory(c.Request.Context(), *targetType, *targetID, page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.SuccessWithPaginate(c, total, page, pageSize, records)
}

// GetReviewerStatistics 获取审核人统计信息
// @Summary 审核人统计
// @Description 获取审核人的统计信息
// @Tags 管理员-审核
// @Produce json
// @Security Bearer
// @Param id path int true "审核人ID"
// @Success 200 {object} response.Response{data=service.ReviewerStatistics}
// @Router /api/v1/admin/reviewers/{id}/statistics [get]
func (h *ReviewHandler) GetReviewerStatistics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的审核人ID")
		return
	}

	stats, err := h.reviewService.GetReviewerStatistics(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, stats)
}
