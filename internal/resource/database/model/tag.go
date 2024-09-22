package model

import (
	"github.com/hcd233/Aris-blog/internal/resource/database"
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
}

// Create 创建标签
func (t *Tag) Create() (err error) {
	err = database.DB.Create(t).Error
	return
}

// Delete 删除标签
func (t *Tag) Delete() (err error) {
	err = database.DB.Delete(t).Error
	return
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
		"name": t.Name,
		"slug": t.Slug,
	}
}

// GetDetailedInfo 获取详细信息
//
//	@receiver t *Tag
//	@return map
//	@author centonhuang
//	@update 2024-09-22 03:17:50
func (t *Tag) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"name":        t.Name,
		"slug":        t.Slug,
		"description": t.Description,
	}
}

// QueryTags 查询标签
//
//	@param limit int
//	@param offset int
//	@param fields []string
//	@return tags []Tag
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 03:15:30
func QueryTags(limit int, offset int, fields []string) (tags []Tag, err error) {
	err = database.DB.Select(fields).Limit(limit).Offset(offset).Find(&tags).Error
	return
}
