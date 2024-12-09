package model

import (
	"time"

	"gorm.io/gorm"
)

// UserView 用户浏览
//
//	@author centonhuang
//	@update 2024-11-01 07:34:14
type UserView struct {
	gorm.Model
	UserID       uint      `json:"user_id" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:用户ID"`
	User         *User     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ArticleID    uint      `json:"article_id" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:文章ID"`
	Article      *Article  `json:"article" gorm:"foreignKey:ArticleID;references:ID"`
	LastViewedAt time.Time `json:"last_viewed_at" gorm:"not null;default:CURRENT_TIMESTAMP(3);comment:最后浏览时间"`
	Progress     int8      `json:"progress" gorm:"not null;comment:浏览进度"`
}

// GetBasicInfo 获取用户浏览基本信息
//
//	@receiver uv *UserView
//	@return map
//	@author centonhuang
//	@update 2024-12-09 16:11:07
func (uv *UserView) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":           uv.ID,
		"progress":     uv.Progress,
		"lastViewedAt": uv.LastViewedAt,
	}
}

// GetDetailedInfo 获取用户浏览详细信息
//
//	@receiver uv *UserView
//	@return map
//	@author centonhuang
//	@update 2024-12-09 16:11:11
func (uv *UserView) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":           uv.ID,
		"viewer":       uv.User.GetBasicInfo(),
		"article":      uv.Article.GetViewInfo(),
		"progress":     uv.Progress,
		"lastViewedAt": uv.LastViewedAt,
	}
}
