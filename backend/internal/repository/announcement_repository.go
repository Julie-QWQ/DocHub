package repository

import (
	"context"
	"errors"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

var (
	// ErrAnnouncementNotFound 公告不存在
	ErrAnnouncementNotFound = errors.New("公告不存在")
)

// AnnouncementRepository 公告数据访问接口
type AnnouncementRepository interface {
	// Create 创建公告
	Create(ctx context.Context, announcement *model.Announcement) error
	// FindByID 根据ID查找公告
	FindByID(ctx context.Context, id uint) (*model.Announcement, error)
	// FindActive 查找所有启用的公告（按发布时间倒序，限制数量）
	FindActive(ctx context.Context, limit int) ([]model.Announcement, error)
	// List 查询公告列表（支持筛选、分页）
	List(ctx context.Context, req *model.AnnouncementListRequest) ([]model.Announcement, int64, error)
	// Update 更新公告
	Update(ctx context.Context, announcement *model.Announcement) error
	// Delete 删除公告
	Delete(ctx context.Context, id uint) error
}

type announcementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository 创建公告数据访问实例
func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}

// Create 创建公告
func (r *announcementRepository) Create(ctx context.Context, announcement *model.Announcement) error {
	return r.db.WithContext(ctx).Create(announcement).Error
}

// FindByID 根据ID查找公告
func (r *announcementRepository) FindByID(ctx context.Context, id uint) (*model.Announcement, error) {
	var announcement model.Announcement
	err := r.db.WithContext(ctx).
		Preload("Author").
		First(&announcement, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, err
	}
	return &announcement, nil
}

// FindActive 查找所有启用的公告（按发布时间倒序，限制数量）
func (r *announcementRepository) FindActive(ctx context.Context, limit int) ([]model.Announcement, error) {
	var announcements []model.Announcement
	query := r.db.WithContext(ctx).
		Preload("Author").
		Where("is_active = ?", true)

	// 检查是否过期
	query = query.Where("(expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)")

	// 检查是否已发布
	query = query.Where("(published_at IS NULL OR published_at <= CURRENT_TIMESTAMP)")

	// 按优先级降序，高优先级在前；同优先级按发布时间倒序
	err := query.
		Order("CASE WHEN priority = 'high' THEN 1 WHEN priority = 'normal' THEN 2 ELSE 3 END, published_at DESC").
		Limit(limit).
		Find(&announcements).Error

	return announcements, err
}

// List 查询公告列表（支持筛选、分页）
func (r *announcementRepository) List(ctx context.Context, req *model.AnnouncementListRequest) ([]model.Announcement, int64, error) {
	var announcements []model.Announcement
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Announcement{})

	// 筛选条件
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}
	if req.AuthorID > 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 查询数据（按优先级降序，高优先级在前；同优先级按创建时间倒序）
	err := query.
		Preload("Author").
		Order("CASE WHEN priority = 'high' THEN 1 WHEN priority = 'normal' THEN 2 ELSE 3 END, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&announcements).Error

	return announcements, total, err
}

// Update 更新公告
func (r *announcementRepository) Update(ctx context.Context, announcement *model.Announcement) error {
	return r.db.WithContext(ctx).Save(announcement).Error
}

// Delete 删除公告
func (r *announcementRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Announcement{}, id).Error
}
