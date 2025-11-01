package dto

// LikeArticleRequestBody 点赞文章请求体
type LikeArticleRequestBody struct {
	ArticleID uint `json:"articleID" doc:"Article ID to like"`
	Undo      bool `json:"undo" doc:"Whether to undo the like"`
}

// LikeArticleRequest 点赞文章请求
type LikeArticleRequest struct {
	Body *LikeArticleRequestBody `json:"body" doc:"Fields for liking article"`
}

// LikeCommentRequestBody 点赞评论请求体
type LikeCommentRequestBody struct {
	CommentID uint `json:"commentID" doc:"Comment ID to like"`
	Undo      bool `json:"undo" doc:"Whether to undo the like"`
}

// LikeCommentRequest 点赞评论请求
type LikeCommentRequest struct {
	Body *LikeCommentRequestBody `json:"body" doc:"Fields for liking comment"`
}

// LikeTagRequestBody 点赞标签请求体
type LikeTagRequestBody struct {
	TagID uint `json:"tagID" doc:"Tag ID to like"`
	Undo  bool `json:"undo" doc:"Whether to undo the like"`
}

// LikeTagRequest 点赞标签请求
type LikeTagRequest struct {
	Body *LikeTagRequestBody `json:"body" doc:"Fields for liking tag"`
}

// LogArticleViewRequestBody 记录文章浏览请求体
type LogArticleViewRequestBody struct {
	ArticleID uint `json:"articleID" doc:"Article ID being viewed"`
	Progress  int8 `json:"progress" doc:"Reading progress percentage (0-100)" minimum:"0" maximum:"100"`
}

// LogArticleViewRequest 记录文章浏览请求
type LogArticleViewRequest struct {
	Body *LogArticleViewRequestBody `json:"body" doc:"Fields for logging article view"`
}
