package model

import "gorm.io/gorm"

// UserProfile 用户画像模型
//
//	author system
//	update 2025-01-19 12:00:00
type UserProfile struct {
	gorm.Model
	UserID        uint   `gorm:"uniqueIndex;not null" json:"userId"`         // 用户ID
	Preferences   string `gorm:"type:text" json:"preferences"`               // 偏好标签及权重(JSON)
	Interests     string `gorm:"type:text" json:"interests"`                 // 兴趣类别(JSON)
	BehaviorStats string `gorm:"type:text" json:"behaviorStats"`             // 行为统计(JSON)
	Metadata      string `gorm:"type:text" json:"metadata"`                  // 元数据(JSON)
	LastUpdated   int64  `gorm:"index;not null" json:"lastUpdated"`          // 最后更新时间
	Version       int    `gorm:"default:1" json:"version"`                   // 版本号

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 表名
func (UserProfile) TableName() string {
	return "user_profiles"
}