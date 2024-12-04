package protocol

import "github.com/hcd233/Aris-blog/internal/resource/database/model"

// CreateTagBody 刷新token请求体
//
//	@Author centonhuang
//	@Update 2024-11-09 02:56:39
type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

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
	Title      string `json:"title" binding:"omitempty"`
	Slug       string `json:"slug" binding:"omitempty"`
	CategoryID uint   `json:"categoryID" binding:"omitempty"`
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
	Content string `json:"content" binding:"required,min=100,max=20000"`
}

// UpdateArticleStatusBody 更新文章状态请求体
//
//	@author centonhuang
//	@update 2024-10-17 09:28:07
type UpdateArticleStatusBody struct {
	Status model.ArticleStatus `json:"status" binding:"required,oneof=draft publish"`
}

// CreateArticleCommentBody 创建文章评论请求体
//
//	@author centonhuang
//	@update 2024-10-24 04:29:01
type CreateArticleCommentBody struct {
	ReplyTo uint   `json:"replyTo" binding:"omitempty"`
	Content string `json:"content" binding:"required,min=1,max=300"`
}

// LikeBody 点赞请求体
//
//	@author centonhuang
//	@update 2024-10-29 06:49:35
type LikeBody struct {
	Undo bool `json:"undo" binding:"omitempty"`
}

// LikeArticleBody 点赞文章请求体
//
//	@author centonhuang
//	@update 2024-10-29 06:49:41
type LikeArticleBody struct {
	LikeBody
	Author      string `json:"author" binding:"required"`
	ArticleSlug string `json:"articleSlug" binding:"required"`
}

// LikeCommentBody 点赞评论请求体
//
//	@author centonhuang
//	@update 2024-10-29 06:59:21
type LikeCommentBody struct {
	LikeBody
	CommentID uint `json:"commentID" binding:"required"`
}

// LikeTagBody 点赞标签请求体
//
//	@author centonhuang
//	@update 2024-10-29 06:50:42
type LikeTagBody struct {
	LikeBody
	TagSlug string `json:"tagSlug" binding:"required"`
}

// LogUserViewArticleBody 记录文章浏览请求体
//
//	@author centonhuang
//	@update 2024-11-09 07:34:14
type LogUserViewArticleBody struct {
	Author      string `json:"author" binding:"required"`
	ArticleSlug string `json:"articleSlug" binding:"required"`
	Progress    int8   `json:"progress" binding:"min=0,max=100"`
}

type Template struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreatePromptBody struct {
	Templates []Template `json:"templates" binding:"required,min=1,dive"`
}

type AIAPPRequestBody struct {
	Temperature float64 `json:"temperature" binding:"omitempty,min=0,max=1"`
}

type GenerateContentCompletionBody struct {
	AIAPPRequestBody
	Context     string `json:"context" binding:"required"`
	Instruction string `json:"instruction" binding:"required"`
	Reference   string `json:"reference" binding:"omitempty"`
}

type GenerateArticleSummaryBody struct {
	AIAPPRequestBody
	ArticleSlug string `json:"articleSlug" binding:"required"`
	Instruction string `json:"instruction" binding:"required"`
}
