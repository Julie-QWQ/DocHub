package repository

import (
	"context"
	"errors"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

var (
	// ErrReviewRecordNotFound 审核记录不存在
	ErrReviewRecordNotFound = errors.New("审核记录不存在")
)

// ReviewRepository 审核记录数据访问层接口
type ReviewRepository interface {
	// CreateReviewRecord 创建审核记录
	CreateReviewRecord(ctx context.Context, record *model.ReviewRecord) error
	// FindReviewRecordByID 根据ID查找审核记录
	FindReviewRecordByID(ctx context.Context, id uint) (*model.ReviewRecord, error)
	// ListReviewRecords 分页获取审核记录列表
	ListReviewRecords(ctx context.Context, page, pageSize int, targetType *model.ReviewTarget, targetID *uint) ([]*model.ReviewRecord, int64, error)
	// ListReviewRecordsByReviewerID 分页获取审核人的审核记录
	ListReviewRecordsByReviewerID(ctx context.Context, reviewerID uint, page, pageSize int) ([]*model.ReviewRecord, int64, error)
	// CountReviewRecords 统计审核记录数量
	CountReviewRecords(ctx context.Context) (int64, error)
}

// reviewRepository 审核记录数据访问层实现
type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository 创建审核记录数据访问层实例
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// CreateReviewRecord 创建审核记录
func (r *reviewRepository) CreateReviewRecord(ctx context.Context, record *model.ReviewRecord) error {
	result := r.db.WithContext(ctx).Create(record)
	return result.Error
}

// FindReviewRecordByID 根据ID查找审核记录
func (r *reviewRepository) FindReviewRecordByID(ctx context.Context, id uint) (*model.ReviewRecord, error) {
	var record model.ReviewRecord
	result := r.db.WithContext(ctx).
		Preload("Reviewer").
		First(&record, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrReviewRecordNotFound
		}
		return nil, result.Error
	}
	return &record, nil
}

// ListReviewRecords 分页获取审核记录列表
func (r *reviewRepository) ListReviewRecords(ctx context.Context, page, pageSize int, targetType *model.ReviewTarget, targetID *uint) ([]*model.ReviewRecord, int64, error) {
	var records []*model.ReviewRecord
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.ReviewRecord{})

	// 如果指定了目标类型，添加过滤
	if targetType != nil {
		query = query.Where("target_type = ?", *targetType)
	}

	// 如果指定了目标ID，添加过滤
	if targetID != nil {
		query = query.Where("target_id = ?", *targetID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	result := query.
		Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return records, total, nil
}

// ListReviewRecordsByReviewerID 分页获取审核人的审核记录
func (r *reviewRepository) ListReviewRecordsByReviewerID(ctx context.Context, reviewerID uint, page, pageSize int) ([]*model.ReviewRecord, int64, error) {
	var records []*model.ReviewRecord
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.ReviewRecord{}).Where("reviewer_id = ?", reviewerID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	result := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return records, total, nil
}

// CountReviewRecords 统计审核记录数量
func (r *reviewRepository) CountReviewRecords(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.ReviewRecord{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
