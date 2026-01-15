package repository

import (
	"context"
	"errors"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

var (
	// ErrNotificationNotFound 通知不存在
	ErrNotificationNotFound = errors.New("通知不存在")
)

// NotificationRepository 通知数据访问层接口
type NotificationRepository interface {
	// CreateNotification 创建通知
	CreateNotification(ctx context.Context, notification *model.Notification) error
	// FindNotificationByID 根据ID查找通知
	FindNotificationByID(ctx context.Context, id uint) (*model.Notification, error)
	// ListNotifications 分页获取用户的通知列表
	ListNotifications(ctx context.Context, userID uint, page, pageSize int) ([]*model.Notification, int64, error)
	// ListUnreadNotifications 获取用户的未读通知列表
	ListUnreadNotifications(ctx context.Context, userID uint) ([]*model.Notification, error)
	// MarkAsRead 标记通知为已读
	MarkAsRead(ctx context.Context, id uint, userID uint) error
	// MarkAllAsRead 标记用户的所有通知为已读
	MarkAllAsRead(ctx context.Context, userID uint) error
	// CountUnreadNotifications 统计用户未读通知数量
	CountUnreadNotifications(ctx context.Context, userID uint) (int64, error)
	// DeleteNotification 删除通知
	DeleteNotification(ctx context.Context, id uint, userID uint) error
}

// notificationRepository 通知数据访问层实现
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 创建通知数据访问层实例
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// CreateNotification 创建通知
func (r *notificationRepository) CreateNotification(ctx context.Context, notification *model.Notification) error {
	result := r.db.WithContext(ctx).Create(notification)
	return result.Error
}

// FindNotificationByID 根据ID查找通知
func (r *notificationRepository) FindNotificationByID(ctx context.Context, id uint) (*model.Notification, error) {
	var notification model.Notification
	result := r.db.WithContext(ctx).
		Preload("User").
		First(&notification, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, result.Error
	}
	return &notification, nil
}

// ListNotifications 分页获取用户的通知列表
func (r *notificationRepository) ListNotifications(ctx context.Context, userID uint, page, pageSize int) ([]*model.Notification, int64, error) {
	var notifications []*model.Notification
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	result := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return notifications, total, nil
}

// ListUnreadNotifications 获取用户的未读通知列表
func (r *notificationRepository) ListUnreadNotifications(ctx context.Context, userID uint) ([]*model.Notification, error) {
	var notifications []*model.Notification
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, model.NotifyUnread).
		Order("created_at DESC").
		Find(&notifications)

	if result.Error != nil {
		return nil, result.Error
	}

	return notifications, nil
}

// MarkAsRead 标记通知为已读
func (r *notificationRepository) MarkAsRead(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"status":   model.NotifyRead,
			"read_at": gorm.Expr("NOW()"),
		})
	return result.Error
}

// MarkAllAsRead 标记用户的所有通知为已读
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND status = ?", userID, model.NotifyUnread).
		Updates(map[string]interface{}{
			"status":   model.NotifyRead,
			"read_at": gorm.Expr("NOW()"),
		})
	return result.Error
}

// CountUnreadNotifications 统计用户未读通知数量
func (r *notificationRepository) CountUnreadNotifications(ctx context.Context, userID uint) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND status = ?", userID, model.NotifyUnread).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// DeleteNotification 删除通知
func (r *notificationRepository) DeleteNotification(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Notification{})
	return result.Error
}
