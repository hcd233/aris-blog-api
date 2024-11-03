package model

import (
	"gorm.io/gorm"
)

// Tag 标签数据库模型
type Tag struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"column:id;primary_key;auto_increment;comment:标签ID"`
	Name        string    `json:"name" gorm:"column:name;unique;not null;comment:标签名称"`
	Slug        string    `json:"slug" gorm:"column:slug;unique;not null;comment:标签slug"`
	Description string    `json:"description" gorm:"column:description;comment:标签描述"`
	CreateBy    uint      `json:"create_by" gorm:"column:create_by;not null;comment:创建者ID"`
	User        *User     `json:"user" gorm:"foreignKey:CreateBy"`
	Articles    []Article `json:"articles" gorm:"many2many:article_tags;"`
	Likes       uint      `json:"likes" gorm:"column:likes;default:0;comment:点赞数"`
}

// GetBasicInfo 获取基本信息
//
//	@receiver t *Tag
//	@return map
//	@author centonhuang
//	@update 2024-09-22 03:17:33
func (t *Tag) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":   t.ID,
		"slug": t.Slug,
	}
}

// GetLikeInfo 获取点赞信息
//
//	@receiver t *Tag
//	@return map
//	@author centonhuang
//	@update 2024-11-03 08:57:20
func (t *Tag) GetLikeInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        t.ID,
		"slug":      t.Slug,
		"createdAt": t.CreatedAt,
		"creator":   t.User.GetBasicInfo(),
		"likes":     t.Likes,
	}
}

// GetDetailedInfo 获取详细信息
//
//	@receiver t *Tag
//	@return map
//	@author centonhuang
//	@update 2024-09-22 05:39:22
func (t *Tag) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"name":        t.Name,
		"slug":        t.Slug,
		"description": t.Description,
	}
}

// GetDetailedInfoWithUser 获取详细信息(包含用户名)
//
//	@receiver t *Tag
//	@return map
//	@author centonhuang
//	@update 2024-09-22 03:17:50
func (t *Tag) GetDetailedInfoWithUser() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"name":        t.Name,
		"slug":        t.Slug,
		"description": t.Description,
		"creator":     t.User.Name,
	}
}
