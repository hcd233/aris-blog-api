package model

import "gorm.io/gorm"

// RecommendationLog 推荐日志模型
//
//	author system
//	update 2025-01-19 12:00:00
type RecommendationLog struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null" json:"userId"`               // 用户ID
	RequestType string `gorm:"index;not null;size:20" json:"requestType"`  // 请求类型(article/tag)
	Algorithm   string `gorm:"index;not null;size:50" json:"algorithm"`    // 算法名称
	ItemIDs     string `gorm:"type:text" json:"itemIds"`                   // 推荐物品ID列表(JSON)
	Scores      string `gorm:"type:text" json:"scores"`                    // 推荐分数列表(JSON)
	Context     string `gorm:"type:text" json:"context"`                   // 推荐上下文(JSON)
	ClickedIDs  string `gorm:"type:text" json:"clickedIds"`                // 被点击的物品ID(JSON)
	Timestamp   int64  `gorm:"index;not null" json:"timestamp"`            // 推荐时间戳
	SessionID   string `gorm:"index;size:100" json:"sessionId"`            // 会话ID

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 表名
func (RecommendationLog) TableName() string {
	return "recommendation_logs"
}