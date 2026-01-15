package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/service"
)

// SystemHandler 系统处理器
type SystemHandler struct {
	adminService service.AdminService
}

// NewSystemHandler 创建系统处理器
func NewSystemHandler(adminService service.AdminService) *SystemHandler {
	return &SystemHandler{
		adminService: adminService,
	}
}

// GetPublicSystemConfig 获取单个公开的系统配置
// @Summary 获取系统配置
// @Description 获取公开的系统配置（如网站名称、描述等）
// @Tags 系统
// @Accept json
// @Produce json
// @Param key path string true "配置键"
// @Success 200 {object} response.Response{data=string}
// @Router /api/v1/system/configs/{key} [get]
func (h *SystemHandler) GetPublicSystemConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.Error(c, response.CodeInvalidParams, "配置键不能为空")
		return
	}

	// 允许的公开配置键
	publicKeys := map[string]bool{
		"site_name":        true,
		"site_description": true,
		"maintenance_mode": true,
	}

	if !publicKeys[key] {
		response.Error(c, response.CodeNotFound, "配置不存在")
		return
	}

	config, err := h.adminService.GetSystemConfig(key)
	if err != nil {
		response.Error(c, response.CodeNotFound, "配置不存在")
		return
	}

	response.Success(c, config.ConfigValue)
}

// GetPublicSystemConfigs 批量获取公开的系统配置
// @Summary 批量获取系统配置
// @Description 批量获取公开的系统配置
// @Tags 系统
// @Accept json
// @Produce json
// @Param keys query string true "配置键，逗号分隔"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/system/configs [get]
func (h *SystemHandler) GetPublicSystemConfigs(c *gin.Context) {
	keysStr := c.Query("keys")
	if keysStr == "" {
		response.Error(c, response.CodeInvalidParams, "配置键不能为空")
		return
	}

	// 分割配置键
	keys := strings.Split(keysStr, ",")

	// 允许的公开配置键
	publicKeys := map[string]bool{
		"site_name":        true,
		"site_description": true,
		"maintenance_mode": true,
	}

	result := make(map[string]string)
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if !publicKeys[key] {
			continue
		}

		config, err := h.adminService.GetSystemConfig(key)
		if err == nil {
			result[key] = config.ConfigValue
		}
	}

	response.Success(c, result)
}

// GetUploadConfig 获取上传配置
// @Summary 获取上传配置
// @Description 获取文件上传的大小限制和允许的文件类型
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/v1/system/upload-config [get]
func (h *SystemHandler) GetUploadConfig(c *gin.Context) {
	// 获取最大上传大小
	maxSizeConfig, err := h.adminService.GetSystemConfig("max_upload_size")
	var maxSize int64 = 52428800 // 默认 50MB
	if err == nil {
		parsedSize := parseInt64(maxSizeConfig.ConfigValue, maxSize)
		maxSize = parsedSize
	}

	// 获取允许的文件类型
	allowedTypesConfig, err := h.adminService.GetSystemConfig("allowed_file_types")
	var allowedTypes []string
	if err == nil {
		// 分割逗号分隔的文件类型
		typesStr := allowedTypesConfig.ConfigValue
		if typesStr != "" {
			allowedTypes = strings.Split(typesStr, ",")
			// 去除空格
			for i, t := range allowedTypes {
				allowedTypes[i] = strings.TrimSpace(t)
			}
		}
	}

	// 如果没有配置,使用默认值
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"pdf", "docx", "doc", "pptx", "ppt", "txt", "md", "zip", "rar"}
	}

	response.Success(c, map[string]interface{}{
		"max_size":        maxSize,
		"allowed_types":   allowedTypes,
		"max_size_mb":     maxSize / 1024 / 1024, // 转换为 MB 供前端使用
		"accept":          "." + strings.Join(allowedTypes, ",."), // 文件选择器的 accept 属性
	})
}

// parseInt64 解析字符串为 int64,失败时返回默认值
func parseInt64(s string, defaultValue int64) int64 {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}
