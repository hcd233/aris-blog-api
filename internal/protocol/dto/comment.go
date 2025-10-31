package dto

// Comment 评论信息
//
//	author centonhuang
//	update 2025-10-31 05:40:00
type Comment struct {
	CommentID uint   `json:"commentID" doc:"评论 ID"`
	Content   string `json:"content" doc:"评论内容"`
	UserID    uint   `json:"userID" doc:"评论用户 ID"`
	ReplyTo   uint   `json:"replyTo,omitempty" doc:"回复评论 ID"`
	CreatedAt string `json:"createdAt" doc:"创建时间"`
	Likes     uint   `json:"likes" doc:"点赞数量"`
}

// CommentPathParam 评论路径参数
type CommentPathParam struct {
	CommentID uint `path:"commentID" doc:"评论 ID"`
}

// CommentCreateRequestBody 创建评论请求体
type CommentCreateRequestBody struct {
	ArticleID uint   `json:"articleID" doc:"文章 ID"`
	ReplyTo   uint   `json:"replyTo" doc:"回复评论 ID"`
	Content   string `json:"content" doc:"评论内容"`
}

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	UserID uint                      `json:"-"`
	Body   *CommentCreateRequestBody `json:"body" doc:"创建评论字段"`
}

// CommentCreateResponse 创建评论响应
type CommentCreateResponse struct {
	Comment *Comment `json:"comment" doc:"评论详情"`
}

// CommentDeleteRequest 删除评论请求
type CommentDeleteRequest struct {
	CommentPathParam
	UserID uint `json:"-"`
}

// CommentDeleteResponse 删除评论响应
type CommentDeleteResponse struct{}

// CommentListArticleRequest 列出文章评论请求
type CommentListArticleRequest struct {
	ArticlePathParam
	UserID uint `json:"-"`
	PaginationQuery
}

// CommentListArticleResponse 列出文章评论响应
type CommentListArticleResponse struct {
	Comments []*Comment `json:"comments" doc:"评论列表"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"分页信息"`
}

// CommentListChildrenRequest 列出子评论请求
type CommentListChildrenRequest struct {
	CommentPathParam
	UserID uint `json:"-"`
	PaginationQuery
}

// CommentListChildrenResponse 列出子评论响应
type CommentListChildrenResponse struct {
	Comments []*Comment `json:"comments" doc:"子评论列表"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"分页信息"`
}
