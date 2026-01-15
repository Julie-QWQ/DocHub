package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
	"gorm.io/gorm"
)

// ReportService 举报服务接口
type ReportService interface {
	// CreateReport 创建举报
	CreateReport(ctx context.Context, userID, materialID uint, req *model.ReportRequest) error
	// HandleReport 处理举报
	HandleReport(ctx context.Context, reportID, handlerID uint, req *model.HandleReportRequest) error
	// ListReports 获取举报列表
	ListReports(ctx context.Context, page, pageSize int, status *model.ReportStatus) ([]*model.ReportResponse, int64, error)
	// GetReport 获取举报详情
	GetReport(ctx context.Context, reportID uint) (*model.ReportResponse, error)
}

// reportService 举报服务实现
type reportService struct {
	reportRepo   repository.ReportRepository
	materialRepo repository.MaterialRepository
}

// NewReportService 创建举报服务实例
func NewReportService(
	reportRepo repository.ReportRepository,
	materialRepo repository.MaterialRepository,
) ReportService {
	return &reportService{
		reportRepo:   reportRepo,
		materialRepo: materialRepo,
	}
}

// CreateReport 创建举报
func (s *reportService) CreateReport(ctx context.Context, userID, materialID uint, req *model.ReportRequest) error {
	// 检查资料是否存在
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return ErrMaterialNotFound
		}
		return fmt.Errorf("获取资料失败: %w", err)
	}

	// 只能举报已审核通过的资料
	if material.Status != model.StatusApproved {
		return ErrAccessDenied
	}

	// 检查是否已有待处理的举报
	_, err = s.reportRepo.FindPendingByUserAndMaterial(ctx, userID, materialID)
	if err == nil {
		// 已有待处理的举报
		return errors.New("已有待处理的举报，请勿重复举报")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查举报状态失败: %w", err)
	}

	// 创建举报
	report := &model.Report{
		UserID:      userID,
		MaterialID:  materialID,
		Reason:      req.Reason,
		Description: req.Description,
		Status:      model.ReportStatusPending,
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return fmt.Errorf("创建举报失败: %w", err)
	}

	return nil
}

// HandleReport 处理举报
func (s *reportService) HandleReport(ctx context.Context, reportID, handlerID uint, req *model.HandleReportRequest) error {
	// 获取举报
	report, err := s.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		if errors.Is(err, repository.ErrReportNotFound) {
			return repository.ErrReportNotFound
		}
		return fmt.Errorf("获取举报失败: %w", err)
	}

	// 检查状态
	if report.Status != model.ReportStatusPending {
		return errors.New("该举报已被处理")
	}

	// 更新举报状态
	report.Status = req.Status
	report.HandlerID = &handlerID
	report.HandleNote = req.HandleNote

	if err := s.reportRepo.Update(ctx, report); err != nil {
		return fmt.Errorf("处理举报失败: %w", err)
	}

	// 如果举报通过，删除资料
	if req.Status == model.ReportStatusApproved {
		if err := s.materialRepo.Delete(ctx, report.MaterialID); err != nil {
			// 记录错误但继续处理
			fmt.Printf("删除被举报资料失败: %v\n", err)
		}
	}

	return nil
}

// ListReports 获取举报列表
func (s *reportService) ListReports(ctx context.Context, page, pageSize int, status *model.ReportStatus) ([]*model.ReportResponse, int64, error) {
	reports, total, err := s.reportRepo.List(ctx, page, pageSize, status)
	if err != nil {
		return nil, 0, fmt.Errorf("获取举报列表失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*model.ReportResponse, 0, len(reports))
	for _, report := range reports {
		response := s.convertToResponse(report)
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetReport 获取举报详情
func (s *reportService) GetReport(ctx context.Context, reportID uint) (*model.ReportResponse, error) {
	report, err := s.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		if errors.Is(err, repository.ErrReportNotFound) {
			return nil, repository.ErrReportNotFound
		}
		return nil, fmt.Errorf("获取举报失败: %w", err)
	}

	return s.convertToResponse(report), nil
}

// convertToResponse 转换为响应格式
func (s *reportService) convertToResponse(report *model.Report) *model.ReportResponse {
	response := &model.ReportResponse{
		ID:          report.ID,
		UserID:      report.UserID,
		MaterialID:  report.MaterialID,
		Reason:      report.Reason,
		Description: report.Description,
		Status:      report.Status,
		HandlerID:   report.HandlerID,
		HandleNote:  report.HandleNote,
		CreatedAt:   report.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 转换举报人信息
	if report.User != nil {
		userInfo := report.User.ToUserInfo()
		response.Reporter = &userInfo
	}

	// 转换资料信息
	if report.Material != nil {
		response.Material = report.Material.ToMaterialResponse()
	}

	// 转换处理人信息
	if report.Handler != nil {
		handlerInfo := report.Handler.ToUserInfo()
		response.Handler = &handlerInfo
	}

	// 转换处理时间
	if report.HandledAt != nil {
		handledAt := report.HandledAt.Format("2006-01-02 15:04:05")
		response.HandledAt = &handledAt
	}

	return response
}
