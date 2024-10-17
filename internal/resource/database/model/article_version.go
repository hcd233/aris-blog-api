package model

import (
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"gorm.io/gorm"
)

// ArticleVersion 文章版本
//
//	@author centonhuang
//	@update 2024-09-21 06:47:31
type ArticleVersion struct {
	gorm.Model
	ArticleID   uint     `json:"article_id" gorm:"column:article_id;uniqueIndex:idx_article_version;comment:文章ID"`
	Article     *Article `json:"article" gorm:"foreignKey:ArticleID"`
	Version     uint     `json:"version" gorm:"column:version;uniqueIndex:idx_article_version;comment:版本号"`
	Content     string   `json:"content" gorm:"column:content;not null;comment:文章内容"`
	Description string   `json:"description" gorm:"column:description;comment:版本描述"`
}

// Create 创建文章版本
//
//	@receiver av *ArticleVersion
//	@return err error
//	@author centonhuang
//	@update 2024-10-16 01:54:18
func (av *ArticleVersion) Create() (err error) {
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
