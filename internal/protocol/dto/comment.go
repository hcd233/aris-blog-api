package dto

// Comment 评论信息
//
//	author centonhuang
//	update 2025-10-31 05:40:00
type Comment struct {
	CommentID uint   `json:"commentID" doc:"Comment ID"`
	Content   string `json:"content" doc:"Comment content"`
	UserID    uint   `json:"userID" doc:"User ID of the comment author"`
	ReplyTo   uint   `json:"replyTo,omitempty" doc:"Parent comment ID if this is a reply"`
	CreatedAt string `json:"createdAt" doc:"Creation timestamp"`
	Likes     uint   `json:"likes" doc:"Number of likes"`
}

// CommentPathParam 评论路径参数
type CommentPathParam struct {
	CommentID uint `path:"commentID" doc:"Comment ID"`
}

// CommentCreateRequestBody 创建评论请求体
type CommentCreateRequestBody struct {
	ArticleID uint   `json:"articleID" doc:"Article ID"`
	ReplyTo   uint   `json:"replyTo" doc:"Parent comment ID if this is a reply"`
	Content   string `json:"content" doc:"Comment content"`
}

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	Body *CommentCreateRequestBody `json:"body" doc:"Fields for creating comment"`
}

// CommentCreateResponse 创建评论响应
type CommentCreateResponse struct {
	Comment *Comment `json:"comment" doc:"Comment details"`
}

// CommentDeleteRequest 删除评论请求
type CommentDeleteRequest struct {
	CommentPathParam
}

// CommentDeleteResponse 删除评论响应
type CommentDeleteResponse struct{}

// CommentListArticleRequest 列出文章评论请求
type CommentListArticleRequest struct {
	ArticlePathParam
	PaginationQuery
}

// CommentListArticleResponse 列出文章评论响应
type CommentListArticleResponse struct {
	Comments []*Comment `json:"comments" doc:"List of comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// CommentListChildrenRequest 列出子评论请求
type CommentListChildrenRequest struct {
	CommentPathParam
	PaginationQuery
}

// CommentListChildrenResponse 列出子评论响应
type CommentListChildrenResponse struct {
	Comments []*Comment `json:"comments" doc:"List of child comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}
