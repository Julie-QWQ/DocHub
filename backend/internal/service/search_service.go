package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
	"gorm.io/gorm"
)

// SearchService 搜索服务接口
type SearchService interface {
	// Search 搜索资料
	Search(ctx context.Context, userID uint, req *model.SearchRequest) (*model.SearchResponse, error)
	// RecordSearchHistory 记录搜索历史
	RecordSearchHistory(ctx context.Context, userID uint, keyword string, resultCount int) error
	// GetUserSearchHistories 获取用户搜索历史
	GetUserSearchHistories(ctx context.Context, userID uint, limit int) ([]*model.SearchHistory, error)
	// ClearSearchHistory 清空搜索历史
	ClearSearchHistory(ctx context.Context, userID uint) error
	// GetHotKeywords 获取热门搜索词
	GetHotKeywords(ctx context.Context, limit int) ([]*model.HotKeyword, error)
}

// RecommendationService 推荐服务接口
type RecommendationService interface {
	// GetRecommendations 获取推荐资料
	GetRecommendations(ctx context.Context, userID uint, req *model.RecommendationRequest) ([]*model.RecommendationResult, error)
	// GetHotMaterials 获取热门资料
	GetHotMaterials(ctx context.Context, limit int) ([]*model.Material, error)
	// GetPersonalizedRecommendations 获取个性化推荐
	GetPersonalizedRecommendations(ctx context.Context, userID uint, limit int) ([]*model.RecommendationResult, error)
	// GetRelatedMaterials 获取相关资料
	GetRelatedMaterials(ctx context.Context, materialID uint, limit int) ([]*model.RecommendationResult, error)
}

// searchService 搜索服务实现
type searchService struct {
	db                *gorm.DB
	materialRepo      repository.MaterialRepository
	searchHistoryRepo repository.SearchHistoryRepository
	hotKeywordRepo    repository.HotKeywordRepository
	downloadRepo      repository.DownloadRecordRepository
}

// NewSearchService 创建搜索服务实例
func NewSearchService(
	db *gorm.DB,
	materialRepo repository.MaterialRepository,
	searchHistoryRepo repository.SearchHistoryRepository,
	hotKeywordRepo repository.HotKeywordRepository,
	downloadRepo repository.DownloadRecordRepository,
) SearchService {
	return &searchService{
		db:                db,
		materialRepo:      materialRepo,
		searchHistoryRepo: searchHistoryRepo,
		hotKeywordRepo:    hotKeywordRepo,
		downloadRepo:      downloadRepo,
	}
}

