package model

import (
	"time"

	"gorm.io/gorm"
)

// UserView 用户浏览
//
//	author centonhuang
//	update 2024-11-01 07:34:14
type UserView struct {
	gorm.Model
	UserID       uint      `json:"user_id" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:用户ID"`
	User         *User     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ArticleID    uint      `json:"article_id" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:文章ID"`
	Article      *Article  `json:"article" gorm:"foreignKey:ArticleID;references:ID"`
	LastViewedAt time.Time `json:"last_viewed_at" gorm:"not null;comment:最后浏览时间"`
	Progress     int8      `json:"progress" gorm:"not null;comment:浏览进度"`
}
