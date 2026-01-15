package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
)

var (
	// ErrHasPendingApplication 已有待审核申请
	ErrHasPendingApplication = errors.New("已存在待审核的学委申请")
	// ErrAlreadyCommittee 已经是学委
	ErrAlreadyCommittee = errors.New("您已经是学委")
	// ErrApplicationNotPending 申请不是待审核状态
	ErrApplicationNotPending = errors.New("该申请不是待审核状态")
	// ErrCannotReviewOwnApplication 不能审核自己的申请
	ErrCannotReviewOwnApplication = errors.New("不能审核自己的申请")
)

// CommitteeService 学委申请服务接口
type CommitteeService interface {
	// ApplyForCommittee 申请学委
	ApplyForCommittee(ctx context.Context, userID uint, reason string) (*model.CommitteeApplication, error)
	// ListMyApplications 获取我的申请列表
	ListMyApplications(ctx context.Context, userID uint, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error)
	// GetApplication 获取申请详情
	GetApplication(ctx context.Context, applicationID uint) (*model.CommitteeApplication, error)
	// ListApplications 获取申请列表（管理员）
	ListApplications(ctx context.Context, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error)
	// ReviewApplication 审核学委申请
	ReviewApplication(ctx context.Context, applicationID, reviewerID uint, approved bool, comment string) error
	// GetPendingCount 获取待审核申请数量
	GetPendingCount(ctx context.Context) (int64, error)
	// CancelApplication 取消申请
	CancelApplication(ctx context.Context, applicationID, userID uint) error
	// SetNotificationService 设置通知服务
	SetNotificationService(notificationSvc NotificationService)
}

// committeeService 学委申请服务实现
type committeeService struct {
	committeeRepo      repository.CommitteeRepository
	userRepo           repository.UserRepository
	reviewRepo         repository.ReviewRepository
	notificationSvc    NotificationService
	notificationSvcSet bool // 标记通知服务是否已设置
}

// NewCommitteeService 创建学委申请服务实例（通知服务可选）
func NewCommitteeService(
	committeeRepo repository.CommitteeRepository,
	userRepo repository.UserRepository,
	reviewRepo repository.ReviewRepository,
) CommitteeService {
	return &committeeService{
		committeeRepo:   committeeRepo,
		userRepo:        userRepo,
		reviewRepo:      reviewRepo,
	}
}

// SetNotificationService 设置通知服务（用于解决循环依赖）
func (s *committeeService) SetNotificationService(notificationSvc NotificationService) {
	s.notificationSvc = notificationSvc
	s.notificationSvcSet = true
}

// ApplyForCommittee 申请学委
func (s *committeeService) ApplyForCommittee(ctx context.Context, userID uint, reason string) (*model.CommitteeApplication, error) {
	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 检查是否已经是学委
	if user.Role == model.RoleCommittee {
		return nil, ErrAlreadyCommittee
	}

	// 检查是否已有待审核申请
	hasPending, err := s.committeeRepo.ExistsPendingApplication(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("检查申请状态失败: %w", err)
	}
	if hasPending {
		return nil, ErrHasPendingApplication
	}

	// 创建申请
	application := &model.CommitteeApplication{
		UserID: userID,
		Status: model.ApplicationPending,
		Reason: reason,
	}

	if err := s.committeeRepo.CreateApplication(ctx, application); err != nil {
		return nil, fmt.Errorf("创建申请失败: %w", err)
	}

	return application, nil
}

// ListMyApplications 获取我的申请列表
func (s *committeeService) ListMyApplications(ctx context.Context, userID uint, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error) {
	applications, total, err := s.committeeRepo.ListApplicationsByUserID(ctx, userID, page, pageSize, status)
	if err != nil {
		return nil, 0, fmt.Errorf("获取申请列表失败: %w", err)
	}
	return applications, total, nil
}

// GetApplication 获取申请详情
func (s *committeeService) GetApplication(ctx context.Context, applicationID uint) (*model.CommitteeApplication, error) {
	application, err := s.committeeRepo.FindApplicationByID(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("获取申请详情失败: %w", err)
	}
	return application, nil
}

