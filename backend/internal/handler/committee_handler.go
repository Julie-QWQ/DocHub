package handler

import (
	"strconv"

	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// CommitteeHandler 学委申请处理器
type CommitteeHandler struct {
	committeeService service.CommitteeService
}

// NewCommitteeHandler 创建学委申请处理器实例
func NewCommitteeHandler(committeeService service.CommitteeService) *CommitteeHandler {
	return &CommitteeHandler{
		committeeService: committeeService,
	}
}

// ApplyForCommitteeRequest 申请学委请求
type ApplyForCommitteeRequest struct {
	Reason string `json:"reason" binding:"required,min=10,max=500"` // 申请理由
}

// ReviewApplicationRequest 审核申请请求
type ReviewApplicationRequest struct {
	Approved bool   `json:"approved" binding:"required"` // 是否通过
	Comment  string `json:"comment" binding:"omitempty,max=500"` // 审核意见
}

// ListApplicationsRequest 申请列表请求参数
type ListApplicationsRequest struct {
	Page   int    `form:"page,default=1"`
	PageSize int   `form:"page_size,default=20"`
	Status string `form:"status"` // 申请状态过滤
}

// ApplyForCommittee 申请学委
// @Summary 申请学委
// @Description 学生申请成为学委
// @Tags 学委申请
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ApplyForCommitteeRequest true "申请信息"
// @Success 200 {object} response.Response{data=model.CommitteeApplication}
// @Router /api/v1/user/apply-committee [post]
func (h *CommitteeHandler) ApplyForCommittee(c *gin.Context) {
	var req ApplyForCommitteeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	application, err := h.committeeService.ApplyForCommittee(c.Request.Context(), userID, req.Reason)
	if err != nil {
		switch err {
		case service.ErrAlreadyCommittee:
			response.Error(c, response.ErrForbidden, "您已经是学委")
		case service.ErrHasPendingApplication:
			response.Error(c, response.ErrForbidden, "已存在待审核的学委申请")
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, application)
}

// ListMyApplications 获取我的申请列表
// @Summary 我的申请列表
// @Description 获取当前用户的学委申请列表
// @Tags 学委申请
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "申请状态" Enums(pending, approved, rejected, cancelled)
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/v1/user/applications [get]
func (h *CommitteeHandler) ListMyApplications(c *gin.Context) {
	var req ListApplicationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	page := req.Page
	pageSize := req.PageSize

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	var status *model.ApplicationStatus
	if req.Status != "" {
		s := model.ApplicationStatus(req.Status)
		status = &s
	}

	applications, total, err := h.committeeService.ListMyApplications(c.Request.Context(), userID, page, pageSize, status)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.SuccessWithPaginate(c, total, page, pageSize, applications)
}

// GetApplication 获取申请详情
// @Summary 获取申请详情
// @Description 获取学委申请详情
// @Tags 学委申请
// @Produce json
// @Security Bearer
// @Param id path int true "申请ID"
// @Success 200 {object} response.Response{data=model.CommitteeApplication}
// @Router /api/v1/user/applications/{id} [get]
func (h *CommitteeHandler) GetApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的申请ID")
		return
	}

	application, err := h.committeeService.GetApplication(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, application)
}

// CancelApplication 取消申请
// @Summary 取消申请
// @Description 取消待审核的学委申请
// @Tags 学委申请
// @Produce json
// @Security Bearer
// @Param id path int true "申请ID"
// @Success 200 {object} response.Response
// @Router /api/v1/user/applications/{id}/cancel [post]
func (h *CommitteeHandler) CancelApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的申请ID")
		return
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.committeeService.CancelApplication(c.Request.Context(), uint(id), userID); err != nil {
		switch err {
		case service.ErrApplicationNotPending:
			response.Error(c, response.ErrForbidden, "该申请不是待审核状态")
		case service.ErrAccessDenied:
			response.Error(c, response.ErrForbidden, "无权操作该申请")
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// ListApplications 获取申请列表（管理员）
// @Summary 学委申请列表
// @Description 管理员获取学委申请列表
// @Tags 管理员-学委申请
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "申请状态"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Router /api/v1/admin/applications [get]
func (h *CommitteeHandler) ListApplications(c *gin.Context) {
	var req ListApplicationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	page := req.Page
	pageSize := req.PageSize

	var status *model.ApplicationStatus
	if req.Status != "" {
		s := model.ApplicationStatus(req.Status)
		status = &s
	}

	applications, total, err := h.committeeService.ListApplications(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.SuccessWithPaginate(c, total, page, pageSize, applications)
}

// ReviewApplication 审核学委申请
// @Summary 审核学委申请
// @Description 管理员审核学委申请
// @Tags 管理员-学委申请
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "申请ID"
// @Param request body ReviewApplicationRequest true "审核信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/applications/{id}/review [post]
func (h *CommitteeHandler) ReviewApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的申请ID")
		return
	}

	var req ReviewApplicationRequest
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

	if err := h.committeeService.ReviewApplication(c.Request.Context(), uint(id), reviewerID, req.Approved, req.Comment); err != nil {
		switch err {
		case service.ErrApplicationNotPending:
			response.Error(c, response.ErrForbidden, "该申请不是待审核状态")
		case service.ErrCannotReviewOwnApplication:
			response.Error(c, response.ErrForbidden, "不能审核自己的申请")
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// GetPendingCount 获取待审核申请数量
// @Summary 待审核申请数量
// @Description 获取待审核学委申请数量
// @Tags 管理员-学委申请
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=int}
// @Router /api/v1/admin/applications/pending/count [get]
func (h *CommitteeHandler) GetPendingCount(c *gin.Context) {
	count, err := h.committeeService.GetPendingCount(c.Request.Context())
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, count)
}
