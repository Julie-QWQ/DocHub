package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/email"
	"github.com/study-upc/backend/internal/pkg/verification"
	"github.com/study-upc/backend/internal/repository"
)

// EmailVerificationService 邮箱验证服务接口
type EmailVerificationService interface {
	// SendVerificationCode 发送验证码
	SendVerificationCode(ctx context.Context, email, purpose string) error

	// VerifyCode 验证验证码
	VerifyCode(ctx context.Context, email, code, purpose string) error

	// RegisterWithEmailCode 使用邮箱验证码注册
	RegisterWithEmailCode(ctx context.Context, username, email, password, code string) error

	// LoginWithEmailCode 使用邮箱验证码登录
	LoginWithEmailCode(ctx context.Context, email, password, code string) (string, *model.User, error)

	// CleanExpiredCodes 清理过期验证码
	CleanExpiredCodes(ctx context.Context) (int64, error)
}

type emailVerificationService struct {
	userRepo    repository.UserRepository
	emailRepo   repository.EmailVerificationRepository
	emailClient *email.SMTPClient
}

// NewEmailVerificationService 创建邮箱验证服务
func NewEmailVerificationService(
	userRepo repository.UserRepository,
	emailRepo repository.EmailVerificationRepository,
	emailClient *email.SMTPClient,
) EmailVerificationService {
	return &emailVerificationService{
		userRepo:    userRepo,
		emailRepo:   emailRepo,
		emailClient: emailClient,
	}
}

// 发送验证码
func (s *emailVerificationService) SendVerificationCode(ctx context.Context, emailAddr, purpose string) error {
	// 去除邮箱首尾空格
	emailAddr = strings.TrimSpace(emailAddr)

	if err := s.emailRepo.DeleteByEmail(ctx, emailAddr); err != nil {
		return fmt.Errorf("删除旧验证码失败: %w", err)
	}

	// 生成新验证码
	code, err := verification.GenerateVerificationCode()
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	// 保存到数据库
	expireTime := time.Now().Add(10 * time.Minute)
	verificationCode := &model.EmailVerificationCode{
		Email:     emailAddr,
		Code:      code,
		ExpiresAt: expireTime,
		Purpose:   purpose,
		IsUsed:    false,
	}

	if err := s.emailRepo.Create(ctx, verificationCode); err != nil {
		return fmt.Errorf("保存验证码失败: %w", err)
	}

	// 发送邮件
	if err := s.emailClient.SendVerificationCode(emailAddr, code, purpose); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// 验证验证码
func (s *emailVerificationService) VerifyCode(ctx context.Context, emailAddr, code, purpose string) error {
	// 查询验证码
	verificationCode, err := s.emailRepo.GetByEmail(ctx, emailAddr, purpose)
	if err != nil {
		return fmt.Errorf("查询验证码失败: %w", err)
	}

	if verificationCode == nil {
		return errors.New("验证码不存在或已过期")
	}

	// 检查验证码是否匹配
	if verificationCode.Code != code {
		return errors.New("验证码错误")
	}

	// 检查是否已使用
	if verificationCode.IsUsed {
		return errors.New("验证码已使用")
	}

	// 检查是否过期
	if time.Now().After(verificationCode.ExpiresAt) {
		return errors.New("验证码已过期")
	}

	// 标记为已使用
	if err := s.emailRepo.MarkAsUsed(ctx, verificationCode.ID); err != nil {
		return fmt.Errorf("标记验证码失败: %w", err)
	}

	return nil
}

// 使用邮箱验证码注册
func (s *emailVerificationService) RegisterWithEmailCode(ctx context.Context, username, emailAddr, password, code string) error {
	// 先验证验证码
	if err := s.VerifyCode(ctx, emailAddr, code, "register"); err != nil {
		return err
	}

	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingEmail, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("查询邮箱失败: %w", err)
	}
	if existingEmail != nil {
		return errors.New("邮箱已被注册")
	}

	// 创建用户(密码加密在Repository层处理)
	user := &model.User{
		Username:      username,
		Email:         emailAddr,
		PasswordHash:  password, // Repository会进行hash处理
		Role:          model.RoleStudent,
		Status:        model.StatusActive,
		EmailVerified: true, // 邮箱已验证
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// 使用邮箱和密码登录
func (s *emailVerificationService) LoginWithEmailCode(ctx context.Context, emailAddr, password, code string) (string, *model.User, error) {
	// 先验证验证码
	if err := s.VerifyCode(ctx, emailAddr, code, "login"); err != nil {
		return "", nil, err
	}

	// 查询用户
	user, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		return "", nil, fmt.Errorf("查询用户失败: %w", err)
	}

	if user == nil {
		return "", nil, errors.New("用户不存在")
	}

	// 如果提供了密码，验证密码；如果密码为空，则跳过密码验证（邮箱验证码登录）
	if password != "" {
		// 验证密码(Repository层处理)
		if err := s.userRepo.VerifyPassword(ctx, user.ID, password); err != nil {
			return "", nil, errors.New("密码错误")
		}
	}

	// 检查账号状态
	if user.Status == model.StatusBanned {
		return "", nil, ErrUserDisabled
	}
	if user.Status == model.StatusInactive {
		return "", nil, ErrUserInactive
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", nil, fmt.Errorf("更新登录时间失败: %w", err)
	}

	// 生成JWT token(在Handler层处理)
	return "", user, nil
}

// 清理过期验证码
func (s *emailVerificationService) CleanExpiredCodes(ctx context.Context) (int64, error) {
	return s.emailRepo.CleanExpired(ctx)
}
