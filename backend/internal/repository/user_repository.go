package repository

import (
	"context"
	"errors"

	"github.com/study-upc/backend/internal/model"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound 用户不存在错误
	ErrUserNotFound = errors.New("用户不存在")
	// ErrUserAlreadyExists 用户已存在错误
	ErrUserAlreadyExists = errors.New("用户已存在")
)

// UserRepository 用户数据访问层接口
type UserRepository interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *model.User) error
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*model.User, error)
	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *model.User) error
	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error
	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, userID uint) error
	// ListUsers 分页获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// CountByRole 统计指定角色的用户数量
	CountByRole(role model.UserRole) (int64, error)
}

// userRepository 用户数据访问层实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问层实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser 创建用户
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		// 检查是否是唯一索引冲突
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Save(user)
	return result.Error
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("password_hash", hashedPassword)
	return result.Error
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("NOW()"))
	return result.Error
}

// ListUsers 分页获取用户列表
func (r *userRepository) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	result := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// CountByRole 统计指定角色的用户数量
func (r *userRepository) CountByRole(role model.UserRole) (int64, error) {
	var count int64
	result := r.db.Model(&model.User{}).Where("role = ? AND deleted_at IS NULL", role).Count(&count)
	return count, result.Error
}
