package service

import (
	"context"
	"fmt"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
)

var (
	// ErrNotificationNotFound 通知不存在
	ErrNotificationNotFound = fmt.Errorf("通知不存在")
	// ErrNotificationAccessDenied 无权访问该通知
	ErrNotificationAccessDenied = fmt.Errorf("无权访问该通知")
)

// NotificationService 通知服务接口
type NotificationService interface {
	// CreateNotification 创建通知
	CreateNotification(ctx context.Context, notification *model.Notification) error
	// ListNotifications 获取用户通知列表
	ListNotifications(ctx context.Context, userID uint, page, pageSize int) ([]*model.Notification, int64, error)
	// GetUnreadNotifications 获取未读通知
	GetUnreadNotifications(ctx context.Context, userID uint) ([]*model.Notification, error)
	// MarkAsRead 标记通知为已读
	MarkAsRead(ctx context.Context, notificationID, userID uint) error
	// MarkAllAsRead 标记所有通知为已读
	MarkAllAsRead(ctx context.Context, userID uint) error
	// GetUnreadCount 获取未读通知数量
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	// DeleteNotification 删除通知
	DeleteNotification(ctx context.Context, notificationID, userID uint) error
	// BroadcastNotification 广播通知（发送给所有用户或特定角色）
	BroadcastNotification(ctx context.Context, title, content string, notifyType model.NotificationType, role *model.UserRole) error
}

// notificationService 通知服务实现
type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

// NewNotificationService 创建通知服务实例
func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

// CreateNotification 创建通知
func (s *notificationService) CreateNotification(ctx context.Context, notification *model.Notification) error {
	if err := s.notificationRepo.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("创建通知失败: %w", err)
	}
	return nil
}

// ListNotifications 获取用户通知列表
func (s *notificationService) ListNotifications(ctx context.Context, userID uint, page, pageSize int) ([]*model.Notification, int64, error) {
	notifications, total, err := s.notificationRepo.ListNotifications(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取通知列表失败: %w", err)
	}
	return notifications, total, nil
}

// GetUnreadNotifications 获取未读通知
func (s *notificationService) GetUnreadNotifications(ctx context.Context, userID uint) ([]*model.Notification, error) {
	notifications, err := s.notificationRepo.ListUnreadNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取未读通知失败: %w", err)
	}
	return notifications, nil
}

// MarkAsRead 标记通知为已读
func (s *notificationService) MarkAsRead(ctx context.Context, notificationID, userID uint) error {
	if err := s.notificationRepo.MarkAsRead(ctx, notificationID, userID); err != nil {
		return fmt.Errorf("标记通知已读失败: %w", err)
	}
	return nil
}

// MarkAllAsRead 标记所有通知为已读
func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uint) error {
	if err := s.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("标记所有通知已读失败: %w", err)
	}
	return nil
}

// GetUnreadCount 获取未读通知数量
func (s *notificationService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	count, err := s.notificationRepo.CountUnreadNotifications(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("获取未读通知数量失败: %w", err)
	}
	return count, nil
}

// DeleteNotification 删除通知
func (s *notificationService) DeleteNotification(ctx context.Context, notificationID, userID uint) error {
	if err := s.notificationRepo.DeleteNotification(ctx, notificationID, userID); err != nil {
		return fmt.Errorf("删除通知失败: %w", err)
	}
	return nil
}

// BroadcastNotification 广播通知（发送给所有用户或特定角色）
func (s *notificationService) BroadcastNotification(ctx context.Context, title, content string, notifyType model.NotificationType, role *model.UserRole) error {
	// 如果指定了角色，只发送给该角色的用户
	// 否则发送给所有用户
	// 这里简化处理，实际可能需要批量插入优化
	var users []*model.User
	var err error

	if role != nil {
		// 获取特定角色的用户（这里需要 user repository 支持按角色查询）
		// 暂时简化处理
		return fmt.Errorf("按角色广播功能待实现")
	} else {
		// 获取所有用户
		users, _, err = s.userRepo.ListUsers(ctx, 1, 10000)
		if err != nil {
			return fmt.Errorf("获取用户列表失败: %w", err)
		}
	}

	// 为每个用户创建通知
	for _, user := range users {
		notification := &model.Notification{
			UserID:  user.ID,
			Type:    notifyType,
			Title:   title,
			Content: content,
			Status:  model.NotifyUnread,
		}
		if err := s.notificationRepo.CreateNotification(ctx, notification); err != nil {
			// 记录错误但继续处理其他用户
			continue
		}
	}

	return nil
}
