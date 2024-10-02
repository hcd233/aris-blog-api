package model

import (
	"time"

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
		"slug": t.Slug,
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

// QueryTags 查询标签
//
//	@param limit int
//	@param offset int
//	@param fields []string
//	@return tags []Tag
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 03:15:30
func QueryTags(limit int, offset int, fields []string) (tags *[]Tag, err error) {
	err = database.DB.Select(fields).Limit(limit).Offset(offset).Find(&tags).Error
	return
}

// QueryTagBySlug 查询标签
//
//	@param tagSlug string
//	@param fields []string
//	@return tag Tag
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 04:45:22
func QueryTagBySlug(tagSlug string, fields []string) (tag *Tag, err error) {
	err = database.DB.Preload("User").Select(fields).Where(&Tag{Slug: tagSlug}).First(&tag).Error
	return
}

// QueryTagsBySlugs 查询多个标签
//
//	@param tagSlugs []string
//	@param fields []string
//	@return tags []Tag
//	@return err error
//	@author centonhuang
//	@update 2024-10-02 01:05:34
func QueryTagsBySlugs(tagSlugs []string, fields []string) (tags *[]Tag, err error) {
	err = database.DB.Preload("User").Select(fields).Where("slug IN ?", tagSlugs).Find(&tags).Error
	return
}

// UpdateTagBySlug 更新标签
//
//	@param tagSlug string
//	@param info map[string]interface{}
//	@return tag Tag
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 04:10:07
func UpdateTagBySlug(tagSlug string, info map[string]interface{}) (tag *Tag, err error) {
	info["updated_at"] = time.Now()
	err = database.DB.Model(&Tag{}).Where(&Tag{Slug: tagSlug}).Updates(info).Error
	if err != nil {
		return
	}
	err = database.DB.Preload("User").Where(&Tag{Slug: tagSlug}).First(&tag).Error
	return
}
