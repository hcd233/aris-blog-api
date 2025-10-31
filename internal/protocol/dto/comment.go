// Package dto 评论DTO
package dto

// Comment 评论
//
//	author centonhuang
//	update 2025-10-30
type Comment struct {
	CommentID uint   `json:"commentID" doc:"Unique identifier for the comment"`
	Content   string `json:"content" doc:"Comment content"`
	UserID    uint   `json:"userID" doc:"ID of the user who created the comment"`
	ReplyTo   uint   `json:"replyTo" doc:"ID of the comment this is replying to (0 for top-level comments)"`
	CreatedAt string `json:"createdAt" doc:"Creation timestamp"`
	Likes     uint   `json:"likes" doc:"Number of likes"`
}

// CreateArticleCommentRequest 创建文章评论请求
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleCommentRequest struct {
	Body *CreateArticleCommentBody `json:"body" doc:"Request body containing comment information"`
}

// CreateArticleCommentBody 创建文章评论请求体
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleCommentBody struct {
	ArticleID uint   `json:"articleID" doc:"ID of the article to comment on"`
	ReplyTo   uint   `json:"replyTo,omitempty" doc:"ID of the comment to reply to (optional)"`
	Content   string `json:"content" minLength:"1" maxLength:"300" doc:"Comment content (1-300 characters)"`
}

// CreateArticleCommentResponse 创建文章评论响应
//
//	author centonhuang
//	update 2025-10-30
type CreateArticleCommentResponse struct {
	Comment *Comment `json:"comment" doc:"Created comment"`
}

// DeleteCommentRequest 删除评论请求
//
//	author centonhuang
//	update 2025-10-30
type DeleteCommentRequest struct {
	CommentID uint `path:"commentID" doc:"Unique identifier of the comment to delete"`
}

// DeleteCommentResponse 删除评论响应
//
//	author centonhuang
//	update 2025-10-30
type DeleteCommentResponse struct{}

// ListArticleCommentsRequest 列出文章评论请求
//
//	author centonhuang
//	update 2025-10-30
type ListArticleCommentsRequest struct {
	ArticleID uint `path:"articleID" doc:"Article ID"`
	Page      *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize  *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListArticleCommentsResponse 列出文章评论响应
//
//	author centonhuang
//	update 2025-10-30
type ListArticleCommentsResponse struct {
	Comments []*Comment `json:"comments" doc:"List of top-level comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}

// ListChildrenCommentsRequest 列出子评论请求
//
//	author centonhuang
//	update 2025-10-30
type ListChildrenCommentsRequest struct {
	CommentID uint `path:"commentID" doc:"Parent comment ID"`
	Page      *int `query:"page" minimum:"1" doc:"Page number (starting from 1)"`
	PageSize  *int `query:"pageSize" minimum:"1" maximum:"50" doc:"Number of items per page"`
}

// ListChildrenCommentsResponse 列出子评论响应
//
//	author centonhuang
//	update 2025-10-30
type ListChildrenCommentsResponse struct {
	Comments []*Comment `json:"comments" doc:"List of child comments"`
	PageInfo *PageInfo  `json:"pageInfo" doc:"Pagination information"`
}
