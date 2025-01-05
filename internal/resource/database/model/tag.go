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
	UserID      uint      `json:"user_id" gorm:"column:user_id;not null;comment:创建者ID"`
	User        *User     `json:"user" gorm:"foreignKey:UserID"`
	Articles    []Article `json:"articles" gorm:"many2many:article_tags;"`
	Likes       uint      `json:"likes" gorm:"column:likes;default:0;comment:点赞数"`
}
