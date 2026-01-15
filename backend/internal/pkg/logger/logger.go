package logger

import (
	"os"

	"github.com/study-upc/backend/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log  *zap.Logger
	sugar *zap.SugaredLogger
)

// Init 初始化日志系统
func Init(cfg *config.Config) error {
	// 日志级别
	level := zapcore.InfoLevel
	switch cfg.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 根据模式选择编码器
	var encoder zapcore.Encoder
	if cfg.Server.Mode == "debug" {
		// 开发模式使用控制台编码器（彩色输出）
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 生产模式使用 JSON 编码器
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 日志轮转
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,    // MB
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,     // days
		Compress:   true,
	}

	// 同时输出到文件和控制台
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	)

	// 创建 logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = log.Sugar()

	return nil
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// GetLogger 获取原生 logger
func GetLogger() *zap.Logger {
	return log
}

// GetSugar 获取 sugar logger
func GetSugar() *zap.SugaredLogger {
	return sugar
}

// Sync 同步日志缓冲区
func Sync() error {
	return log.Sync()
}
