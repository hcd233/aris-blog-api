package model

import "gorm.io/gorm"

// Comment 评论
//
//	@author centonhuang
//	@update 2024-09-21 06:45:57
type Comment struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"column:id;primary_key;auto_increment;comment:'评论ID'"`
	ArticleID uint      `json:"article_id" gorm:"column:article_id;not null;comment:'文章ID'"`
	UserID    uint      `json:"user_id" gorm:"column:user_id;not null;comment:'用户ID'"`
	Content   string    `json:"content" gorm:"column:content;not null;comment:'评论内容'"`
	ParentID  uint      `json:"parent_id" gorm:"column:parent_id;default:NULL;comment:'父评论ID'"`
	Likes     uint      `json:"likes" gorm:"column:likes;default:0;comment:'点赞数'"`
	Parent    *Comment  `json:"parent" gorm:"foreignKey:ParentID"`
	Children  []Comment `json:"children" gorm:"foreignKey:ParentID"`
}

// GetBasicInfo 获取基本信息
//
//	@receiver c *Comment
//	@return map
//	@author centonhuang
//	@update 2024-10-24 05:45:04
func (c *Comment) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        c.ID,
		"userID":    c.UserID,
		"createdAt": c.CreatedAt,
		"parentID":  c.ParentID,
		"content":   c.Content,
		"likes":     c.Likes,
	}
}
