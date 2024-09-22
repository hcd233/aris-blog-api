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
	Name     string     `json:"name" gorm:"column:name;not null;unique;comment:类别名称"`
	ParentID uint      `json:"parent_id" gorm:"column:parent_id;comment:父类别ID"`
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

// GetChildren 获取子类别
//
//	@receiver c *Category
//	@return children []Category, err error
func (c *Category) GetChildren() (children []Category, err error) {
	err = database.DB.Where("parent_id = ?", c.ID).Find(&children).Error
	return
}
