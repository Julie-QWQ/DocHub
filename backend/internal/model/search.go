package model

import (
	"time"

	"gorm.io/gorm"
)

// SearchHistory 搜索历史模型
type SearchHistory struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint   `gorm:"not null;index:idx_user_search" json:"user_id"`
	Keyword     string `gorm:"type:varchar(200);not null" json:"keyword"`          // 搜索关键词
	ResultCount int    `gorm:"not null;default:0" json:"result_count"`              // 结果数量
}

// TableName 指定表名
func (SearchHistory) TableName() string {
	return "search_histories"
}

// HotKeyword 热门搜索词模型
type HotKeyword struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Keyword      string    `gorm:"type:varchar(200);not null;uniqueIndex" json:"keyword"` // 搜索关键词
	SearchCount  int       `gorm:"not null;default:0" json:"search_count"`               // 搜索次数
	LastSearchedAt time.Time `json:"last_searched_at"`                                   // 最后搜索时间
}

// TableName 指定表名
func (HotKeyword) TableName() string {
	return "hot_keywords"
}

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Keyword    string            `form:"keyword"`                                       // 搜索关键词
	Category   *MaterialCategory `form:"category"`                                      // 分类筛选
	CourseName string            `form:"course_name"`                                   // 课程名称
	Tags       []string          `form:"tags"`                                          // 标签筛选
	Status     *MaterialStatus   `form:"status"`                                        // 状态筛选
	StartDate  string            `form:"start_date"`                                    // 开始日期
	EndDate    string            `form:"end_date"`                                      // 结束日期
	SortBy     string            `form:"sort_by,default:created_at"`                    // 排序字段: created_at, download_count, favorite_count, view_count, relevance
	SortOrder  string            `form:"sort_order,default:desc"`                       // 排序方向: asc, desc
	Page       int               `form:"page,default=1"`                                // 页码
	PageSize   int               `form:"page_size,default=20"`                          // 每页数量
}

// SearchResult 搜索结果
type SearchResult struct {
	Material   *Material `json:"material"`
	Relevance  float64   `json:"relevance"`  // 相关度分数 (0-1)
	Highlighted string   `json:"highlighted"` // 高亮显示的文本片段
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Results     []*SearchResult `json:"results"`
	Total       int64           `json:"total"`
	Page        int             `json:"page"`
	PageSize    int             `json:"page_size"`
	TotalPages  int             `json:"total_pages"`
	DidYouMean  []string        `json:"did_you_mean,omitempty"`  // 拼写建议
}

// RecommendationRequest 推荐请求参数
type RecommendationRequest struct {
	Type     string `form:"type,default:hot"`        // 推荐类型: hot, personalized, related, downloaded
	MaterialID *uint `form:"material_id"`            // 资料ID(用于相关推荐)
	Limit    int    `form:"limit,default=10"`        // 返回数量
}

// RecommendationResult 推荐结果
type RecommendationResult struct {
	Material     *Material `json:"material"`
	Reason       string    `json:"reason"`         // 推荐理由
	Score        float64   `json:"score"`          // 推荐分数
}
