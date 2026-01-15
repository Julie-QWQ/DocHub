package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"code"`              // 业务错误码
	Message string `json:"message"`           // 错误信息
	Err     error  `json:"-"`                 // 原始错误（不输出到日志）
	Details string `json:"details,omitempty"` // 详细信息（可选）
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持 errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus 返回 HTTP 状态码
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeInvalidParams, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeBusinessFail:
		return http.StatusBadRequest
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// 错误码常量
const (
	ErrCodeInternal     = 50000 // 内部错误
	ErrCodeInvalidParams = 40001 // 参数错误
	ErrCodeValidation    = 40002 // 验证错误
	ErrCodeUnauthorized  = 40101 // 未授权
	ErrCodeForbidden     = 40301 // 禁止访问
	ErrCodeNotFound      = 40401 // 资源不存在
	ErrCodeConflict      = 40901 // 冲突
	ErrCodeBusinessFail  = 40003 // 业务失败
	ErrCodeRateLimit     = 42901 // 超出限流
	ErrCodeUserExists    = 40902 // 用户已存在
	ErrCodeInvalidCredentials = 40102 // 凭证无效
)

// New 创建新的应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// 预定义常用错误
var (
	ErrInternal       = New(ErrCodeInternal, "内部服务器错误", nil)
	ErrInvalidParams  = New(ErrCodeInvalidParams, "参数错误", nil)
	ErrValidation     = New(ErrCodeValidation, "数据验证失败", nil)
	ErrUnauthorized   = New(ErrCodeUnauthorized, "未授权访问", nil)
	ErrForbidden      = New(ErrCodeForbidden, "禁止访问", nil)
	ErrNotFound       = New(ErrCodeNotFound, "资源不存在", nil)
	ErrConflict       = New(ErrCodeConflict, "资源冲突", nil)
	ErrBusinessFail   = New(ErrCodeBusinessFail, "业务处理失败", nil)
	ErrRateLimit      = New(ErrCodeRateLimit, "请求过于频繁", nil)
	ErrUserExists     = New(ErrCodeUserExists, "用户已存在", nil)
	ErrInvalidCreds   = New(ErrCodeInvalidCredentials, "用户名或密码错误", nil)
)

// Wrap 包装错误
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}
	return New(ErrCodeInternal, message, err)
}

// Wrapf 包装错误（格式化）
func Wrapf(err error, format string, args ...interface{}) *AppError {
	if err == nil {
		return nil
	}
	return New(ErrCodeInternal, fmt.Sprintf(format, args...), err)
}
