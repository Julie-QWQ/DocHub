package model

import (
	"time"

	"gorm.io/gorm"
)

// ApplicationStatus 申请状态
type ApplicationStatus string

const (
	ApplicationPending   ApplicationStatus = "pending"   // 待审核
	ApplicationApproved  ApplicationStatus = "approved"  // 已通过
	ApplicationRejected  ApplicationStatus = "rejected"  // 已拒绝
	ApplicationCancelled ApplicationStatus = "cancelled" // 已取消
)

// CommitteeApplication 学委申请模型
type CommitteeApplication struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID   uint              `gorm:"not null;index" json:"user_id"`                              // 申请人ID
	User     *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`                    // 用户信息
	Status   ApplicationStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`  // 申请状态
	Reason   string            `gorm:"type:text;not null" json:"reason"`                           // 申请理由
	ReviewerID *uint           `gorm:"index" json:"reviewer_id,omitempty"`                         // 审核人ID
	Reviewer  *User            `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`            // 审核人信息
	ReviewedAt *time.Time      `json:"reviewed_at,omitempty"`                                      // 审核时间
	ReviewComment string       `gorm:"type:text" json:"review_comment,omitempty"`                  // 审核意见
}

// TableName 指定表名
func (CommitteeApplication) TableName() string {
	return "committee_applications"
}
