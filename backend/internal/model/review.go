package model

import (
	"time"

	"gorm.io/gorm"
)

// ReviewAction 审核操作类型
type ReviewAction string

const (
	ReviewApprove ReviewAction = "approve" // 通过
	ReviewReject  ReviewAction = "reject"  // 拒绝
)

// ReviewTarget 审核目标类型
type ReviewTarget string

const (
	TargetMaterial     ReviewTarget = "material"      // 资料
	TargetCommittee    ReviewTarget = "committee"     // 学委申请
	TargetReport       ReviewTarget = "report"        // 举报
)

// ReviewRecord 审核记录模型
type ReviewRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ReviewerID    uint         `gorm:"not null;index" json:"reviewer_id"`                       // 审核人ID
	Reviewer      *User        `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`         // 审核人信息
	TargetType    ReviewTarget `gorm:"type:varchar(20);not null;index" json:"target_type"`      // 目标类型
	TargetID      uint         `gorm:"not null;index" json:"target_id"`                         // 目标ID
	Action        ReviewAction `gorm:"type:varchar(20);not null" json:"action"`                 // 审核操作
	Comment       string       `gorm:"type:text" json:"comment,omitempty"`                      // 审核意见
	OriginalData  string       `gorm:"type:jsonb" json:"original_data,omitempty"`              // 原始数据（快照）
}

// TableName 指定表名
func (ReviewRecord) TableName() string {
	return "review_records"
}

// ReviewerStatistics 审核人统计信息
type ReviewerStatistics struct {
	TotalReviews     int64 `json:"total_reviews"`      // 总审核数
	ApprovedCount    int64 `json:"approved_count"`     // 通过数
	RejectedCount    int64 `json:"rejected_count"`     // 拒绝数
	MaterialReviews  int64 `json:"material_reviews"`   // 资料审核数
	CommitteeReviews int64 `json:"committee_reviews"`  // 学委申请审核数
	ReportReviews    int64 `json:"report_reviews"`     // 举报处理数
}
