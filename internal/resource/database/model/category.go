package model

import (
	"gorm.io/gorm"
)

// Category 文章类别数据库模型
//
//	@author centonhuang
//	@update 2024-09-22 10:00:00
type Category struct {
	gorm.Model
	ID       uint       `json:"id" gorm:"column:id;primary_key;auto_increment;comment:类别ID"`
	Name     string     `json:"name" gorm:"column:name;not null;uniqueIndex:pid_name;comment:类别名称"`
	ParentID uint       `json:"parent_id" gorm:"column:parent_id;default:NULL;uniqueIndex:pid_name;comment:父类别ID"`
	Parent   *Category  `json:"parent" gorm:"foreignKey:ParentID"`
	UserID   uint       `json:"user_id" gorm:"column:user_id;comment:用户ID"`
	User     *User      `json:"user" gorm:"foreignKey:UserID"`
	Children []Category `json:"children" gorm:"foreignKey:ParentID"`
	Articles []Article  `json:"articles" gorm:"foreignKey:CategoryID"`
}
