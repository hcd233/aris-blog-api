package model

import "gorm.io/gorm"

// Tag 标签数据库模型
type Tag struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"column:id;primary_key;auto_increment;uniqueIndex:idx_name_slug"`
	Name        string    `json:"name" gorm:"column:name;not null;uniqueIndex:idx_name_slug"`
	Slug        string    `json:"slug" gorm:"column:slug;unique;not null"`
	Description string    `json:"description" gorm:"column:description"`
	Articles    []Article `json:"articles" gorm:"many2many:article_tags;"`
}
