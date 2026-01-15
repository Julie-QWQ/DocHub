package service

import (
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
)

// StatisticsService 统计服务接口
type StatisticsService interface {
	// 概览统计
	GetOverviewStatistics() (*model.OverviewStatistics, error)

	// 用户统计
	GetUserStatistics() (*model.UserStatistics, error)
	GetUserTrend(days int) ([]model.TrendData, error)

	// 资料统计
	GetMaterialStatistics() (*model.MaterialStatistics, error)
	GetMaterialTrend(days int) ([]model.TrendData, error)

	// 下载统计
	GetDownloadStatistics() (*model.DownloadStatistics, error)
	GetDownloadTrend(days int) ([]model.TrendData, error)

	// 学委申请统计
	GetApplicationStatistics() (*model.ApplicationStatistics, error)

	// 访问统计
	GetVisitStatistics() (*model.VisitStatistics, error)
	GetVisitTrend(days int) ([]model.TrendData, error)

	// 记录访问日志
	RecordAccess(userID *uint, ip, path, method, userAgent, referer string) error

	// 记录登录日志
	RecordLoginLog(userID uint, ip, userAgent string, success bool) error
}

type statisticsService struct {
	statsRepo repository.StatisticsRepository
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(statsRepo repository.StatisticsRepository) StatisticsService {
	return &statisticsService{
		statsRepo: statsRepo,
	}
}

// GetOverviewStatistics 获取概览统计
func (s *statisticsService) GetOverviewStatistics() (*model.OverviewStatistics, error) {
	return s.statsRepo.GetOverviewStatistics()
}

// GetUserStatistics 获取用户统计
func (s *statisticsService) GetUserStatistics() (*model.UserStatistics, error) {
	return s.statsRepo.GetUserStatistics()
}

// GetUserTrend 获取用户趋势
func (s *statisticsService) GetUserTrend(days int) ([]model.TrendData, error) {
	if days <= 0 || days > 365 {
		days = 30 // 默认30天
	}
	return s.statsRepo.GetUserTrend(days)
}

// GetMaterialStatistics 获取资料统计
func (s *statisticsService) GetMaterialStatistics() (*model.MaterialStatistics, error) {
	return s.statsRepo.GetMaterialStatistics()
}

// GetMaterialTrend 获取资料趋势
func (s *statisticsService) GetMaterialTrend(days int) ([]model.TrendData, error) {
	if days <= 0 || days > 365 {
		days = 30 // 默认30天
	}
	return s.statsRepo.GetMaterialTrend(days)
}

// GetDownloadStatistics 获取下载统计
func (s *statisticsService) GetDownloadStatistics() (*model.DownloadStatistics, error) {
	return s.statsRepo.GetDownloadStatistics()
}

// GetDownloadTrend 获取下载趋势
func (s *statisticsService) GetDownloadTrend(days int) ([]model.TrendData, error) {
	if days <= 0 || days > 365 {
		days = 30 // 默认30天
	}
	return s.statsRepo.GetDownloadTrend(days)
}

// GetApplicationStatistics 获取学委申请统计
func (s *statisticsService) GetApplicationStatistics() (*model.ApplicationStatistics, error) {
	return s.statsRepo.GetApplicationStatistics()
}

// GetVisitStatistics 获取访问统计
func (s *statisticsService) GetVisitStatistics() (*model.VisitStatistics, error) {
	return s.statsRepo.GetVisitStatistics()
}

// GetVisitTrend 获取访问趋势
func (s *statisticsService) GetVisitTrend(days int) ([]model.TrendData, error) {
	if days <= 0 || days > 365 {
		days = 30 // 默认30天
	}
	return s.statsRepo.GetVisitTrend(days)
}

// RecordAccess 记录访问日志
func (s *statisticsService) RecordAccess(userID *uint, ip, path, method, userAgent, referer string) error {
	// 过滤掉静态资源和健康检查接口
	if path == "/health" || path == "/liveness" {
		return nil
	}

	log := &model.AccessLog{
		UserID:    userID,
		IPAddress: ip,
		Path:      path,
		Method:    method,
		UserAgent: userAgent,
		Referer:   referer,
	}

	return s.statsRepo.CreateAccessLog(log)
}

// RecordLoginLog 记录登录日志
func (s *statisticsService) RecordLoginLog(userID uint, ip, userAgent string, success bool) error {
	log := &model.LoginLog{
		UserID:    userID,
		IPAddress: ip,
		UserAgent: userAgent,
		Success:   success,
	}

	return s.statsRepo.CreateLoginLog(log)
}
