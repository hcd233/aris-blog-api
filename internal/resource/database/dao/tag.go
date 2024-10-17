package dao

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// TagDAO 标签数据访问对象
//
//	@author centonhuang
//	@update 2024-10-17 02:30:24
type TagDAO struct {
	baseDAO[model.Tag]
}

// GetBySlug 通过slug获取标签
//
//	@receiver dao *TagDAO
//	@param db *gorm.DB
//	@param slug string
//	@param fields []string
//	@return tag *model.Tag
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 05:58:06
func (dao *TagDAO) GetBySlug(db *gorm.DB, slug string, fields []string) (tag *model.Tag, err error) {
	err = db.Select(fields).Where(model.Tag{Slug: slug}).First(&tag).Error
	return
}
