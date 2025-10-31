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

// CreateCommentRequestBody 创建评论请求体
type CreateCommentRequestBody struct {
	ArticleID uint   `json:"articleID" doc:"Article ID"`
	ReplyTo   uint   `json:"replyTo" doc:"Parent comment ID if this is a reply"`
	Content   string `json:"content" doc:"Comment content"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Body *CreateCommentRequestBody `json:"body" doc:"Fields for creating comment"`
}

// CreateCommentResponse 创建评论响应
type CreateCommentResponse struct {
	Comment *Comment `json:"comment" doc:"Comment details"`
}

// DeleteCommentRequest 删除评论请求
type DeleteCommentRequest struct {
	CommentPathParam
}

// ListArticleCommentRequest 列出文章评论请求
type ListArticleCommentRequest struct {
	ArticlePathParam
	CommonParam
}

// ListArticleCommentResponse 列出文章评论响应
type ListArticleCommentResponse struct {
	Comments []*Comment `json:"comments" doc:"List of comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ListChildrenCommentRequest 列出子评论请求
type ListChildrenCommentRequest struct {
	CommentPathParam
	CommonParam
}

// ListChildrenCommentResponse 列出子评论响应
type ListChildrenCommentResponse struct {
	Comments []*Comment `json:"comments" doc:"List of child comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}
