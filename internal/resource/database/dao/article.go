package dao

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// ArticleDAO 标签DAO
//
//	author centonhuang
//	update 2024-10-17 06:34:00
type ArticleDAO struct {
	baseDAO[model.Article]
}

// Delete 删除文章
//
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param article *model.Article
//	return err error
//	author centonhuang
//	update 2024-10-17 06:52:28
func (dao *ArticleDAO) Delete(db *gorm.DB, article *model.Article) (err error) {
	UUID := uuid.New().String()
	err = db.Model(article).Updates(map[string]interface{}{"slug": fmt.Sprintf("%s-%s", article.Slug, UUID), "deleted_at": time.Now().UTC()}).Error
	return
}

// GetByIDAndStatus 通过ID和状态获取文章
//
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param slug string
//	param userID uint
//	param fields []string
//	return article *model.Article
//	return err error
//	author centonhuang
//	update 2024-10-17 07:17:59
func (dao *ArticleDAO) GetByIDAndStatus(db *gorm.DB, articleID uint, status model.ArticleStatus, fields, preloads []string) (article *model.Article, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Article{ID: articleID, Status: status}).First(&article).Error
	return
}

// GetByIDAndUserID 通过ID和用户ID获取文章
//
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param articleID uint
//	param userID uint
//	param fields []string
//	param preloads []string
//	return article *model.Article
//	return err error
//	author centonhuang
//	update 2025-01-18 17:00:00
func (dao *ArticleDAO) GetByIDAndUserID(db *gorm.DB, articleID uint, userID uint, fields, preloads []string) (article *model.Article, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Article{ID: articleID, UserID: userID}).First(&article).Error
	return
}

// GetBySlugAndUserID 通过Slug和用户ID获取文章
//
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param slug string
//	param userID uint
//	param fields []string
//	param preloads []string
//	return article *model.Article
//	return err error
//	author centonhuang
//	update 2025-01-19 15:23:26
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
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param userID uint
//	param fields []string
//	param page int
//	param pageSize int
//	return articles *[]model.Article
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 07:09:20
func (dao *ArticleDAO) PaginateByUserID(db *gorm.DB, userID uint, fields, preloads []string, param *PaginateParam) (articles *[]model.Article, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	
	// 添加模糊查询支持
	if param.Query != "" && len(param.QueryFields) > 0 {
		sql = sql.Where("? LIKE ?", param.QueryFields[0], "%"+param.Query+"%")
		for _, field := range param.QueryFields[1:] {
			sql = sql.Or("? LIKE ?", field, "%"+param.Query+"%")
		}
	}
	
	err = sql.Where(&model.Article{UserID: userID}).Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&articles).Where(&model.Article{UserID: userID}).Count(&pageInfo.Total).Error
	return
}

// PaginateByCategoryID 通过类别ID获取文章列表
//
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param categoryID uint
//	param fields []string
//	param page int
//	param pageSize int
//	return articles *[]model.Article
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 07:09:26
func (dao *ArticleDAO) PaginateByCategoryID(db *gorm.DB, categoryID uint, fields []string, preloads []string, param *PaginateParam) (articles *[]model.Article, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	
	// 添加模糊查询支持
	if param.Query != "" && len(param.QueryFields) > 0 {
		sql = sql.Where("? LIKE ?", param.QueryFields[0], "%"+param.Query+"%")
		for _, field := range param.QueryFields[1:] {
			sql = sql.Or("? LIKE ?", field, "%"+param.Query+"%")
		}
	}
	
	err = sql.Where(&model.Article{CategoryID: categoryID}).Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&articles).Where(&model.Article{CategoryID: categoryID}).Count(&pageInfo.Total).Error

	return
}

// PaginateByStatus 列出已发布的文章
//
//	receiver dao *ArticleDAO
//	param db *gorm.DB
//	param status model.ArticleStatus
//	param fields []string
//	param preloads []string
//	param page int
//	param pageSize int
//	return articles *[]model.Article
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 05:33:37
func (dao *ArticleDAO) PaginateByStatus(db *gorm.DB, status model.ArticleStatus, fields []string, preloads []string, param *PaginateParam) (articles *[]model.Article, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	
	// 添加模糊查询支持
	if param.Query != "" && len(param.QueryFields) > 0 {
		sql = sql.Where("? LIKE ?", param.QueryFields[0], "%"+param.Query+"%")
		for _, field := range param.QueryFields[1:] {
			sql = sql.Or("? LIKE ?", field, "%"+param.Query+"%")
		}
	}
	
	err = sql.Where(&model.Article{Status: status}).Limit(limit).Offset(offset).Find(&articles).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&articles).Where(&model.Article{Status: model.ArticleStatusPublish}).Count(&pageInfo.Total).Error
	return
}
