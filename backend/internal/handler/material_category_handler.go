package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"
)

type MaterialCategoryHandler struct {
	categoryService *service.MaterialCategoryService
}

func NewMaterialCategoryHandler(categoryService *service.MaterialCategoryService) *MaterialCategoryHandler {
	return &MaterialCategoryHandler{
		categoryService: categoryService,
	}
}

// List 获取资料类型列表
// @Summary 获取资料类型列表
// @Tags MaterialCategory
// @Accept json
// @Produce json
// @Param active_only query bool false "是否只获取启用的类型"
// @Success 200 {object} response.Response{data=[]model.MaterialCategoryConfig}
// @Router /api/v1/material-categories [get]
func (h *MaterialCategoryHandler) List(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "false") == "true"

	categories, err := h.categoryService.List(activeOnly)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, categories)
}

// GetByID 根据ID获取资料类型
// @Summary 获取资料类型详情
// @Tags MaterialCategory
// @Accept json
// @Produce json
// @Param id path int true "资料类型ID"
// @Success 200 {object} response.Response{data=model.MaterialCategoryConfig}
// @Router /api/v1/material-categories/:id [get]
func (h *MaterialCategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的ID")
		return
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		response.Error(c, response.ErrNotFound, "资料类型不存在")
		return
	}

	response.Success(c, category)
}

// Create 创建资料类型(管理员)
// @Summary 创建资料类型
// @Tags MaterialCategory
// @Accept json
// @Produce json
// @Param request body model.MaterialCategoryRequest true "资料类型信息"
// @Success 200 {object} response.Response{data=model.MaterialCategoryConfig}
// @Router /api/v1/admin/material-categories [post]
func (h *MaterialCategoryHandler) Create(c *gin.Context) {
	var req model.MaterialCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, "参数错误: "+err.Error())
		return
	}

	category, err := h.categoryService.Create(req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, category)
}

// Update 更新资料类型(管理员)
// @Summary 更新资料类型
// @Tags MaterialCategory
// @Accept json
// @Produce json
// @Param id path int true "资料类型ID"
// @Param request body model.MaterialCategoryRequest true "资料类型信息"
// @Success 200 {object} response.Response{data=model.MaterialCategoryConfig}
// @Router /api/v1/admin/material-categories/:id [put]
func (h *MaterialCategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的ID")
		return
	}

	var req model.MaterialCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, "参数错误: "+err.Error())
		return
	}

	category, err := h.categoryService.Update(uint(id), req)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, category)
}

// Delete 删除资料类型(管理员)
// @Summary 删除资料类型
// @Tags MaterialCategory
// @Accept json
// @Produce json
// @Param id path int true "资料类型ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/material-categories/:id [delete]
func (h *MaterialCategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的ID")
		return
	}

	if err := h.categoryService.Delete(uint(id)); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// ToggleStatus 切换启用状态(管理员)
// @Summary 切换资料类型启用状态
// @Tags MaterialCategory
// @Accept json
// @Produce json
// @Param id path int true "资料类型ID"
// @Success 200 {object} response.Response{data=model.MaterialCategoryConfig}
// @Router /api/v1/admin/material-categories/:id/toggle [post]
func (h *MaterialCategoryHandler) ToggleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, response.ErrInvalidParams, "无效的ID")
		return
	}

	category, err := h.categoryService.ToggleStatus(uint(id))
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, category)
}
