package model

import (
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
	Content     string   `json:"content" gorm:"column:content;type:LONGTEXT null;comment:文章内容"`
	Description string   `json:"description" gorm:"column:description;comment:版本描述"`
	Summary     string   `json:"summary" gorm:"column:summary;type:LONGTEXT null;comment:版本摘要"`
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
		"createdAt": av.CreatedAt,
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
		"createdAt": av.CreatedAt,
		"version":   av.Version,
		"content":   av.Content,
		"summary":   av.Summary,
	}
}
