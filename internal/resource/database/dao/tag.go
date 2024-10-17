package dao

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// TagDAO 标签DAO
//
//	@author centonhuang
//	@update 2024-10-17 02:30:24
type TagDAO struct {
	baseDAO[model.Tag]
}

// Delete 删除标签
//
//	@receiver dao *TagDAO
//	@param db *gorm.DB
//	@param tag *model.Tag
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 06:58:14
func (dao *TagDAO) Delete(db *gorm.DB, tag *model.Tag) (err error) {
	UUID := uuid.New().String()
	err = db.Model(tag).Updates(map[string]interface{}{"name": fmt.Sprintf("%s-%s", tag.Name, UUID), "slug": fmt.Sprintf("%s-%s", tag.Slug, UUID), "deleted_at": time.Now()}).Error
	return
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
