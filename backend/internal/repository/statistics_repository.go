package repository

import (
	"fmt"
	"time"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

// StatisticsRepository 统计数据访问接口
type StatisticsRepository interface {
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
	CreateAccessLog(log *model.AccessLog) error

	// 记录登录日志
	CreateLoginLog(log *model.LoginLog) error
}

type statisticsRepository struct {
	db *gorm.DB
}

// NewStatisticsRepository 创建统计仓库
func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

// GetUserStatistics 获取用户统计
func (r *statisticsRepository) GetUserStatistics() (*model.UserStatistics, error) {
	stats := &model.UserStatistics{
		ByRole: make(map[string]int64),
	}

	// 总用户数
	if err := r.db.Model(&model.User{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 今日新增用户
	today := time.Now().Format("2006-01-02")
	if err := r.db.Model(&model.User{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.Today).Error; err != nil {
		return nil, err
	}

	// 本周新增用户
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	if err := r.db.Model(&model.User{}).
		Where("created_at >= ?", weekStart.Format("2006-01-02")).
		Count(&stats.Week).Error; err != nil {
		return nil, err
	}

	// 本月新增用户
	monthStart := time.Now().Format("2006-01") + "-01"
	if err := r.db.Model(&model.User{}).
		Where("created_at >= ?", monthStart).
		Count(&stats.Month).Error; err != nil {
		return nil, err
	}

	// 今日活跃用户数 (今日登录成功的用户数)
	if err := r.db.Model(&model.LoginLog{}).
		Where("DATE(created_at) = ? AND success = ?", today, true).
		Distinct("user_id").
		Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	// 按角色统计
	var roleStats []struct {
		Role  string
		Count int64
	}
	if err := r.db.Model(&model.User{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Scan(&roleStats).Error; err != nil {
		return nil, err
	}

	for _, rs := range roleStats {
		stats.ByRole[rs.Role] = rs.Count
	}

	return stats, nil
}

// GetUserTrend 获取用户趋势
func (r *statisticsRepository) GetUserTrend(days int) ([]model.TrendData, error) {
	var trend []model.TrendData

	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as count
		FROM users
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	startDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.Raw(query, startDate).Scan(&trend).Error; err != nil {
		return nil, err
	}

	return trend, nil
}

// GetMaterialStatistics 获取资料统计
func (r *statisticsRepository) GetMaterialStatistics() (*model.MaterialStatistics, error) {
	stats := &model.MaterialStatistics{
		ByCategory: make(map[string]int64),
	}

	// 总资料数
	if err := r.db.Model(&model.Material{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 各状态资料数
	var statusStats []struct {
		Status string
		Count  int64
	}
	if err := r.db.Model(&model.Material{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats).Error; err != nil {
		return nil, err
	}

	for _, ss := range statusStats {
		switch ss.Status {
		case "approved":
			stats.Approved = ss.Count
		case "pending":
			stats.Pending = ss.Count
		case "rejected":
			stats.Rejected = ss.Count
		case "offline":
			stats.Offline = ss.Count
		}
	}

	// 今日新增资料
	today := time.Now().Format("2006-01-02")
	if err := r.db.Model(&model.Material{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.Today).Error; err != nil {
		return nil, err
	}

	// 本周新增资料
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	if err := r.db.Model(&model.Material{}).
		Where("created_at >= ?", weekStart.Format("2006-01-02")).
		Count(&stats.Week).Error; err != nil {
		return nil, err
	}

	// 按分类统计
	var categoryStats []struct {
		Category string
		Count    int64
	}
	if err := r.db.Model(&model.Material{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Scan(&categoryStats).Error; err != nil {
		return nil, err
	}

	for _, cs := range categoryStats {
		stats.ByCategory[cs.Category] = cs.Count
	}

	return stats, nil
}

// GetMaterialTrend 获取资料趋势
func (r *statisticsRepository) GetMaterialTrend(days int) ([]model.TrendData, error) {
	var trend []model.TrendData

	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as count
		FROM materials
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	startDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.Raw(query, startDate).Scan(&trend).Error; err != nil {
		return nil, err
	}

	return trend, nil
}

// GetDownloadStatistics 获取下载统计
func (r *statisticsRepository) GetDownloadStatistics() (*model.DownloadStatistics, error) {
	stats := &model.DownloadStatistics{}

	// 总下载次数
	if err := r.db.Model(&model.DownloadRecord{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 今日下载次数
	today := time.Now().Format("2006-01-02")
	if err := r.db.Model(&model.DownloadRecord{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.Today).Error; err != nil {
		return nil, err
	}

	// 本周下载次数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	if err := r.db.Model(&model.DownloadRecord{}).
		Where("created_at >= ?", weekStart.Format("2006-01-02")).
		Count(&stats.Week).Error; err != nil {
		return nil, err
	}

	// 本月下载次数
	monthStart := time.Now().Format("2006-01") + "-01"
	if err := r.db.Model(&model.DownloadRecord{}).
		Where("created_at >= ?", monthStart).
		Count(&stats.Month).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetDownloadTrend 获取下载趋势
func (r *statisticsRepository) GetDownloadTrend(days int) ([]model.TrendData, error) {
	var trend []model.TrendData

	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as count
		FROM download_records
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	startDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.Raw(query, startDate).Scan(&trend).Error; err != nil {
		return nil, err
	}

	return trend, nil
}

// GetApplicationStatistics 获取学委申请统计
func (r *statisticsRepository) GetApplicationStatistics() (*model.ApplicationStatistics, error) {
	stats := &model.ApplicationStatistics{}

	// 总申请数
	if err := r.db.Model(&model.CommitteeApplication{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 各状态申请数
	var statusStats []struct {
		Status string
		Count  int64
	}
	if err := r.db.Model(&model.CommitteeApplication{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats).Error; err != nil {
		return nil, err
	}

	for _, ss := range statusStats {
		switch ss.Status {
		case "pending":
			stats.Pending = ss.Count
		case "approved":
			stats.Approved = ss.Count
		case "rejected":
			stats.Rejected = ss.Count
		}
	}

	return stats, nil
}

// GetVisitStatistics 获取访问统计
func (r *statisticsRepository) GetVisitStatistics() (*model.VisitStatistics, error) {
	stats := &model.VisitStatistics{}

	// 总访问次数
	if err := r.db.Model(&model.AccessLog{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 今日访问次数
	today := time.Now().Format("2006-01-02")
	if err := r.db.Model(&model.AccessLog{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.Today).Error; err != nil {
		return nil, err
	}

	// 本周访问次数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	if err := r.db.Model(&model.AccessLog{}).
		Where("created_at >= ?", weekStart.Format("2006-01-02")).
		Count(&stats.Week).Error; err != nil {
		return nil, err
	}

	// 本月访问次数
	monthStart := time.Now().Format("2006-01") + "-01"
	if err := r.db.Model(&model.AccessLog{}).
		Where("created_at >= ?", monthStart).
		Count(&stats.Month).Error; err != nil {
		return nil, err
	}

	// 独立访客数 (按IP去重)
	if err := r.db.Model(&model.AccessLog{}).
		Select("COUNT(DISTINCT ip_address)").
		Scan(&stats.Unique).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetVisitTrend 获取访问趋势
func (r *statisticsRepository) GetVisitTrend(days int) ([]model.TrendData, error) {
	var trend []model.TrendData

	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as count
		FROM access_logs
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	startDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.Raw(query, startDate).Scan(&trend).Error; err != nil {
		return nil, err
	}

	return trend, nil
}

// CreateAccessLog 创建访问日志
func (r *statisticsRepository) CreateAccessLog(log *model.AccessLog) error {
	return r.db.Create(log).Error
}

// CreateLoginLog 创建登录日志
func (r *statisticsRepository) CreateLoginLog(log *model.LoginLog) error {
	return r.db.Create(log).Error
}

// GetOverviewStatistics 获取概览统计
func (r *statisticsRepository) GetOverviewStatistics() (*model.OverviewStatistics, error) {
	overview := &model.OverviewStatistics{}

	// 并发获取各类统计数据
	var err error
	var wgWait = make(chan struct{})
	var wg int

	go func() {
		userStats, e := r.GetUserStatistics()
		if e != nil {
			err = e
		} else {
			overview.Users = *userStats
		}
		close(wgWait)
		wg++
	}()

	go func() {
		materialStats, e := r.GetMaterialStatistics()
		if e != nil {
			err = e
		} else {
			overview.Materials = *materialStats
		}
		wg++
	}()

	go func() {
		downloadStats, e := r.GetDownloadStatistics()
		if e != nil {
			err = e
		} else {
			overview.Downloads = *downloadStats
		}
		wg++
	}()

	go func() {
		applicationStats, e := r.GetApplicationStatistics()
		if e != nil {
			err = e
		} else {
			overview.Applications = *applicationStats
		}
		wg++
	}()

	go func() {
		visitStats, e := r.GetVisitStatistics()
		if e != nil {
			err = e
		} else {
			overview.Visits = *visitStats
		}
		wg++
	}()

	// 等待所有goroutine完成
	for i := 0; i < 5; i++ {
		<-wgWait
	}

	if err != nil {
		return nil, fmt.Errorf("获取统计数据失败: %w", err)
	}

	return overview, nil
}
