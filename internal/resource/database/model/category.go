package model

import (
	"github.com/hcd233/Aris-blog/internal/resource/database"
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
	UserID   uint       `json:"user_id" gorm:"column:user_id;comment:用户ID"`
	Parent   *Category  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children" gorm:"foreignKey:ParentID"`
	Articles []Article  `json:"articles" gorm:"foreignKey:CategoryID"`
}

// Create 创建类别
//
//	@receiver c *Category
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 10:05:00
func (c *Category) Create() (err error) {
	err = database.DB.Create(c).Error
	return
}

// Delete 删除类别
//
//	@receiver c *Category
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 10:10:00
func (c *Category) Delete() (err error) {
	err = database.DB.Delete(c).Error
	return
}

// GetBasicInfo 获取基本信息
//
//	@receiver c * Category
//	@return map
//	@author centonhuang
//	@update 2024-09-28 07:09:27
func (c *Category) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":       c.ID,
		"name":     c.Name,
		"parentID": c.ParentID,
	}
}

// QueryChildren 获取子类别
//
//	@receiver c *Category
//	@return children []Category, err error
func (c *Category) QueryChildren() (children []Category, err error) {
	err = database.DB.Where(&Category{ParentID: c.ID}).Find(&children).Error
	return
}

// QueryParent 获取父类别
//
//	@receiver c *Category
//	@return parent Category
//	@return err error
//	@author centonhuang
//	@update 2024-09-28 07:08:23
func (c *Category) QueryParent() (parent Category, err error) {
	err = database.DB.Where(&Category{ID: c.ParentID}).First(&parent).Error
	return
}

// QueryCategoryByID 使用ID查询指定类别
//
//	@param categoryID uint
//	@param fields []string
//	@return category Category
//	@return err error
//	@author centonhuang
//	@update 2024-10-01 05:01:44
func QueryCategoryByID(categoryID uint, fields []string) (category Category, err error) {
	err = database.DB.Select(fields).Where(&Category{ID: categoryID}).First(&category).Error
	return
}

// QueryRootCategoriesByUserID 查询指定用户的梗类别
//
//	@param userID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return categories []Category
//	@return err error
//	@author centonhuang
//	@update 2024-10-01 03:55:57
func QueryRootCategoriesByUserID(userID uint, fields []string, limit, offset int) (categories []Category, err error) {
	err = database.DB.Select(fields).Where(&Category{UserID: userID}).Where("parent_id IS NULL").Limit(limit).Offset(offset).Find(&categories).Error
	return
}

// QueryChildrenCategoriesByUserID 查询指定用户的子类别
//
//	@param parentID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return categories []Category
//	@return err error
//	@author centonhuang
//	@update 2024-10-01 05:11:22
func QueryChildrenCategoriesByUserID(parentID uint, fields []string, limit, offset int) (categories []Category, err error) {
	err = database.DB.Select(fields).Where(&Category{ParentID: parentID}).Limit(limit).Offset(offset).Find(&categories).Error
	return
}
