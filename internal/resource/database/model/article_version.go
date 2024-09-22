package model

import "gorm.io/gorm"

// ArticleVersion 文章版本
//
//	@author centonhuang
//	@update 2024-09-21 06:47:31
type ArticleVersion struct {
	gorm.Model
	ArticleID uint   `json:"article_id" gorm:"column:article_id;uniqueIndex:idx_article_version;comment:文章ID"`
	Version   int    `json:"version" gorm:"column:version_number;uniqueIndex:idx_article_version;comment:版本号"`
	Content   string `json:"content" gorm:"column:content;not null;comment:文章内容"`
}
