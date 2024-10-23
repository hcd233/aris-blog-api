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
func (dao *ArticleDAO) GetBySlugAndUserID(db *gorm.DB, slug string, userID uint, fields []string) (article *model.Article, err error) {
	err = db.Select(fields).Where(&model.Article{Slug: slug, UserID: userID}).First(&article).Error
	return
}

// GetAllBySlugAndUserID  通过slug和用户ID获取文章全部字段
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param slug string
//	@param userID uint
//	@return article *model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-10-18 02:21:20
func (dao *ArticleDAO) GetAllBySlugAndUserID(db *gorm.DB, slug string, userID uint) (article *model.Article, err error) {
	err = db.Preload("Category").Preload("Tags").Preload("User").Where(&model.Article{Slug: slug, UserID: userID}).First(&article).Error
	return
}

// ListByUserID 通过用户ID获取文章列表
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param userID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return articles *[]model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 07:34:43
func (dao *ArticleDAO) ListByUserID(db *gorm.DB, userID uint, fields []string, limit, offset int) (articles *[]model.Article, err error) {
	err = db.Select(fields).Where(&model.Article{UserID: userID}).Limit(limit).Offset(offset).Find(&articles).Error
	return
}

// ListByCategoryID 通过类别ID获取文章列表
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param categoryID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return articles *[]model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 07:34:57
func (dao *ArticleDAO) ListByCategoryID(db *gorm.DB, categoryID uint, fields []string, limit, offset int) (articles *[]model.Article, err error) {
	err = db.Select(fields).Where(&model.Article{CategoryID: categoryID}).Limit(limit).Offset(offset).Find(&articles).Error
	return
}

// ListByPublished 列出已发布的文章
//
//	@receiver dao *ArticleDAO
//	@param db *gorm.DB
//	@param userID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return articles *[]model.Article
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 07:34:43
func (dao *ArticleDAO) ListByPublished(db *gorm.DB, fields []string, limit, offset int) (articles *[]model.Article, err error) {
	err = db.Select(fields).Where(&model.Article{Status: model.ArticleStatusPublish}).Limit(limit).Offset(offset).Find(&articles).Error
	return
}
