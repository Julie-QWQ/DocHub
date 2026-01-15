package model

import "time"

// Announcement 公告模型
type Announcement struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Priority    string     `gorm:"type:varchar(20);not null;default:'normal';check:priority IN ('normal', 'high')" json:"priority"`
	AuthorID    uint       `gorm:"not null" json:"author_id"`
	Author      *User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	IsActive    bool       `gorm:"not null;default:true" json:"is_active"`
	PublishedAt *time.Time `json:"published_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// TableName 指定表名
func (Announcement) TableName() string {
	return "announcements"
}

// AnnouncementListRequest 公告列表请求
type AnnouncementListRequest struct {
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"page_size" json:"page_size"`
	Priority  string `form:"priority" json:"priority"`   // 优先级筛选
	IsActive  *bool  `form:"is_active" json:"is_active"` // 是否启用
	AuthorID  uint   `form:"author_id" json:"author_id"` // 发布者ID
}

// CreateAnnouncementRequest 创建公告请求
type CreateAnnouncementRequest struct {
	Title      string    `json:"title" binding:"required,min=2,max=200"`
	Content    string    `json:"content" binding:"required,min=10,max=10000"`
	Priority   string    `json:"priority" binding:"required,oneof=normal high"`
	IsActive   bool      `json:"is_active"`
	PublishedAt *time.Time `json:"published_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// UpdateAnnouncementRequest 更新公告请求
type UpdateAnnouncementRequest struct {
	Title       string     `json:"title" binding:"required,min=2,max=200"`
	Content     string     `json:"content" binding:"required,min=10,max=10000"`
	Priority    string     `json:"priority" binding:"required,oneof=normal high"`
	IsActive    bool       `json:"is_active"`
	PublishedAt *time.Time `json:"published_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// AnnouncementResponse 公告响应
type AnnouncementResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Priority    string     `json:"priority"`
	Author      *User      `json:"author"`
	IsActive    bool       `json:"is_active"`
	PublishedAt *time.Time `json:"published_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToAnnouncementResponse 转换为响应格式
func (a *Announcement) ToAnnouncementResponse() *AnnouncementResponse {
	return &AnnouncementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		Priority:    a.Priority,
		Author:      a.Author,
		IsActive:    a.IsActive,
		PublishedAt: a.PublishedAt,
		ExpiresAt:   a.ExpiresAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
