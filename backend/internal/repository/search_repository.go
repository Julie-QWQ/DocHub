package repository

import (
	"context"
	"time"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

// SearchHistoryRepository 搜索历史仓储接口
type SearchHistoryRepository interface {
	// CreateSearchHistory 创建搜索历史记录
	CreateSearchHistory(ctx context.Context, history *model.SearchHistory) error
	// GetUserSearchHistories 获取用户搜索历史
	GetUserSearchHistories(ctx context.Context, userID uint, limit int) ([]*model.SearchHistory, error)
	// DeleteSearchHistory 删除搜索历史
	DeleteSearchHistory(ctx context.Context, id uint) error
	// ClearUserSearchHistories 清空用户搜索历史
	ClearUserSearchHistories(ctx context.Context, userID uint) error
}

// HotKeywordRepository 热门搜索词仓储接口
type HotKeywordRepository interface {
	// IncrementKeywordCount 增加关键词搜索次数
	IncrementKeywordCount(ctx context.Context, keyword string) error
	// GetHotKeywords 获取热门搜索词
	GetHotKeywords(ctx context.Context, limit int) ([]*model.HotKeyword, error)
	// UpdateLastSearchedAt 更新最后搜索时间
	UpdateLastSearchedAt(ctx context.Context, keyword string) error
}

// searchHistoryRepository 搜索历史仓储实现
type searchHistoryRepository struct {
	db *gorm.DB
}

// NewSearchHistoryRepository 创建搜索历史仓储实例
func NewSearchHistoryRepository(db *gorm.DB) SearchHistoryRepository {
	return &searchHistoryRepository{db: db}
}

// CreateSearchHistory 创建搜索历史记录
func (r *searchHistoryRepository) CreateSearchHistory(ctx context.Context, history *model.SearchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetUserSearchHistories 获取用户搜索历史
func (r *searchHistoryRepository) GetUserSearchHistories(ctx context.Context, userID uint, limit int) ([]*model.SearchHistory, error) {
	var histories []*model.SearchHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

// DeleteSearchHistory 删除搜索历史
func (r *searchHistoryRepository) DeleteSearchHistory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.SearchHistory{}, id).Error
}

// ClearUserSearchHistories 清空用户搜索历史
func (r *searchHistoryRepository) ClearUserSearchHistories(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.SearchHistory{}).Error
}

// hotKeywordRepository 热门搜索词仓储实现
type hotKeywordRepository struct {
	db *gorm.DB
}

// NewHotKeywordRepository 创建热门搜索词仓储实例
func NewHotKeywordRepository(db *gorm.DB) HotKeywordRepository {
	return &hotKeywordRepository{db: db}
}

// IncrementKeywordCount 增加关键词搜索次数
func (r *hotKeywordRepository) IncrementKeywordCount(ctx context.Context, keyword string) error {
	// 使用 ON CONFLICT 来实现 upsert
	return r.db.WithContext(ctx).
		Exec(`
			INSERT INTO hot_keywords (keyword, search_count, last_searched_at, created_at, updated_at)
			VALUES (?, 1, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT (keyword)
			DO UPDATE SET
				search_count = hot_keywords.search_count + 1,
				last_searched_at = ?,
				updated_at = CURRENT_TIMESTAMP
		`, keyword, time.Now(), time.Now()).Error
}

// GetHotKeywords 获取热门搜索词
func (r *hotKeywordRepository) GetHotKeywords(ctx context.Context, limit int) ([]*model.HotKeyword, error) {
	var keywords []*model.HotKeyword
	err := r.db.WithContext(ctx).
		Order("search_count DESC, last_searched_at DESC").
		Limit(limit).
		Find(&keywords).Error
	return keywords, err
}

// UpdateLastSearchedAt 更新最后搜索时间
func (r *hotKeywordRepository) UpdateLastSearchedAt(ctx context.Context, keyword string) error {
	return r.db.WithContext(ctx).
		Model(&model.HotKeyword{}).
		Where("keyword = ?", keyword).
		Update("last_searched_at", time.Now()).Error
}
