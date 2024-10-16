package model

import (
	"errors"

	"github.com/hcd233/Aris-blog/internal/resource/database"
	"gorm.io/gorm"
)

// ArticleVersion 文章版本
//
//	@author centonhuang
//	@update 2024-09-21 06:47:31
type ArticleVersion struct {
	gorm.Model
	ArticleID uint     `json:"article_id" gorm:"column:article_id;uniqueIndex:idx_article_version;comment:文章ID"`
	Article   *Article `json:"article" gorm:"foreignKey:ArticleID"`
	Version   uint     `json:"version" gorm:"column:version;uniqueIndex:idx_article_version;comment:版本号"`
	Content   string   `json:"content" gorm:"column:content;not null;comment:文章内容"`
}

// Create 创建文章版本
//
//	@receiver av *ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-16 01:54:18
func (av *ArticleVersion) Create() (err error) {
	latestVersion, err := QueryLatestArticleVersionByArticleID(av.ArticleID, []string{"version_number"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	av.Version = latestVersion.Version + 1
	err = database.DB.Create(av).Error
	return
}

// GetBasicInfo 获取文章基本信息
//
//	@receiver a *Article
//	@return map
//	@author centonhuang
//	@update 2024-09-21 09:35:50
func (av *ArticleVersion) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        av.ID,
		"create_at": av.CreatedAt,
		"version":   av.Version,
	}
}

// GetDetailedInfo 获取文章详细信息
//
//	@receiver a *Article
//	@return map
//	@author centonhuang
//	@update 2024-09-21 09:21:50
func (av *ArticleVersion) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        av.ID,
		"create_at": av.CreatedAt,
		"version":   av.Version,
		"content":   av.Content,
	}
}

// QueryLatestArticleVersionByArticleID 查询指定 ArticleID 的最新文章版本
//
//	@param articleID uint
//	@param fields []string
//	@return latestVersion ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-16 10:08:53
func QueryLatestArticleVersionByArticleID(articleID uint, fields []string) (latestVersion *ArticleVersion, err error) {
	// err = database.DB.Where(&ArticleVersion{ArticleID: articleID}).Order("version_number DESC").First(&latestVersion).Error
	err = database.DB.Select(fields).Where(&ArticleVersion{ArticleID: articleID}).Last(&latestVersion).Error
	return
}

// QueryArticleVersionsByArticleID 查询指定 ArticleID 的所有文章版本
//
//	@param articleID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return versions []ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-16 10:08:48
func QueryArticleVersionsByArticleID(articleID uint, fields []string, limit, offset int) (versions *[]ArticleVersion, err error) {
	err = database.DB.Select(fields).Where(&ArticleVersion{ArticleID: articleID}).Limit(limit).Offset(offset).Find(&versions).Error
	return
}
