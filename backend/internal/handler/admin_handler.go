package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	adminService service.AdminService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// ============ 系统配置管理 ============

// ListSystemConfigs 获取系统配置列表
// @Summary 获取系统配置列表
// @Description 获取系统配置列表
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "配置分类"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=[]model.SystemConfig}
// @Router /api/v1/admin/configs [get]
func (h *AdminHandler) ListSystemConfigs(c *gin.Context) {
	var req model.SystemConfigListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	configs, total, err := h.adminService.ListSystemConfigs(&req)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取系统配置失败")
		return
	}

	response.SuccessWithPaginate(c, total, req.Page, req.PageSize, configs)
}

// GetSystemConfig 获取单个系统配置
// @Summary 获取单个系统配置
// @Description 根据配置键获取系统配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "配置键"
// @Success 200 {object} response.Response{data=model.SystemConfig}
// @Router /api/v1/admin/configs/{key} [get]
func (h *AdminHandler) GetSystemConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.Error(c, response.CodeInvalidParams, "配置键不能为空")
		return
	}

	config, err := h.adminService.GetSystemConfig(key)
	if err != nil {
		response.Error(c, response.CodeNotFound, "配置不存在")
		return
	}

	response.Success(c, config)
}

// UpdateSystemConfig 更新系统配置
// @Summary 更新系统配置
// @Description 更新系统配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UpdateSystemConfigRequest true "配置信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/configs [put]
func (h *AdminHandler) UpdateSystemConfig(c *gin.Context) {
	var req model.UpdateSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	if err := h.adminService.UpdateSystemConfig(req.ConfigKey, req.ConfigValue); err != nil {
		response.Error(c, response.CodeServerError, "更新配置失败")
		return
	}

	response.Success(c, nil)
}

// CreateSystemConfig 创建系统配置
// @Summary 创建系统配置
// @Description 创建系统配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.SystemConfig true "配置信息"
// @Success 200 {object} response.Response{data=model.SystemConfig}
// @Router /api/v1/admin/configs [post]
func (h *AdminHandler) CreateSystemConfig(c *gin.Context) {
	var config model.SystemConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	if err := h.adminService.CreateSystemConfig(&config); err != nil {
		response.Error(c, response.CodeServerError, "创建配置失败")
		return
	}

	response.Success(c, config)
}

// DeleteSystemConfig 删除系统配置
// @Summary 删除系统配置
// @Description 删除系统配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "配置键"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/configs/{key} [delete]
func (h *AdminHandler) DeleteSystemConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.Error(c, response.CodeInvalidParams, "配置键不能为空")
		return
	}

	if err := h.adminService.DeleteSystemConfig(key); err != nil {
		response.Error(c, response.CodeServerError, "删除配置失败")
		return
	}

	response.Success(c, nil)
}

// ============ 用户管理 ============

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param role query string false "角色筛选"
// @Param status query string false "状态筛选"
// @Param major query string false "专业筛选"
// @Param class query string false "班级筛选"
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方向"
// @Success 200 {object} response.Response{data=[]model.User}
// @Router /api/v1/admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var req model.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	users, total, err := h.adminService.ListUsers(&req)
	if err != nil {
		response.Error(c, response.CodeServerError, "获取用户列表失败")
		return
	}

	response.SuccessWithPaginate(c, total, req.Page, req.PageSize, users)
}

// GetUserDetail 获取用户详情
// @Summary 获取用户详情
// @Description 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.UserDetailResponse}
// @Router /api/v1/admin/users/{id} [get]
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams, "用户ID格式错误")
		return
	}

	detail, err := h.adminService.GetUserDetail(uint(id))
	if err != nil {
		response.Error(c, response.CodeNotFound, "获取用户详情失败")
		return
	}

	response.Success(c, detail)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户状态 (激活/禁用/封禁)
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body model.UpdateUserStatusRequest true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams, "用户ID格式错误")
		return
	}

	var req model.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	// 获取当前用户ID
	currentUserID, _ := middleware.GetUserID(c)

	// 不能修改自己的状态
	if uint(id) == currentUserID {
		response.Error(c, response.CodeInvalidParams, "不能修改自己的状态")
		return
	}

	if err := h.adminService.UpdateUserStatus(uint(id), req.Status, req.Reason); err != nil {
		response.Error(c, response.CodeServerError, "更新用户状态失败")
		return
	}

	response.Success(c, nil)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 管理员更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body object true "用户信息"
// @Success 200 {object} response.Response{data=model.User}
// @Router /api/v1/admin/users/{id} [put]
func (h *AdminHandler) UpdateUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams, "用户ID格式错误")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, response.CodeInvalidParams, "参数错误")
		return
	}

	if err := h.adminService.UpdateUserInfo(uint(id), updates); err != nil {
		response.Error(c, response.CodeServerError, "更新用户信息失败")
		return
	}

	// 获取更新后的用户信息
	user, err := h.adminService.GetUserDetail(uint(id))
	if err != nil {
		response.Error(c, response.CodeServerError, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户 (软删除)
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.CodeInvalidParams, "用户ID格式错误")
		return
	}

	// 获取当前用户ID
	currentUserID, _ := middleware.GetUserID(c)

	// 不能删除自己
	if uint(id) == currentUserID {
		response.Error(c, response.CodeInvalidParams, "不能删除自己")
		return
	}

	if err := h.adminService.DeleteUser(uint(id)); err != nil {
		response.Error(c, response.CodeServerError, "删除用户失败")
		return
	}

	response.Success(c, nil)
}
