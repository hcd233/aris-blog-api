package protocol

import "github.com/hcd233/Aris-blog/internal/resource/database/model"

// UpdateUserBody 更新用户请求体
//
//	@author centonhuang
//	@update 2024-09-18 02:39:31
type UpdateUserBody struct {
	UserName string `json:"userName" binding:"required"`
}

// CreateArticleBody 创建文章请求体
//
//	@author centonhuang
//	@update 2024-09-21 09:59:55
type CreateArticleBody struct {
	Title      string   `json:"title" binding:"required"`
	Slug       string   `json:"slug"`
	Tags       []string `json:"tags"`
	CategoryID uint     `json:"categoryID" binding:"omitempty"`
}

// UpdateArticleBody 更新文章请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:56:09
type UpdateArticleBody struct {
	Title      string              `json:"title" binding:"omitempty"`
	Slug       string              `json:"slug" binding:"omitempty"`
	CategoryID uint                `json:"categoryID" binding:"omitempty"`
	Status     model.ArticleStatus `json:"status" binding:"omitempty,oneof=draft publish"`
}

// UpdateTagBody 更新标签请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:20:00
type UpdateTagBody struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
}

// CreateTagBody 创建标签请求体
//
//	@author centonhuang
//	@update 2024-09-22 03:20:00
type CreateTagBody struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
}

// CreateCategoryBody 创建分类请求体
//
//	@author centonhuang
//	@update 2024-09-28 07:02:11
type CreateCategoryBody struct {
	ParentID uint   `json:"parentID" binding:"omitempty"`
	Name     string `json:"name" binding:"required"`
}

// UpdateCategoryBody 更新分类请求体
//
//	@author centonhuang
//	@update 2024-10-02 03:46:26
type UpdateCategoryBody struct {
	Name     string `json:"name" binding:"omitempty"`
	ParentID uint   `json:"parentID" binding:"omitempty"`
}

// CreateArticleVersionBody 创建文章版本请求体
//
//	@author centonhuang
//	@update 2024-10-17 12:43:32
type CreateArticleVersionBody struct {
	Content string `json:"content" binding:"required"`
}

// UpdateArticleStatusBody 更新文章状态请求体
//
//	@author centonhuang
//	@update 2024-10-17 09:28:07
type UpdateArticleStatusBody struct {
	Status model.ArticleStatus `json:"status" binding:"required,oneof=draft publish"`
}
