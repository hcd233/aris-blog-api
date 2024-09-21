package model

import "gorm.io/gorm"

// Comment 评论
//
//	@author centonhuang
//	@update 2024-09-21 06:45:57
type Comment struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"column:id;primary_key;auto_increment"`
	ArticleID uint   `json:"article_id" gorm:"column:article_id;not null"`
	UserID    uint   `json:"user_id" gorm:"column:user_id;not null"`
	Content   string `json:"content" gorm:"column:content;not null"`
	ParentID  uint   `json:"parent_id" gorm:"column:parent_id"`
	Likes     uint   `json:"likes" gorm:"column:likes;default:0"`
}
