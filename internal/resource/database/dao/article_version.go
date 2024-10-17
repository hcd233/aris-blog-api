package dao

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// ArticleVersionDAO 标签DAO
//
//	@author centonhuang
//	@update 2024-10-17 02:30:24
type ArticleVersionDAO struct {
	baseDAO[model.ArticleVersion]
}

// GetLatestByArticleID 通过文章ID获取最新文章版本
//
//	@receiver dao *ArticleVersionDAO
//	@param db *gorm.DB
//	@param articleID uint
//	@param fields []string
//	@return articleVersion *model.ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 08:14:09
func (dao *ArticleVersionDAO) GetLatestByArticleID(db *gorm.DB, articleID uint, fields []string) (articleVersion *model.ArticleVersion, err error) {
	err = db.Select(fields).Where(&model.ArticleVersion{ArticleID: articleID}).Last(&articleVersion).Error
	return
}

// GetByArticleIDAndVersion 通过文章ID和版本号获取文章版本
//
//	@receiver dao *ArticleVersionDAO
//	@param db *gorm.DB
//	@param articleID uint
//	@param version uint
//	@param fields []string
//	@return articleVersion *model.ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-18 03:17:06
func (dao *ArticleVersionDAO) GetByArticleIDAndVersion(db *gorm.DB, articleID, version uint, fields []string) (articleVersion *model.ArticleVersion, err error) {
	err = db.Select(fields).Where(&model.ArticleVersion{ArticleID: articleID, Version: version}).Last(&articleVersion).Error
	return
}

// ListByArticleID 通过文章ID获取文章版本列表
//
//	@receiver dao *ArticleVersionDAO
//	@param db *gorm.DB
//	@param articleID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return articleVersions []*model.ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 08:15:44
func (dao *ArticleVersionDAO) ListByArticleID(db *gorm.DB, articleID uint, fields []string, limit, offset int) (articleVersions *[]model.ArticleVersion, err error) {
	err = db.Select(fields).Where(&model.ArticleVersion{ArticleID: articleID}).Limit(limit).Offset(offset).Find(&articleVersions).Error
	return
}
