package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
)

var (
	// ErrAlreadyReviewed 已经审核过
	ErrAlreadyReviewed = errors.New("已经审核过")
)

// ReviewService 审核服务接口
type ReviewService interface {
	// ReviewMaterial 审核资料
	ReviewMaterial(ctx context.Context, materialID, reviewerID uint, approved bool, comment string) error
	// HandleReport 处理举报
	HandleReport(ctx context.Context, reportID, handlerID uint, approved bool, note string) error
	// GetReviewHistory 获取审核历史
	GetReviewHistory(ctx context.Context, targetType model.ReviewTarget, targetID uint, page, pageSize int) ([]*model.ReviewRecord, int64, error)
	// GetReviewerStatistics 获取审核人统计信息
	GetReviewerStatistics(ctx context.Context, reviewerID uint) (*model.ReviewerStatistics, error)
	// SetNotificationService 设置通知服务
	SetNotificationService(notificationSvc NotificationService)
}

// reviewService 审核服务实现
type reviewService struct {
	materialRepo      repository.MaterialRepository
	committeeRepo     repository.CommitteeRepository
	reportRepo        repository.ReportRepository
	reviewRepo        repository.ReviewRepository
	userRepo          repository.UserRepository
	notificationSvc   NotificationService
	notificationSvcSet bool // 标记通知服务是否已设置
}

// NewReviewService 创建审核服务实例（通知服务可选）
func NewReviewService(
	materialRepo repository.MaterialRepository,
	committeeRepo repository.CommitteeRepository,
	reportRepo repository.ReportRepository,
	reviewRepo repository.ReviewRepository,
	userRepo repository.UserRepository,
) ReviewService {
	return &reviewService{
		materialRepo:  materialRepo,
		committeeRepo: committeeRepo,
		reportRepo:    reportRepo,
		reviewRepo:    reviewRepo,
		userRepo:      userRepo,
	}
}

// SetNotificationService 设置通知服务（用于解决循环依赖）
func (s *reviewService) SetNotificationService(notificationSvc NotificationService) {
	s.notificationSvc = notificationSvc
	s.notificationSvcSet = true
}

// ReviewMaterial 审核资料
func (s *reviewService) ReviewMaterial(ctx context.Context, materialID, reviewerID uint, approved bool, comment string) error {
	// 获取资料信息
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		return fmt.Errorf("获取资料信息失败: %w", err)
	}

	// 检查资料状态
	if material.Status != model.StatusPending {
		return ErrMaterialAlreadyReviewed
	}

	// 确定新状态
	var newStatus model.MaterialStatus
	var action model.ReviewAction
	if approved {
		newStatus = model.StatusApproved
		action = model.ReviewApprove
	} else {
		newStatus = model.StatusRejected
		action = model.ReviewReject
	}

	// 更新资料审核状态
	rejectionReason := ""
	if !approved && comment != "" {
		rejectionReason = comment
	}
	if err := s.materialRepo.UpdateReviewStatus(ctx, materialID, newStatus, &reviewerID, rejectionReason); err != nil {
		return fmt.Errorf("更新资料状态失败: %w", err)
	}

	// 创建审核记录
	reviewRecord := &model.ReviewRecord{
		ReviewerID:   reviewerID,
		TargetType:   model.TargetMaterial,
		TargetID:     materialID,
		Action:       action,
		Comment:      comment,
		OriginalData: fmt.Sprintf(`{"title": "%s", "category": "%s", "status": "%s"}`, material.Title, material.Category, material.Status),
	}
	if err := s.reviewRepo.CreateReviewRecord(ctx, reviewRecord); err != nil {
		return fmt.Errorf("创建审核记录失败: %w", err)
	}

	// 发送通知给上传者
	s.sendReviewNotification(ctx, material.UploaderID, "资料审核结果",
		fmt.Sprintf("您上传的资料《%s》已%s", material.Title, getReviewStatusText(approved)), model.NotifyMaterial, "")

	return nil
}

// HandleReport 处理举报
func (s *reviewService) HandleReport(ctx context.Context, reportID, handlerID uint, approved bool, note string) error {
	// 获取举报信息
	report, err := s.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("获取举报信息失败: %w", err)
	}

	// 检查举报状态
	if report.Status != model.ReportStatusPending {
		return ErrAlreadyReviewed
	}

	// 确定新状态
	var newStatus model.ReportStatus
	if approved {
		newStatus = model.ReportStatusApproved
	} else {
		newStatus = model.ReportStatusRejected
	}

	// 更新举报状态
	report.Status = newStatus
	report.HandlerID = &handlerID
	report.HandleNote = note
	if err := s.reportRepo.Update(ctx, report); err != nil {
		return fmt.Errorf("更新举报状态失败: %w", err)
	}

	// 如果举报通过，需要处理被举报的资料（例如删除或标记）
	if approved {
		// 这里可以根据业务需求处理被举报的资料
		// 例如：软删除资料
		if err := s.materialRepo.Delete(ctx, report.MaterialID); err != nil {
			return fmt.Errorf("处理被举报资料失败: %w", err)
		}
	}

	// 创建审核记录
	reviewRecord := &model.ReviewRecord{
		ReviewerID: handlerID,
		TargetType: model.TargetReport,
		TargetID:   reportID,
		Action:     model.ReviewApprove,
		Comment:    note,
		OriginalData: fmt.Sprintf(`{"material_id": %d, "reason": "%s", "description": "%s"}`,
			report.MaterialID, report.Reason, report.Description),
	}
	if err := s.reviewRepo.CreateReviewRecord(ctx, reviewRecord); err != nil {
		return fmt.Errorf("创建审核记录失败: %w", err)
	}

	// 发送通知给举报人
	s.sendReviewNotification(ctx, report.UserID, "举报处理结果",
		fmt.Sprintf("您对资料的举报已%s", getReviewStatusText(approved)), model.NotifyReport, "")

	return nil
}

// GetReviewHistory 获取审核历史
func (s *reviewService) GetReviewHistory(ctx context.Context, targetType model.ReviewTarget, targetID uint, page, pageSize int) ([]*model.ReviewRecord, int64, error) {
	records, total, err := s.reviewRepo.ListReviewRecords(ctx, page, pageSize, &targetType, &targetID)
	if err != nil {
		return nil, 0, fmt.Errorf("获取审核历史失败: %w", err)
	}
	return records, total, nil
}

// GetReviewerStatistics 获取审核人统计信息
func (s *reviewService) GetReviewerStatistics(ctx context.Context, reviewerID uint) (*model.ReviewerStatistics, error) {
	records, _, err := s.reviewRepo.ListReviewRecordsByReviewerID(ctx, reviewerID, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("获取审核记录失败: %w", err)
	}

	stats := &model.ReviewerStatistics{
		TotalReviews: int64(len(records)),
	}

	for _, record := range records {
		switch record.Action {
		case model.ReviewApprove:
			stats.ApprovedCount++
		case model.ReviewReject:
			stats.RejectedCount++
		}

		switch record.TargetType {
		case model.TargetMaterial:
			stats.MaterialReviews++
		case model.TargetCommittee:
			stats.CommitteeReviews++
		case model.TargetReport:
			stats.ReportReviews++
		}
	}

	return stats, nil
}

// sendReviewNotification 发送审核通知
func (s *reviewService) sendReviewNotification(ctx context.Context, userID uint, title, content string, notifyType model.NotificationType, link string) {
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
