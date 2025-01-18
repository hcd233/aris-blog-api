package protocol

import "github.com/hcd233/aris-blog-api/internal/resource/database/model"

// RefreshTokenBody 刷新token请求体
//
//	Author centonhuang
//	Update 2024-11-09 02:56:39
type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// UpdateUserBody 更新用户请求体
//
//	author centonhuang
//	update 2024-09-18 02:39:31
type UpdateUserBody struct {
	UserName string `json:"userName" binding:"required"`
}

// CreateArticleBody 创建文章请求体
//
//	author centonhuang
//	update 2024-09-21 09:59:55
type CreateArticleBody struct {
	Title      string   `json:"title" binding:"required"`
	Slug       string   `json:"slug" binding:"required"`
	Tags       []string `json:"tags" binding:"required"`
	CategoryID uint     `json:"categoryID" binding:"required"`
}

// UpdateArticleBody 更新文章请求体
//
//	author centonhuang
//	update 2024-09-22 03:56:09
type UpdateArticleBody struct {
	Title      string `json:"title" binding:"omitempty"`
	Slug       string `json:"slug" binding:"omitempty"`
	CategoryID uint   `json:"categoryID" binding:"omitempty"`
}

// UpdateTagBody 更新标签请求体
//
//	author centonhuang
//	update 2024-09-22 03:20:00
type UpdateTagBody struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
}

// CreateTagBody 创建标签请求体
//
//	author centonhuang
//	update 2024-09-22 03:20:00
type CreateTagBody struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
}

// CreateCategoryBody 创建分类请求体
//
//	author centonhuang
//	update 2024-09-28 07:02:11
type CreateCategoryBody struct {
	ParentID uint   `json:"parentID" binding:"omitempty"`
	Name     string `json:"name" binding:"required"`
}

// UpdateCategoryBody 更新分类请求体
//
//	author centonhuang
//	update 2024-10-02 03:46:26
type UpdateCategoryBody struct {
	Name     string `json:"name" binding:"omitempty"`
	ParentID uint   `json:"parentID" binding:"omitempty"`
}

// CreateArticleVersionBody 创建文章版本请求体
//
//	author centonhuang
//	update 2024-10-17 12:43:32
type CreateArticleVersionBody struct {
	Content string `json:"content" binding:"required,min=100,max=20000"`
}

// UpdateArticleStatusBody 更新文章状态请求体
//
//	author centonhuang
//	update 2024-10-17 09:28:07
type UpdateArticleStatusBody struct {
	Status model.ArticleStatus `json:"status" binding:"required,oneof=draft publish"`
}

// CreateArticleCommentBody 创建文章评论请求体
//
//	author centonhuang
//	update 2024-10-24 04:29:01
type CreateArticleCommentBody struct {
	ArticleID uint   `json:"articleID" binding:"required"`
	ReplyTo   uint   `json:"replyTo" binding:"omitempty"`
	Content   string `json:"content" binding:"required,min=1,max=300"`
}

// LikeBody 点赞请求体
//
//	author centonhuang
//	update 2024-10-29 06:49:35
type LikeBody struct {
	Undo bool `json:"undo" binding:"omitempty"`
}

// LikeArticleBody 点赞文章请求体
//
//	author centonhuang
//	update 2024-10-29 06:49:41
type LikeArticleBody struct {
	LikeBody
	ArticleID uint `json:"articleID" binding:"required"`
}

// LikeCommentBody 点赞评论请求体
//
//	author centonhuang
//	update 2024-10-29 06:59:21
type LikeCommentBody struct {
	LikeBody
	CommentID uint `json:"commentID" binding:"required"`
}

// LikeTagBody 点赞标签请求体
//
//	author centonhuang
//	update 2024-10-29 06:50:42
type LikeTagBody struct {
	LikeBody
	TagID uint `json:"tagID" binding:"required"`
}

// LogUserViewArticleBody 记录文章浏览请求体
//
//	author centonhuang
//	update 2024-11-09 07:34:14
type LogUserViewArticleBody struct {
	ArticleID uint `json:"articleID" binding:"required"`
	Progress  int8 `json:"progress" binding:"min=0,max=100"`
}

// Template 提示词模板
//
//	author centonhuang
//	update 2024-12-08 16:39:28
type Template struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreatePromptBody 创建提示词请求体
//
//	author centonhuang
//	update 2024-12-08 16:39:23
type CreatePromptBody struct {
	Templates []Template `json:"templates" binding:"required,min=1,dive"`
}

// AIAPPRequestBody AI应用请求体
//
//	author centonhuang
//	update 2024-12-08 16:39:20
type AIAPPRequestBody struct {
	Temperature float64 `json:"temperature" binding:"omitempty,min=0,max=1"`
}

// GenerateContentCompletionBody 生成内容补全请求体
//
//	author centonhuang
//	update 2024-12-08 16:38:59
type GenerateContentCompletionBody struct {
	AIAPPRequestBody
	Context     string `json:"context" binding:"required"`
	Instruction string `json:"instruction" binding:"required"`
	Reference   string `json:"reference" binding:"omitempty"`
}

// GenerateArticleSummaryBody 生成文章摘要请求体
//
//	author centonhuang
//	update 2024-12-08 16:39:09
type GenerateArticleSummaryBody struct {
	AIAPPRequestBody
	ArticleID   uint   `json:"articleID" binding:"required"`
	Instruction string `json:"instruction" binding:"required"`
}

// GenerateArticleQABody 生成文章问答请求体
//
//	author centonhuang
//	update 2024-12-08 16:38:59
type GenerateArticleQABody struct {
	AIAPPRequestBody
	ArticleID uint   `json:"articleID" binding:"required"`
	Question  string `json:"question" binding:"required"`
}

// GenerateTermExplainationBody 生成术语解释请求体
//
//	author centonhuang
//	update 2024-12-08 16:38:59
type GenerateTermExplainationBody struct {
	AIAPPRequestBody
	ArticleID uint   `json:"articleID" binding:"required"`
	Term      string `json:"term" binding:"required"`
	Position  uint   `json:"position" binding:"required"`
}
