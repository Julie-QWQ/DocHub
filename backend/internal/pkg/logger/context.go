package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// WithContext 创建带上下文的日志字段
func WithContext(ctx context.Context) []zap.Field {
	var fields []zap.Field

	// 从上下文中提取 request_id
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, zap.Any("request_id", requestID))
	}

	// 从上下文中提取 user_id
	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, zap.Any("user_id", userID))
	}

	return fields
}

// LogOperation 记录业务操作日志
func LogOperation(ctx context.Context, operation string, success bool, duration time.Duration, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("operation", operation),
		zap.Bool("success", success),
		zap.Duration("duration_ms", duration),
	}, fields...)

	allFields = append(allFields, WithContext(ctx)...)

	if success {
		Info("operation_success", allFields...)
	} else {
		Error("operation_failed", allFields...)
	}
}

// LogAPIRequest 记录API请求日志
func LogAPIRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	fields := append([]zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration_ms", duration),
	}, WithContext(ctx)...)

	Info("api_request", fields...)
}

// LogDBQuery 记录数据库查询日志
func LogDBQuery(ctx context.Context, table, operation string, duration time.Duration, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("table", table),
		zap.String("operation", operation),
		zap.Duration("duration_ms", duration),
	}, fields...)

	allFields = append(allFields, WithContext(ctx)...)

	Debug("db_query", allFields...)
}

// LogError 记录错误日志（带上下文）
func LogError(ctx context.Context, msg string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.Error(err),
	}, fields...)

	allFields = append(allFields, WithContext(ctx)...)

	Error(msg, allFields...)
}
