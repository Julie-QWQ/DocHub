package repository

import (
	"context"
	"errors"
	"time"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

// EmailVerificationRepository 邮箱验证码仓储接口
type EmailVerificationRepository interface {
	// Create 创建验证码记录
	Create(ctx context.Context, code *model.EmailVerificationCode) error

	// GetByEmail 获取指定邮箱的验证码
	GetByEmail(ctx context.Context, email string, purpose string) (*model.EmailVerificationCode, error)

	// Delete 删除验证码记录
	Delete(ctx context.Context, id uint) error

	// DeleteByEmail 删除指定邮箱和用途的所有验证码
	DeleteByEmail(ctx context.Context, email string) error

	// MarkAsUsed 标记验证码已使用
	MarkAsUsed(ctx context.Context, id uint) error

	// CleanExpired 清理过期的验证码
	CleanExpired(ctx context.Context) (int64, error)
}

type emailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository 创建邮箱验证码仓储
func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, code *model.EmailVerificationCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *emailVerificationRepository) GetByEmail(ctx context.Context, email string, purpose string) (*model.EmailVerificationCode, error) {
	var code model.EmailVerificationCode
	err := r.db.WithContext(ctx).
		Where("email = ? AND purpose = ? AND is_used = ? AND expires_at > ?",
			email, purpose, false, time.Now()).
		Order("created_at DESC").
		First(&code).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &code, nil
}

func (r *emailVerificationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&model.EmailVerificationCode{}, id).Error
}

func (r *emailVerificationRepository) DeleteByEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Where("email = ?", email).
		Delete(&model.EmailVerificationCode{}).Error
}

func (r *emailVerificationRepository) MarkAsUsed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.EmailVerificationCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		}).Error
}

func (r *emailVerificationRepository) CleanExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("expires_at < ?", time.Now()).
		Delete(&model.EmailVerificationCode{})

	return result.RowsAffected, result.Error
}
