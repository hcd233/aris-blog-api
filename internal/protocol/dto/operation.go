package dto

// LikeArticleRequestBody 点赞文章请求体
type LikeArticleRequestBody struct {
    ArticleID uint `json:"articleID" doc:"文章 ID"`
    Undo      bool `json:"undo" doc:"是否撤销点赞"`
}

// LikeArticleRequest 点赞文章请求
type LikeArticleRequest struct {
    Body *LikeArticleRequestBody `json:"body" doc:"点赞文章的请求体"`
}

// LikeCommentRequestBody 点赞评论请求体
type LikeCommentRequestBody struct {
    CommentID uint `json:"commentID" doc:"评论 ID"`
    Undo      bool `json:"undo" doc:"是否撤销点赞"`
}

// LikeCommentRequest 点赞评论请求
type LikeCommentRequest struct {
    Body *LikeCommentRequestBody `json:"body" doc:"点赞评论的请求体"`
}

// LikeTagRequestBody 点赞标签请求体
type LikeTagRequestBody struct {
    TagID uint `json:"tagID" doc:"标签 ID"`
    Undo  bool `json:"undo" doc:"是否撤销点赞"`
}

// LikeTagRequest 点赞标签请求
type LikeTagRequest struct {
    Body *LikeTagRequestBody `json:"body" doc:"点赞标签的请求体"`
}

// LogArticleViewRequestBody 记录文章浏览请求体
type LogArticleViewRequestBody struct {
    ArticleID uint `json:"articleID" doc:"文章 ID"`
    Progress  int8 `json:"progress" minimum:"0" maximum:"100" doc:"阅读进度，范围 0-100"`
}

// LogArticleViewRequest 记录文章浏览请求
type LogArticleViewRequest struct {
    Body *LogArticleViewRequestBody `json:"body" doc:"记录文章浏览的请求体"`
}

