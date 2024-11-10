package dao

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// ArticleDAO 标签DAO
//
//	@author centonhuang
//	@update 2024-10-17 06:34:00
type ArticleDAO struct {
	baseDAO[model.Article]
}

// Delete 删除文章
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param article *model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 06:52:28
func (dao *ArticleDAO) Delete(db *gorm.DB, article *model.Article) (err error) {
	UUID := uuid.New().String()
	err = db.Model(article).Updates(map[string]interface{}{"slug": fmt.Sprintf("%s-%s", article.Slug, UUID), "deleted_at": time.Now()}).Error
	return
}

// GetBySlugAndUserID 通过slug获取文章
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param slug string
//	@param userID uint
//	@param fields []string
//	@return article *model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 07:17:59
func (dao *ArticleDAO) GetBySlugAndUserID(db *gorm.DB, slug string, userID uint, fields, preloads []string) (article *model.Article, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Article{Slug: slug, UserID: userID}).First(&article).Error
	return
}

// PaginateByUserID 通过用户ID获取文章列表
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param userID uint
//	@param fields []string
//	@param page int
//	@param pageSize int
//	@return articles *[]model.Article
//	@return pageInfo *PageInfo
//	@return err error
//	@author centonhuang
//	@update 2024-11-01 07:09:20
func (dao *ArticleDAO) PaginateByUserID(db *gorm.DB, userID uint, fields, preloads []string, page, pageSize int) (articles *[]model.Article, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Article{UserID: userID}).Limit(limit).Offset(offset).Find(&articles).Error

	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Where(&model.Article{UserID: userID}).Count(&pageInfo.Total).Error
	return
}

// PaginateByCategoryID 通过类别ID获取文章列表
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param categoryID uint
//	@param fields []string
//	@param page int
//	@param pageSize int
//	@return articles *[]model.Article
//	@return pageInfo *PageInfo
//	@return err error
//	@author centonhuang
//	@update 2024-11-01 07:09:26
func (dao *ArticleDAO) PaginateByCategoryID(db *gorm.DB, categoryID uint, fields []string, preloads []string, page, pageSize int) (articles *[]model.Article, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Article{CategoryID: categoryID}).Limit(limit).Offset(offset).Find(&articles).Error

	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Where(&model.Article{CategoryID: categoryID}).Count(&pageInfo.Total).Error

	return
}

// PaginateByPublished 列出已发布的文章
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param fields []string
//	@param page int
//	@param pageSize int
//	@return articles *[]model.Article
//	@return pageInfo *PageInfo
//	@return err error
//	@author centonhuang
//	@update 2024-11-01 05:33:37
func (dao *ArticleDAO) PaginateByPublished(db *gorm.DB, fields []string, page, pageSize int) (articles *[]model.Article, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize
	err = db.Select(fields).Where(&model.Article{Status: model.ArticleStatusPublish}).Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Where(&model.Article{Status: model.ArticleStatusPublish}).Count(&pageInfo.Total).Error
	return
}

// BatchGetAllByIDs 批量获取文章
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param ids []uint
//	@return articles *[]model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-11-03 07:35:58
func (dao *ArticleDAO) BatchGetAllByIDs(db *gorm.DB, ids []uint) (articles *[]model.Article, err error) {
	err = db.Preload("Category").Preload("Tags").Preload("User").Where("id IN ?", ids).Find(&articles).Error
	return
}
