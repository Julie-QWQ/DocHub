package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationType 通知类型
type NotificationType string

const (
	NotifySystem      NotificationType = "system"       // 系统通知
	NotifyMaterial    NotificationType = "material"     // 资料审核通知
	NotifyCommittee   NotificationType = "committee"    // 学委申请通知
	NotifyReport      NotificationType = "report"       // 举报处理通知
)

// NotificationStatus 通知状态
type NotificationStatus string

const (
	NotifyUnread NotificationStatus = "unread" // 未读
	NotifyRead   NotificationStatus = "read"   // 已读
)

// Notification 通知模型
type Notification struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID  uint               `gorm:"not null;index" json:"user_id"`                           // 接收用户ID
	User    *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`                 // 用户信息
	Type    NotificationType   `gorm:"type:varchar(20);not null;index" json:"type"`             // 通知类型
	Title   string             `gorm:"type:varchar(255);not null" json:"title"`                 // 通知标题
	Content string             `gorm:"type:text;not null" json:"content"`                       // 通知内容
	Status  NotificationStatus `gorm:"type:varchar(20);not null;default:'unread'" json:"status"` // 通知状态
	Link    string             `gorm:"type:varchar(255)" json:"link,omitempty"`                 // 相关链接
	ReadAt  *time.Time         `json:"read_at,omitempty"`                                       // 阅读时间
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}
