package dao

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// ArticleVersionDAO 标签DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type ArticleVersionDAO struct {
	baseDAO[model.ArticleVersion]
}

// GetLatestByArticleID 通过文章ID获取最新文章版本
//
//	receiver dao *ArticleVersionDAO
//	param db *gorm.DB
//	param articleID uint
//	param fields []string
//	return articleVersion *model.ArticleVersion
//	return err error
//	author centonhuang
//	update 2024-10-17 08:14:09
func (dao *ArticleVersionDAO) GetLatestByArticleID(db *gorm.DB, articleID uint, fields, preloads []string) (articleVersion *model.ArticleVersion, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.ArticleVersion{ArticleID: articleID}).Last(&articleVersion).Error
	return
}

// GetByArticleIDAndVersion 通过文章ID和版本号获取文章版本
//
//	receiver dao *ArticleVersionDAO
//	param db *gorm.DB
//	param articleID uint
//	param version uint
//	param fields []string
//	return articleVersion *model.ArticleVersion
//	return err error
//	author centonhuang
//	update 2024-10-18 03:17:06
func (dao *ArticleVersionDAO) GetByArticleIDAndVersion(db *gorm.DB, articleID, version uint, fields, preloads []string) (articleVersion *model.ArticleVersion, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.ArticleVersion{ArticleID: articleID, Version: version}).Last(&articleVersion).Error
	return
}

// PaginateByArticleID 通过文章ID获取文章版本列表
//
//	receiver dao *ArticleVersionDAO
//	param db *gorm.DB
//	param articleID uint
//	param fields []string
//	param page int
//	param pageSize int
//	return articleVersions *[]model.ArticleVersion
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 07:08:50
func (dao *ArticleVersionDAO) PaginateByArticleID(db *gorm.DB, articleID uint, fields, preloads []string, param *PaginateParam) (articleVersions *[]model.ArticleVersion, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	if param.Query != "" && len(param.QueryFields) > 0 {
		sql = sql.Where("? LIKE ?", param.QueryFields[0], "%"+param.Query+"%")
		for _, field := range param.QueryFields[1:] {
			sql = sql.Or("? LIKE ?", field, "%"+param.Query+"%")
		}
	}

	err = sql.Where(&model.ArticleVersion{ArticleID: articleID}).Limit(limit).Offset(offset).Find(&articleVersions).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&articleVersions).Where(&model.ArticleVersion{ArticleID: articleID}).Count(&pageInfo.Total).Error
	return
}
