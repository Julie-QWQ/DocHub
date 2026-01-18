package model

import (
	"time"

	"gorm.io/gorm"
)

// EmailVerificationCode 邮箱验证码模型
type EmailVerificationCode struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`        // 邮箱地址
	Code         string    `gorm:"type:varchar(10);not null" json:"-"`                         // 验证码(不返回给前端)
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`                                 // 过期时间
	IsUsed       bool      `gorm:"default:false" json:"is_used"`                               // 是否已使用
	UsedAt       *time.Time `json:"used_at,omitempty"`                                         // 使用时间
	Purpose      string    `gorm:"type:varchar(20);not null" json:"purpose"`                   // 用途: register/login/reset_password
}

// TableName 指定表名
func (EmailVerificationCode) TableName() string {
	return "email_verification_codes"
}