// Search 搜索资料
func (s *searchService) Search(ctx context.Context, userID uint, req *model.SearchRequest) (*model.SearchResponse, error) {
	// 基础查询
	query := s.db.WithContext(ctx).
		Where("status = ?", model.StatusApproved) // 只搜索已审核通过的资料

	// 关键词全文搜索
	if req.Keyword != "" {
		// 使用 PostgreSQL 的全文搜索
		query = query.Where(`
			search_vector @@ to_tsquery('simple', ?)
			OR title ILIKE ?
			OR description ILIKE ?
			OR course_name ILIKE ?
		`,
			fmt.Sprintf("%s:*", req.Keyword), // to_tsquery 格式
			fmt.Sprintf("%%%s%%", req.Keyword),
			fmt.Sprintf("%%%s%%", req.Keyword),
			fmt.Sprintf("%%%s%%", req.Keyword),
		)

		// 添加相关度排序
		query = query.Order(fmt.Sprintf(
			"ts_rank(search_vector, to_tsquery('simple', '%s:*')) DESC",
			req.Keyword,
		))
	}

	// 分类筛选
	if req.Category != nil {
		query = query.Where("category = ?", *req.Category)
	}

	// 课程名称筛选
	if req.CourseName != "" {
		query = query.Where("course_name ILIKE ?", fmt.Sprintf("%%%s%%", req.CourseName))
	}

	// 标签筛选 (使用数组包含查询)
	if len(req.Tags) > 0 {
		query = query.Where("tags && ?", "{"+strings.Join(req.Tags, ",")+"}")
	}

	// 时间范围筛选
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("created_at >= ?", startDate)
		}
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			// 添加一天以包含结束日期当天
			endDate = endDate.AddDate(0, 0, 1)
			query = query.Where("created_at < ?", endDate)
		}
	}

	// 排序
	switch req.SortBy {
	case "created_at":
		order := "created_at"
		if req.SortOrder == "asc" {
			order += " ASC"
		} else {
			order += " DESC"
		}
		if req.Keyword != "" {
			// 如果有搜索关键词，相关度优先
			query = query.Order(fmt.Sprintf(
				"ts_rank(search_vector, to_tsquery('simple', '%s:*')) DESC, %s",
				req.Keyword, order,
			))
		} else {
			query = query.Order(order)
		}
	case "download_count":
		order := "download_count"
		if req.SortOrder == "asc" {
			order += " ASC"
		} else {
			order += " DESC"
		}
		query = query.Order(order)
	case "favorite_count":
		order := "favorite_count"
		if req.SortOrder == "asc" {
			order += " ASC"
		} else {
			order += " DESC"
		}
		query = query.Order(order)
	case "view_count":
		order := "view_count"
		if req.SortOrder == "asc" {
			order += " ASC"
		} else {
			order += " DESC"
		}
		query = query.Order(order)
	case "relevance":
		// 相关度排序（仅在有关键词时有效）
		if req.Keyword != "" {
			query = query.Order(fmt.Sprintf(
				"ts_rank(search_vector, to_tsquery('simple', '%s:*')) DESC",
				req.Keyword,
			))
		} else {
			// 没有关键词时默认按创建时间排序
			query = query.Order("created_at DESC")
		}
	default:
		// 默认排序
		if req.Keyword != "" {
			query = query.Order(fmt.Sprintf(
				"ts_rank(search_vector, to_tsquery('simple', '%s:*')) DESC, created_at DESC",
				req.Keyword,
			))
		} else {
			query = query.Order("created_at DESC")
		}
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 执行查询
	var materials []*model.Material
	if err := query.Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("搜索失败: %w", err)
	}

	// 获取总数
	var total int64
	countQuery := s.db.WithContext(ctx).Model(&model.Material{}).Where("status = ?", model.StatusApproved)

	// 重复应用筛选条件（用于计数）
	if req.Keyword != "" {
		countQuery = countQuery.Where(`
			search_vector @@ to_tsquery('simple', ?)
			OR title ILIKE ?
			OR description ILIKE ?
			OR course_name ILIKE ?
		`,
			fmt.Sprintf("%s:*", req.Keyword),
			fmt.Sprintf("%%%s%%", req.Keyword),
			fmt.Sprintf("%%%s%%", req.Keyword),
			fmt.Sprintf("%%%s%%", req.Keyword),
		)
	}
	if req.Category != nil {
		countQuery = countQuery.Where("category = ?", *req.Category)
	}
	if req.CourseName != "" {
		countQuery = countQuery.Where("course_name ILIKE ?", fmt.Sprintf("%%%s%%", req.CourseName))
	}
	if len(req.Tags) > 0 {
		countQuery = countQuery.Where("tags && ?", "{"+strings.Join(req.Tags, ",")+"}")
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			countQuery = countQuery.Where("created_at >= ?", startDate)
		}
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			endDate = endDate.AddDate(0, 0, 1)
			countQuery = countQuery.Where("created_at < ?", endDate)
		}
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %w", err)
	}

	// 构建搜索结果
	results := make([]*model.SearchResult, 0, len(materials))
	for _, material := range materials {
		results = append(results, &model.SearchResult{
			Material:  material,
			Relevance: 0.8, // 简化的相关度计算
		})
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	// 异步记录搜索历史和更新热门搜索词
	go func() {
		// 记录搜索历史
		if userID != 0 {
			_ = s.RecordSearchHistory(context.Background(), userID, req.Keyword, int(total))
		}
		// 更新热门搜索词
		_ = s.hotKeywordRepo.IncrementKeywordCount(context.Background(), req.Keyword)
	}()

	return &model.SearchResponse{
		Results:    results,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// RecordSearchHistory 记录搜索历史
func (s *searchService) RecordSearchHistory(ctx context.Context, userID uint, keyword string, resultCount int) error {
	history := &model.SearchHistory{
		UserID:      userID,
		Keyword:     keyword,
		ResultCount: resultCount,
	}
	return s.searchHistoryRepo.CreateSearchHistory(ctx, history)
}

// GetUserSearchHistories 获取用户搜索历史
func (s *searchService) GetUserSearchHistories(ctx context.Context, userID uint, limit int) ([]*model.SearchHistory, error) {
	return s.searchHistoryRepo.GetUserSearchHistories(ctx, userID, limit)
}

// ClearSearchHistory 清空搜索历史
func (s *searchService) ClearSearchHistory(ctx context.Context, userID uint) error {
	return s.searchHistoryRepo.ClearUserSearchHistories(ctx, userID)
}

// GetHotKeywords 获取热门搜索词
func (s *searchService) GetHotKeywords(ctx context.Context, limit int) ([]*model.HotKeyword, error) {
	return s.hotKeywordRepo.GetHotKeywords(ctx, limit)
}
