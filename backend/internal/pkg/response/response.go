package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 常用响应码定义
const (
	CodeSuccess       = 0     // 成功
	CodeInvalidParams = 10001 // 参数错误
	CodeUnauthorized  = 10002 // 未授权
	CodeForbidden     = 10003 // 禁止访问
	CodeNotFound      = 10004 // 资源不存在
	CodeServerError   = 10005 // 服务器错误
	CodeDuplicate     = 10006 // 资源已存在
	CodeDatabaseError = 10007 // 数据库错误

	// 认证相关错误码
	CodeInvalidCredentials = 10101 // 无效的登录凭证
	CodeUserDisabled       = 10102 // 用户已被禁用
	CodeWrongPassword      = 10103 // 旧密码错误
	CodeInvalidToken       = 10104 // Token 无效或已过期
	CodeUserExists         = 10105 // 用户已存在
)

// 错误变量（用于代码中的错误匹配）
var (
	ErrInvalidParams     = CodeInvalidParams
	ErrUnauthorized      = CodeUnauthorized
	ErrForbidden         = CodeForbidden
	ErrNotFound          = CodeNotFound
	ErrInternal          = CodeServerError
	ErrDuplicate         = CodeDuplicate
	ErrInvalidCredentials = CodeInvalidCredentials
	ErrUserDisabled      = CodeUserDisabled
	ErrWrongPassword     = CodeWrongPassword
	ErrInvalidToken      = CodeInvalidToken
	ErrUserExists        = CodeUserExists
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FailWithStatus 带 HTTP 状态码的失败响应
func FailWithStatus(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	Fail(c, CodeInvalidParams, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	FailWithStatus(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	FailWithStatus(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	Fail(c, CodeNotFound, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	Fail(c, CodeServerError, message)
}

// Paginate 分页响应
type PaginateData struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	List  any   `json:"list"`
}

// SuccessWithPaginate 分页成功响应
func SuccessWithPaginate(c *gin.Context, total int64, page int, size int, list any) {
	Success(c, PaginateData{
		Total: total,
		Page:  page,
		Size:  size,
		List:  list,
	})
}