// ListApplications 获取申请列表（管理员）
func (s *committeeService) ListApplications(ctx context.Context, page, pageSize int, status *model.ApplicationStatus) ([]*model.CommitteeApplication, int64, error) {
	applications, total, err := s.committeeRepo.ListApplications(ctx, page, pageSize, status)
	if err != nil {
		return nil, 0, fmt.Errorf("获取申请列表失败: %w", err)
	}
	return applications, total, nil
}

// ReviewApplication 审核学委申请
func (s *committeeService) ReviewApplication(ctx context.Context, applicationID, reviewerID uint, approved bool, comment string) error {
	// 获取申请信息
	application, err := s.committeeRepo.FindApplicationByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("获取申请信息失败: %w", err)
	}

	// 检查申请状态
	if application.Status != model.ApplicationPending {
		return ErrApplicationNotPending
	}

	// 不能审核自己的申请
	if application.UserID == reviewerID {
		return ErrCannotReviewOwnApplication
	}

	// 确定新状态
	var newStatus model.ApplicationStatus
	var action model.ReviewAction
	if approved {
		newStatus = model.ApplicationApproved
		action = model.ReviewApprove
	} else {
		newStatus = model.ApplicationRejected
		action = model.ReviewReject
	}

	// 更新申请状态
	if err := s.committeeRepo.UpdateApplicationStatus(ctx, applicationID, newStatus, &reviewerID, comment); err != nil {
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	// 如果通过，更新用户角色
	if approved {
		if err := s.updateUserRole(ctx, application.UserID, model.RoleCommittee); err != nil {
			return fmt.Errorf("更新用户角色失败: %w", err)
		}
	}

	// 创建审核记录
	reviewRecord := &model.ReviewRecord{
		ReviewerID:   reviewerID,
		TargetType:   model.TargetCommittee,
		TargetID:     applicationID,
		Action:       action,
		Comment:      comment,
		OriginalData: fmt.Sprintf(`{"user_id": %d, "reason": "%s"}`, application.UserID, application.Reason),
	}
	if err := s.reviewRepo.CreateReviewRecord(ctx, reviewRecord); err != nil {
		return fmt.Errorf("创建审核记录失败: %w", err)
	}

	// 发送通知
	s.sendReviewNotification(ctx, application.UserID, "学委申请审核结果",
		fmt.Sprintf("您的学委申请已%s", getReviewStatusText(approved)), model.NotifyCommittee, "")

	return nil
}

// GetPendingCount 获取待审核申请数量
func (s *committeeService) GetPendingCount(ctx context.Context) (int64, error) {
	count, err := s.committeeRepo.CountPendingApplications(ctx)
	if err != nil {
		return 0, fmt.Errorf("获取待审核申请数量失败: %w", err)
	}
	return count, nil
}

// CancelApplication 取消申请
func (s *committeeService) CancelApplication(ctx context.Context, applicationID, userID uint) error {
	// 获取申请信息
	application, err := s.committeeRepo.FindApplicationByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("获取申请信息失败: %w", err)
	}

	// 检查是否是申请人
	if application.UserID != userID {
		return ErrAccessDenied
	}

	// 检查申请状态
	if application.Status != model.ApplicationPending {
		return ErrApplicationNotPending
	}

	// 更新状态为已取消
	if err := s.committeeRepo.UpdateApplicationStatus(ctx, applicationID, model.ApplicationCancelled, nil, ""); err != nil {
		return fmt.Errorf("取消申请失败: %w", err)
	}

	return nil
}

// updateUserRole 更新用户角色
func (s *committeeService) updateUserRole(ctx context.Context, userID uint, role model.UserRole) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Role = role
	return s.userRepo.UpdateUser(ctx, user)
}

// sendReviewNotification 发送审核通知
func (s *committeeService) sendReviewNotification(ctx context.Context, userID uint, title, content string, notifyType model.NotificationType, link string) {
	notification := &model.Notification{
		UserID:  userID,
		Type:    notifyType,
		Title:   title,
		Content: content,
		Status:  model.NotifyUnread,
		Link:    link,
	}
	_ = s.notificationSvc.CreateNotification(ctx, notification)
}

// getReviewStatusText 获取审核状态文本
func getReviewStatusText(approved bool) string {
	if approved {
		return "通过"
	}
	return "被拒绝"
}
