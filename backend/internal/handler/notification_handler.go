package handler

import (
	"strconv"

	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notificationService service.NotificationService
}

// NewNotificationHandler 创建通知处理器实例
func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// ListNotificationsRequest 通知列表请求参数
type ListNotificationsRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// ListNotifications 获取通知列表
// @Summary 通知列表
// @Description 获取当前用户的通知列表
// @Tags 通知
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PaginateData}
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	page, pageSize := GetPageParams(c)

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	notifications, total, err := h.notificationService.ListNotifications(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.SuccessWithPaginate(c, total, page, pageSize, notifications)
}

// GetUnreadNotifications 获取未读通知
// @Summary 未读通知
// @Description 获取当前用户的未读通知
// @Tags 通知
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=[]model.Notification}
// @Router /api/v1/notifications/unread [get]
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	notifications, err := h.notificationService.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, notifications)
}

// GetUnreadCount 获取未读通知数量
// @Summary 未读通知数量
// @Description 获取当前用户的未读通知数量
// @Tags 通知
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=int}
// @Router /api/v1/notifications/unread/count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, count)
}

// MarkAsRead 标记通知为已读
// @Summary 标记已读
// @Description 标记通知为已读
// @Tags 通知
// @Produce json
// @Security Bearer
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的通知ID")
		return
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.notificationService.MarkAsRead(c.Request.Context(), uint(id), userID); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// MarkAllAsRead 标记所有通知为已读
// @Summary 全部标记已读
// @Description 标记当前用户的所有通知为已读
// @Tags 通知
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteNotification 删除通知
// @Summary 删除通知
// @Description 删除指定通知
// @Tags 通知
// @Produce json
// @Security Bearer
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的通知ID")
		return
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.notificationService.DeleteNotification(c.Request.Context(), uint(id), userID); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}
