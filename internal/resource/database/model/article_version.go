package model

import (
	"gorm.io/gorm"
)

// ArticleVersion 文章版本
//
//	author centonhuang
//	update 2024-09-21 06:47:31
type ArticleVersion struct {
	gorm.Model
	ArticleID   uint     `json:"article_id" gorm:"column:article_id;uniqueIndex:idx_article_version;comment:文章ID"`
	Article     *Article `json:"article" gorm:"foreignKey:ArticleID"`
	Version     uint     `json:"version" gorm:"column:version;uniqueIndex:idx_article_version;comment:版本号"`
	Content     string   `json:"content" gorm:"column:content;type:LONGTEXT null;comment:文章内容"`
	Description string   `json:"description" gorm:"column:description;comment:版本描述"`
	Summary     string   `json:"summary" gorm:"column:summary;type:LONGTEXT null;comment:版本摘要"`
}
