package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
)

var (
	// ErrFavoriteNotFound 收藏不存在
	ErrFavoriteNotFound = errors.New("收藏不存在")
	// ErrAlreadyFavorited 已收藏
	ErrAlreadyFavorited = errors.New("已收藏该资料")
)

// FavoriteService 收藏服务接口
type FavoriteService interface {
	// AddFavorite 添加收藏
	AddFavorite(ctx context.Context, userID, materialID uint) error
	// RemoveFavorite 取消收藏
	RemoveFavorite(ctx context.Context, userID, materialID uint) error
	// ListFavorites 获取用户收藏列表
	ListFavorites(ctx context.Context, userID uint, page, pageSize int) ([]*model.FavoriteResponse, int64, error)
	// IsFavorited 检查是否已收藏
	IsFavorited(ctx context.Context, userID, materialID uint) (bool, error)
}

// favoriteService 收藏服务实现
type favoriteService struct {
	favoriteRepo repository.FavoriteRepository
	materialRepo repository.MaterialRepository
}

// NewFavoriteService 创建收藏服务实例
func NewFavoriteService(
	favoriteRepo repository.FavoriteRepository,
	materialRepo repository.MaterialRepository,
) FavoriteService {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
		materialRepo: materialRepo,
	}
}

// AddFavorite 添加收藏
func (s *favoriteService) AddFavorite(ctx context.Context, userID, materialID uint) error {
	// 检查资料是否存在
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil {
		if errors.Is(err, repository.ErrMaterialNotFound) {
			return ErrMaterialNotFound
		}
		return fmt.Errorf("获取资料失败: %w", err)
	}

	// 只能收藏已审核通过的资料
	if material.Status != model.StatusApproved {
		return ErrAccessDenied
	}

	// 检查是否已收藏
	exists, err := s.favoriteRepo.Exists(ctx, userID, materialID)
	if err != nil {
		return fmt.Errorf("检查收藏状态失败: %w", err)
	}
	if exists {
		return ErrAlreadyFavorited
	}

	// 创建收藏
	favorite := &model.Favorite{
		UserID:     userID,
		MaterialID: materialID,
	}

	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		return fmt.Errorf("添加收藏失败: %w", err)
	}

	// 增加收藏次数
	if err := s.materialRepo.IncrementFavoriteCount(ctx, materialID); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("增加收藏次数失败: %v\n", err)
	}

	return nil
}

// RemoveFavorite 取消收藏
func (s *favoriteService) RemoveFavorite(ctx context.Context, userID, materialID uint) error {
	// 删除收藏
	if err := s.favoriteRepo.Delete(ctx, userID, materialID); err != nil {
		if errors.Is(err, repository.ErrFavoriteNotFound) {
			return ErrFavoriteNotFound
		}
		return fmt.Errorf("取消收藏失败: %w", err)
	}

	// 减少收藏次数
	if err := s.materialRepo.DecrementFavoriteCount(ctx, materialID); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("减少收藏次数失败: %v\n", err)
	}

	return nil
}

// ListFavorites 获取用户收藏列表
func (s *favoriteService) ListFavorites(ctx context.Context, userID uint, page, pageSize int) ([]*model.FavoriteResponse, int64, error) {
	favorites, total, err := s.favoriteRepo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收藏列表失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*model.FavoriteResponse, 0, len(favorites))
	for _, favorite := range favorites {
		response := &model.FavoriteResponse{
			ID:         favorite.ID,
			UserID:     favorite.UserID,
			MaterialID: favorite.MaterialID,
			CreatedAt:  favorite.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if favorite.Material != nil {
			response.Material = favorite.Material.ToMaterialResponse()
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// IsFavorited 检查是否已收藏
func (s *favoriteService) IsFavorited(ctx context.Context, userID, materialID uint) (bool, error) {
	exists, err := s.favoriteRepo.Exists(ctx, userID, materialID)
	if err != nil {
		return false, fmt.Errorf("检查收藏状态失败: %w", err)
	}
	return exists, nil
}
