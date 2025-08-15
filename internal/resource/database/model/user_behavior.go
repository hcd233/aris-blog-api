package model

import "gorm.io/gorm"

// UserBehavior 用户行为模型
//
//	author system
//	update 2025-01-19 12:00:00
type UserBehavior struct {
	gorm.Model
	UserID       uint                   `gorm:"index;not null" json:"userId"`                 // 用户ID
	ItemID       uint                   `gorm:"index;not null" json:"itemId"`                 // 物品ID
	ItemType     string                 `gorm:"index;not null;size:20" json:"itemType"`       // 物品类型(article/tag)
	BehaviorType string                 `gorm:"index;not null;size:20" json:"behaviorType"`   // 行为类型
	Score        float64                `gorm:"default:0" json:"score"`                       // 评分
	Weight       float64                `gorm:"default:1.0" json:"weight"`                    // 权重
	Context      string                 `gorm:"type:text" json:"context"`                     // 上下文信息(JSON)
	SessionID    string                 `gorm:"index;size:100" json:"sessionId"`              // 会话ID
	IPAddress    string                 `gorm:"size:45" json:"ipAddress"`                     // IP地址
	UserAgent    string                 `gorm:"size:500" json:"userAgent"`                    // 用户代理
	Timestamp    int64                  `gorm:"index;not null" json:"timestamp"`              // 时间戳

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 表名
func (UserBehavior) TableName() string {
	return "user_behaviors"
}