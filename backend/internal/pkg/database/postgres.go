package database

import (
	"fmt"
	"time"

	"github.com/study-upc/backend/internal/pkg/config"
	"github.com/study-upc/backend/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitPostgres 初始化 PostgreSQL 连接
func InitPostgres(cfg *config.Config) error {
	dsn := cfg.Database.GetDSN()

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
		// 禁用外键约束（在应用层处理）
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 如果是调试模式，打印 SQL
	if cfg.Server.Mode == "debug" {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库 Ping 失败: %w", err)
	}

	DB = db
	logger.Info("数据库连接成功", zap.String("host", cfg.Database.Host))
	return nil
}

// ClosePostgres 关闭数据库连接
func ClosePostgres() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
