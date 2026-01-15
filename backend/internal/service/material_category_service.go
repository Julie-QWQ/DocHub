package service

import (
	"errors"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
	"gorm.io/gorm"
)

type MaterialCategoryService struct {
	categoryRepo *repository.MaterialCategoryRepository
}

func NewMaterialCategoryService(categoryRepo *repository.MaterialCategoryRepository) *MaterialCategoryService {
	return &MaterialCategoryService{
		categoryRepo: categoryRepo,
	}
}

// List 获取资料类型列表
func (s *MaterialCategoryService) List(activeOnly bool) ([]model.MaterialCategoryConfig, error) {
	return s.categoryRepo.List(activeOnly)
}

// GetByID 根据ID获取资料类型
func (s *MaterialCategoryService) GetByID(id uint) (*model.MaterialCategoryConfig, error) {
	return s.categoryRepo.GetByID(id)
}

// Create 创建资料类型
func (s *MaterialCategoryService) Create(req model.MaterialCategoryRequest) (*model.MaterialCategoryConfig, error) {
	// 检查代码是否已存在
	exists, err := s.categoryRepo.ExistsByCode(req.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("资料类型代码已存在")
	}

	// 设置默认值
	category := &model.MaterialCategoryConfig{
		Code:      req.Code,
		Name:      req.Name,
		Description: req.Description,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	} else {
		category.IsActive = true
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Update 更新资料类型
func (s *MaterialCategoryService) Update(id uint, req model.MaterialCategoryRequest) (*model.MaterialCategoryConfig, error) {
	// 获取现有资料类型
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资料类型不存在")
		}
		return nil, err
	}

	// 检查代码是否已被其他记录使用
	exists, err := s.categoryRepo.ExistsByCode(req.Code, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("资料类型代码已存在")
	}

	// 更新字段
	category.Code = req.Code
	category.Name = req.Name
	category.Description = req.Description
	category.Icon = req.Icon
	category.SortOrder = req.SortOrder

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete 删除资料类型
func (s *MaterialCategoryService) Delete(id uint) error {
	// 检查是否被使用
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("资料类型不存在")
		}
		return err
	}

	count, err := s.categoryRepo.CheckUsage(category.Code)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该资料类型正在使用中,无法删除")
	}

	return s.categoryRepo.Delete(id)
}

// ToggleStatus 切换启用状态
func (s *MaterialCategoryService) ToggleStatus(id uint) (*model.MaterialCategoryConfig, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资料类型不存在")
		}
		return nil, err
	}

	category.IsActive = !category.IsActive
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}
