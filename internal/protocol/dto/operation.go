// Package dto 操作DTO
package dto

// LikeArticleRequest 点赞文章请求
//
//	author centonhuang
//	update 2025-10-30
type LikeArticleRequest struct {
	Body *LikeArticleBody `json:"body" doc:"Request body containing like information"`
}

// LikeArticleBody 点赞文章请求体
//
//	author centonhuang
//	update 2025-10-30
type LikeArticleBody struct {
	ArticleID uint `json:"articleID" doc:"ID of the article to like"`
	Undo      bool `json:"undo,omitempty" doc:"Set to true to unlike the article"`
}

// LikeArticleResponse 点赞文章响应
//
//	author centonhuang
//	update 2025-10-30
type LikeArticleResponse struct{}

// LikeCommentRequest 点赞评论请求
//
//	author centonhuang
//	update 2025-10-30
type LikeCommentRequest struct {
	Body *LikeCommentBody `json:"body" doc:"Request body containing like information"`
}

// LikeCommentBody 点赞评论请求体
//
//	author centonhuang
//	update 2025-10-30
type LikeCommentBody struct {
	CommentID uint `json:"commentID" doc:"ID of the comment to like"`
	Undo      bool `json:"undo,omitempty" doc:"Set to true to unlike the comment"`
}

// LikeCommentResponse 点赞评论响应
//
//	author centonhuang
//	update 2025-10-30
type LikeCommentResponse struct{}

// LikeTagRequest 点赞标签请求
//
//	author centonhuang
//	update 2025-10-30
type LikeTagRequest struct {
	Body *LikeTagBody `json:"body" doc:"Request body containing like information"`
}

// LikeTagBody 点赞标签请求体
//
//	author centonhuang
//	update 2025-10-30
type LikeTagBody struct {
	TagID uint `json:"tagID" doc:"ID of the tag to like"`
	Undo  bool `json:"undo,omitempty" doc:"Set to true to unlike the tag"`
}

// LikeTagResponse 点赞标签响应
//
//	author centonhuang
//	update 2025-10-30
type LikeTagResponse struct{}

// LogArticleViewRequest 记录文章浏览请求
//
//	author centonhuang
//	update 2025-10-30
type LogArticleViewRequest struct {
	Body *LogArticleViewBody `json:"body" doc:"Request body containing view information"`
}

// LogArticleViewBody 记录文章浏览请求体
//
//	author centonhuang
//	update 2025-10-30
type LogArticleViewBody struct {
	ArticleID uint `json:"articleID" doc:"ID of the article being viewed"`
	Progress  int8 `json:"progress" minimum:"0" maximum:"100" doc:"Reading progress percentage (0-100)"`
}

// LogArticleViewResponse 记录文章浏览响应
//
//	author centonhuang
//	update 2025-10-30
type LogArticleViewResponse struct{}
