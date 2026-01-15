package repository

import (
	"context"
	"errors"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

var (
	// ErrApplicationNotFound 申请记录不存在
	ErrApplicationNotFound = errors.New("申请记录不存在")
	// ErrApplicationAlreadyExists 申请已存在
	ErrApplicationAlreadyExists = errors.New("申请已存在")
	// ErrInvalidApplicationStatus 无效的申请状态
	ErrInvalidApplicationStatus = errors.New("无效的申请状态")
)

// CommitteeRepository 学委申请数据访问层接口
type CommitteeRepository interface {
	// CreateApplication 创建学委申请
	CreateApplication(ctx context.Context, application *model.CommitteeApplication) error
	// FindApplicationByID 根据ID查找申请
	FindApplicationByID(ctx context.Context, id uint) (*model.CommitteeApplication, error)
	// FindPendingApplicationByUserID 查找用户的待审核申请
	FindPendingApplicationByUserID(ctx context.Context, userID uint) (*model.CommitteeApplication, error)
	// ListApplications 分页获取申请列表
	ListApplications(ctx context.Context, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error)
	// ListApplicationsByUserID 分页获取用户的申请列表
	ListApplicationsByUserID(ctx context.Context, userID uint, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error)
	// UpdateApplication 更新申请信息
	UpdateApplication(ctx context.Context, application *model.CommitteeApplication) error
	// UpdateApplicationStatus 更新申请状态
	UpdateApplicationStatus(ctx context.Context, id uint, status model.ApplicationStatus, reviewerID *uint, comment string) error
	// CountPendingApplications 统计待审核申请数量
	CountPendingApplications(ctx context.Context) (int64, error)
	// ExistsPendingApplication 检查用户是否有待审核的申请
	ExistsPendingApplication(ctx context.Context, userID uint) (bool, error)
}

// committeeRepository 学委申请数据访问层实现
type committeeRepository struct {
	db *gorm.DB
}

// NewCommitteeRepository 创建学委申请数据访问层实例
func NewCommitteeRepository(db *gorm.DB) CommitteeRepository {
	return &committeeRepository{db: db}
}

// CreateApplication 创建学委申请
func (r *committeeRepository) CreateApplication(ctx context.Context, application *model.CommitteeApplication) error {
	// 检查是否已有待审核的申请（只能有一个待审核申请）
	hasPending, err := r.ExistsPendingApplication(ctx, application.UserID)
	if err != nil {
		return err
	}
	if hasPending {
		return ErrApplicationAlreadyExists
	}

	// 直接创建新申请（允许用户有多个已取消/已拒绝的申请记录）
	result := r.db.WithContext(ctx).Create(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindApplicationByID 根据ID查找申请
func (r *committeeRepository) FindApplicationByID(ctx context.Context, id uint) (*model.CommitteeApplication, error) {
	var application model.CommitteeApplication
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Reviewer").
		First(&application, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, result.Error
	}
	return &application, nil
}

// FindPendingApplicationByUserID 查找用户的待审核申请
func (r *committeeRepository) FindPendingApplicationByUserID(ctx context.Context, userID uint) (*model.CommitteeApplication, error) {
	var application model.CommitteeApplication
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, model.ApplicationPending).
		First(&application)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, result.Error
	}
	return &application, nil
}

// ListApplications 分页获取申请列表
func (r *committeeRepository) ListApplications(ctx context.Context, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error) {
	var applications []*model.CommitteeApplication
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.CommitteeApplication{})

	// 如果指定了状态，添加状态过滤
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	result := query.
		Preload("User").
		Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return applications, total, nil
}

// ListApplicationsByUserID 分页获取用户的申请列表
func (r *committeeRepository) ListApplicationsByUserID(ctx context.Context, userID uint, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error) {
	var applications []*model.CommitteeApplication
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.CommitteeApplication{}).Where("user_id = ?", userID)

	// 如果指定了状态，添加状态过滤
	if status != nil {
		query = query.Where("status = ?", *status)
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
		Find(&applications)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return applications, total, nil
}

// UpdateApplication 更新申请信息
func (r *committeeRepository) UpdateApplication(ctx context.Context, application *model.CommitteeApplication) error {
	result := r.db.WithContext(ctx).Save(application)
	return result.Error
}

// UpdateApplicationStatus 更新申请状态
func (r *committeeRepository) UpdateApplicationStatus(ctx context.Context, id uint, status model.ApplicationStatus, reviewerID *uint, comment string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if reviewerID != nil {
		updates["reviewer_id"] = *reviewerID
		updates["reviewed_at"] = gorm.Expr("NOW()")
	}

	if comment != "" {
		updates["review_comment"] = comment
	}

	result := r.db.WithContext(ctx).Model(&model.CommitteeApplication{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// CountPendingApplications 统计待审核申请数量
func (r *committeeRepository) CountPendingApplications(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&model.CommitteeApplication{}).
		Where("status = ?", model.ApplicationPending).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// ExistsPendingApplication 检查用户是否有待审核的申请
func (r *committeeRepository) ExistsPendingApplication(ctx context.Context, userID uint) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&model.CommitteeApplication{}).
		Where("user_id = ? AND status = ?", userID, model.ApplicationPending).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
