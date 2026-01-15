package repository

import (
	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

type MaterialCategoryRepository struct {
	db *gorm.DB
}

func NewMaterialCategoryRepository(db *gorm.DB) *MaterialCategoryRepository {
	return &MaterialCategoryRepository{db: db}
}

// List 获取资料类型列表
func (r *MaterialCategoryRepository) List(activeOnly bool) ([]model.MaterialCategoryConfig, error) {
	var categories []model.MaterialCategoryConfig

	query := r.db.Order("sort_order ASC, id ASC")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Find(&categories).Error
	return categories, err
}

// GetByID 根据ID获取资料类型
func (r *MaterialCategoryRepository) GetByID(id uint) (*model.MaterialCategoryConfig, error) {
	var category model.MaterialCategoryConfig
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByCode 根据代码获取资料类型
func (r *MaterialCategoryRepository) GetByCode(code string) (*model.MaterialCategoryConfig, error) {
	var category model.MaterialCategoryConfig
	err := r.db.Where("code = ?", code).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建资料类型
func (r *MaterialCategoryRepository) Create(category *model.MaterialCategoryConfig) error {
	return r.db.Create(category).Error
}

// Update 更新资料类型
func (r *MaterialCategoryRepository) Update(category *model.MaterialCategoryConfig) error {
	return r.db.Save(category).Error
}

// Delete 删除资料类型
func (r *MaterialCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.MaterialCategoryConfig{}, id).Error
}

// ExistsByCode 检查代码是否已存在
func (r *MaterialCategoryRepository) ExistsByCode(code string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.MaterialCategoryConfig{}).Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CheckUsage 检查资料类型是否被使用
func (r *MaterialCategoryRepository) CheckUsage(code string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Material{}).Where("category = ?", code).Count(&count).Error
	return count, err
}
