package model

import (
	"time"

	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// ArticleStatus 文章状态
//
//	@author centonhuang
//	@update 2024-09-21 06:53:10
type ArticleStatus string

const (

	// ArticleStatusDraft ArticleStatus 文稿状态
	//	@update 2024-09-21 06:52:56
	ArticleStatusDraft ArticleStatus = "draft"

	// ArticleStatusPublish ArticleStatus 发布状态
	//	@update 2024-09-21 06:53:04
	ArticleStatusPublish ArticleStatus = "publish"
)

// Article 文章数据库模型
//
//	@author centonhuang
//	@update 2024-09-21 06:46:05
type Article struct {
	gorm.Model
	ID          uint             `json:"id" gorm:"column:id;primary_key;auto_increment"`
	Title       string           `json:"title" gorm:"column:title"`
	Slug        string           `json:"slug" gorm:"column:slug;not null;uniqueIndex:idx_user_slug"`
	UserID      uint             `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_slug"`
	Category    string           `json:"category" gorm:"column:category;not null"`
	Status      ArticleStatus    `json:"status" gorm:"column:status;not null;default:'draft'"`
	PublishedAt time.Time        `json:"published_at" gorm:"column:published_at;default:NULL"`
	Tags        []Tag            `json:"tags" gorm:"many2many:article_tags;"`
	Comments    []Comment        `json:"comments" gorm:"foreignKey:ArticleID"`
	Views       uint             `json:"views" gorm:"column:views;default:0"`
	Likes       uint             `json:"likes" gorm:"column:likes;default:0"`
	Versions    []ArticleVersion `json:"versions" gorm:"foreignKey:ArticleID"`
}

// Create 创建文章
//
//	@receiver a *Article
//	@return err error
//	@author centonhuang
//	@update 2024-09-21 09:21:16
func (a *Article) Create() (err error) {
	err = database.DB.Create(a).Error
	return
}

// GetBasicInfo 获取文章基本信息
//
//	@receiver a *Article
//	@return map
//	@author centonhuang
//	@update 2024-09-21 09:35:50
func (a *Article) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":    a.ID,
		"title": a.Title,
		"slug":  a.Slug,
		// TODO support tags and category
	}
}

// GetDetailedInfo 获取文章详细信息
//
//	@receiver a *Article
//	@return map
//	@author centonhuang
//	@update 2024-09-21 09:21:50
func (a *Article) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":           a.ID,
		"title":        a.Title,
		"slug":         a.Slug,
		"user_id":      a.UserID,
		"category":     a.Category,
		"status":       a.Status,
		"published_at": a.PublishedAt,
		"tags":         a.Tags,
		"comments":     a.Comments,
		"views":        a.Views,
		"likes":        a.Likes,
		"versions":     a.Versions,
	}
}

// QueryArticleBySlugAndUserName 根据文章名和用户名查询文章
//
//	@param articleSlug string
//	@param fields []string
//	@return article *Article
//	@return err error
//	@author centonhuang
//	@update 2024-09-21 09:17:50
func QueryArticleBySlugAndUserName(articleSlug string, userName string, fields []string) (article *Article, err error) {
	err = database.DB.Select(
		lo.Map(fields, func(field string, idx int) string { return "articles." + field }),
	).Joins(
		"JOIN users ON user_id = users.id",
	).Where(
		"slug = ? AND users.name = ?", articleSlug, userName,
	).First(&article).Error
	return
}

// QueryArticlesByUserName 根据用户名查询文章
//
//	@param userName string
//	@param limit int
//	@param offset int
//	@return articles []Article
//	@return err error
//	@author centonhuang
//	@update 2024-09-21 09:07:55
func QueryArticlesByUserName(userID uint, limit int, offset int) (articles []Article, err error) {
	err = database.DB.Preload("Tags").Preload("Comments").Preload("Versions").Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&articles).Error
	return
}
