package model

import (
	"time"

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
	ID          uint             `json:"id" gorm:"column:id;primary_key;auto_increment;comment:文章ID"`
	Title       string           `json:"title" gorm:"column:title;not null;comment:文章标题"`
	Slug        string           `json:"slug" gorm:"column:slug;not null;uniqueIndex:idx_user_slug;comment:文章slug"`
	UserID      uint             `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_slug;comment:用户ID"`
	User        *User            `json:"user" gorm:"foreignKey:UserID"`
	CategoryID  uint             `json:"category_id" gorm:"column:category_id;null;comment:类别ID"`
	Category    *Category        `json:"category" gorm:"foreignKey:CategoryID"`
	Status      ArticleStatus    `json:"status" gorm:"column:status;not null;default:'draft';comment:文章状态"`
	PublishedAt time.Time        `json:"published_at" gorm:"column:published_at;default:NULL;comment:发布时间"`
	Views       uint             `json:"views" gorm:"column:views;default:0;comment:浏览数"`
	Likes       uint             `json:"likes" gorm:"column:likes;default:0;comment:点赞数"`
	Tags        []Tag            `json:"tags" gorm:"many2many:article_tags;"`
	Comments    []Comment        `json:"comments" gorm:"foreignKey:ArticleID"`
	Versions    []ArticleVersion `json:"versions" gorm:"foreignKey:ArticleID"`
}

// GetBasicInfo 获取文章基本信息
//
//	@receiver a *Article
//	@return map
//	@author centonhuang
//	@update 2024-09-21 09:35:50
func (a *Article) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":     a.ID,
		"title":  a.Title,
		"slug":   a.Slug,
		"status": a.Status,
	}
}

// GetLikeInfo 获取文章点赞信息
//
//	@receiver a *Article
//	@return map
//	@author centonhuang
//	@update 2024-11-03 07:34:08
func (a *Article) GetLikeInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":           a.ID,
		"title":        a.Title,
		"slug":         a.Slug,
		"published_at": a.PublishedAt,
		"user":         a.User.GetBasicInfo(),
		"tags":         lo.Map(a.Tags, func(tag Tag, idx int) map[string]interface{} { return tag.GetBasicInfo() }),
		"likes":        a.Likes,
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
		"category":     a.CategoryID,
		"status":       a.Status,
		"published_at": a.PublishedAt,
		"tags":         lo.Map(a.Tags, func(tag Tag, idx int) map[string]interface{} { return tag.GetBasicInfo() }),
		"comments":     a.Comments,
		"views":        a.Views,
		"likes":        a.Likes,
		"versions":     lo.Map(a.Versions, func(version ArticleVersion, idx int) map[string]interface{} { return version.GetBasicInfo() }),
	}
}
