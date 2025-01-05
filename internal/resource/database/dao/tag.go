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
//	@param userID uint
//	@param fields []string
//	@return tag *model.Tag
//	@return err error
//	@author centonhuang
//	@update 2024-10-23 12:53:25
func (dao *TagDAO) GetBySlug(db *gorm.DB, slug string, fields, preloads []string) (tag *model.Tag, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(model.Tag{Slug: slug}).First(&tag).Error
	return
}

// PaginateByUserID 通过用户ID获取标签
//
//	@receiver dao *TagDAO
//	@param db *gorm.DB
//	@param userID uint
//	@param fields []string
//	@param page int
//	@param pageSize int
//	@return tags *[]model.Tag
//	@return pageInfo *PageInfo
//	@return err error
//	@author centonhuang
//	@update 2024-11-01 07:10:06
func (dao *TagDAO) PaginateByUserID(db *gorm.DB, userID uint, fields, preloads []string, page, pageSize int) (tags *[]model.Tag, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(model.Tag{UserID: userID}).Limit(limit).Offset(offset).Find(&tags).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Model(&tags).Where(model.Tag{UserID: userID}).Count(&pageInfo.Total).Error
	return
}
