package service

import (
	"context"
	"errors"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
)

var (
	// ErrAnnouncementNotFound 公告不存在
	ErrAnnouncementNotFound = errors.New("公告不存在")
	// ErrInvalidExpiresAt 过期时间无效
	ErrInvalidExpiresAt = errors.New("过期时间必须晚于发布时间")
)

// AnnouncementService 公告服务接口
type AnnouncementService interface {
	// CreateAnnouncement 创建公告
	CreateAnnouncement(ctx context.Context, authorID uint, req *model.CreateAnnouncementRequest) (*model.AnnouncementResponse, error)
	// GetAnnouncement 获取公告详情
	GetAnnouncement(ctx context.Context, id uint) (*model.AnnouncementResponse, error)
	// GetActiveAnnouncements 获取启用的公告列表（用于首页公告栏）
	GetActiveAnnouncements(ctx context.Context, limit int) ([]model.Announcement, error)
	// ListAnnouncements 查询公告列表（支持筛选、分页）
	ListAnnouncements(ctx context.Context, req *model.AnnouncementListRequest) ([]model.Announcement, int64, error)
	// UpdateAnnouncement 更新公告
	UpdateAnnouncement(ctx context.Context, id, authorID uint, req *model.UpdateAnnouncementRequest) (*model.AnnouncementResponse, error)
	// DeleteAnnouncement 删除公告
	DeleteAnnouncement(ctx context.Context, id, authorID uint) error
}

type announcementService struct {
	announcementRepo repository.AnnouncementRepository
	userRepo         repository.UserRepository
}

// NewAnnouncementService 创建公告服务实例
func NewAnnouncementService(
	announcementRepo repository.AnnouncementRepository,
	userRepo repository.UserRepository,
) AnnouncementService {
	return &announcementService{
		announcementRepo: announcementRepo,
		userRepo:         userRepo,
	}
}

// CreateAnnouncement 创建公告
func (s *announcementService) CreateAnnouncement(ctx context.Context, authorID uint, req *model.CreateAnnouncementRequest) (*model.AnnouncementResponse, error) {
	// 验证发布者是否存在
	_, err := s.userRepo.FindByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	// 验证过期时间必须晚于发布时间
	if req.PublishedAt != nil && req.ExpiresAt != nil {
		if req.ExpiresAt.Before(*req.PublishedAt) || req.ExpiresAt.Equal(*req.PublishedAt) {
			return nil, ErrInvalidExpiresAt
		}
	}

	// 创建公告
	announcement := &model.Announcement{
		Title:       req.Title,
		Content:     req.Content,
		Priority:    req.Priority,
		AuthorID:    authorID,
		IsActive:    req.IsActive,
		PublishedAt: req.PublishedAt,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := s.announcementRepo.Create(ctx, announcement); err != nil {
		return nil, err
	}

	// 重新加载包含作者信息的公告
	announcement, err = s.announcementRepo.FindByID(ctx, announcement.ID)
	if err != nil {
		return nil, err
	}

	return announcement.ToAnnouncementResponse(), nil
}

// GetAnnouncement 获取公告详情
func (s *announcementService) GetAnnouncement(ctx context.Context, id uint) (*model.AnnouncementResponse, error) {
	announcement, err := s.announcementRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAnnouncementNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, err
	}

	return announcement.ToAnnouncementResponse(), nil
}

// GetActiveAnnouncements 获取启用的公告列表（用于首页公告栏）
func (s *announcementService) GetActiveAnnouncements(ctx context.Context, limit int) ([]model.Announcement, error) {
	if limit <= 0 {
		limit = 5 // 默认显示5条
	}
	if limit > 20 {
		limit = 20 // 最多显示20条
	}

	return s.announcementRepo.FindActive(ctx, limit)
}

// ListAnnouncements 查询公告列表（支持筛选、分页）
func (s *announcementService) ListAnnouncements(ctx context.Context, req *model.AnnouncementListRequest) ([]model.Announcement, int64, error) {
	return s.announcementRepo.List(ctx, req)
}

// UpdateAnnouncement 更新公告
func (s *announcementService) UpdateAnnouncement(ctx context.Context, id, authorID uint, req *model.UpdateAnnouncementRequest) (*model.AnnouncementResponse, error) {
	// 查找公告
	announcement, err := s.announcementRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAnnouncementNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, err
	}

	// 权限检查：只有作者和管理员可以修改
	// TODO: 添加管理员权限检查
	if announcement.AuthorID != authorID {
		return nil, errors.New("无权修改此公告")
	}

	// 验证过期时间必须晚于发布时间
	if req.PublishedAt != nil && req.ExpiresAt != nil {
		if req.ExpiresAt.Before(*req.PublishedAt) || req.ExpiresAt.Equal(*req.PublishedAt) {
			return nil, ErrInvalidExpiresAt
		}
	}

	// 更新字段
	announcement.Title = req.Title
	announcement.Content = req.Content
	announcement.Priority = req.Priority
	announcement.IsActive = req.IsActive
	announcement.PublishedAt = req.PublishedAt
	announcement.ExpiresAt = req.ExpiresAt

	if err := s.announcementRepo.Update(ctx, announcement); err != nil {
		return nil, err
	}

	// 重新加载包含最新信息的公告
	announcement, err = s.announcementRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return announcement.ToAnnouncementResponse(), nil
}

// DeleteAnnouncement 删除公告
func (s *announcementService) DeleteAnnouncement(ctx context.Context, id, authorID uint) error {
	// 查找公告
	announcement, err := s.announcementRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAnnouncementNotFound) {
			return ErrAnnouncementNotFound
		}
		return err
	}

	// 权限检查：只有作者和管理员可以删除
	// TODO: 添加管理员权限检查
	if announcement.AuthorID != authorID {
		return errors.New("无权删除此公告")
	}

	return s.announcementRepo.Delete(ctx, id)
}
