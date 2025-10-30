package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CommentHandler 评论处理器
//
//	author centonhuang
//	update 2025-10-31 05:10:00
type CommentHandler interface {
	HandleCreateArticleComment(ctx context.Context, req *CommentCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateArticleCommentResponse], error)
	HandleDeleteComment(ctx context.Context, req *CommentDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteCommentResponse], error)
	HandleListArticleComments(ctx context.Context, req *CommentListArticleRequest) (*protocol.HumaHTTPResponse[*protocol.ListArticleCommentsResponse], error)
	HandleListChildrenComments(ctx context.Context, req *CommentListChildrenRequest) (*protocol.HumaHTTPResponse[*protocol.ListChildrenCommentsResponse], error)
}

type commentHandler struct {
	svc service.CommentService
}

// NewCommentHandler 创建评论处理器
//
//	return CommentHandler
//	author centonhuang
//	update 2025-10-31 05:10:00
func NewCommentHandler() CommentHandler {
	return &commentHandler{
		svc: service.NewCommentService(),
	}
}

// CommentPathParam 评论路径参数
type CommentPathParam struct {
	CommentID uint `path:"commentID" doc:"评论 ID"`
}

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	Body *protocol.CreateArticleCommentBody `json:"body" doc:"创建评论请求体"`
}

// CommentDeleteRequest 删除评论请求
type CommentDeleteRequest struct {
	CommentPathParam
}

// CommentListArticleRequest 列出文章评论请求
type CommentListArticleRequest struct {
	ArticlePathParam
	PaginationQuery
}

// CommentListChildrenRequest 列出子评论请求
type CommentListChildrenRequest struct {
	CommentPathParam
	PaginationQuery
}

// HandleCreateArticleComment 创建文章评论
//
//	receiver h *commentHandler
//	param ctx context.Context
//	param req *CommentCreateRequest
//	return *protocol.HumaHTTPResponse[*protocol.CreateArticleCommentResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:10:00
func (h *commentHandler) HandleCreateArticleComment(ctx context.Context, req *CommentCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateArticleCommentResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.CreateArticleCommentRequest{
		UserID:    userID,
		ArticleID: req.Body.ArticleID,
		Content:   req.Body.Content,
		ReplyTo:   req.Body.ReplyTo,
	}

	return util.WrapHTTPResponse(h.svc.CreateArticleComment(ctx, serviceReq))
}

// HandleDeleteComment 删除评论
//
//	receiver h *commentHandler
//	param ctx context.Context
//	param req *CommentDeleteRequest
//	return *protocol.HumaHTTPResponse[*protocol.DeleteCommentResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:10:00
func (h *commentHandler) HandleDeleteComment(ctx context.Context, req *CommentDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteCommentResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.DeleteCommentRequest{
		UserID:    userID,
		CommentID: req.CommentID,
	}

	return util.WrapHTTPResponse(h.svc.DeleteComment(ctx, serviceReq))
}

// HandleListArticleComments 列出文章评论
//
//	receiver h *commentHandler
//	param ctx context.Context
//	param req *CommentListArticleRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListArticleCommentsResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:10:00
func (h *commentHandler) HandleListArticleComments(ctx context.Context, req *CommentListArticleRequest) (*protocol.HumaHTTPResponse[*protocol.ListArticleCommentsResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.ListArticleCommentsRequest{
		UserID:        userID,
		ArticleID:     req.ArticleID,
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListArticleComments(ctx, serviceReq))
}

// HandleListChildrenComments 列出子评论
//
//	receiver h *commentHandler
//	param ctx context.Context
//	param req *CommentListChildrenRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListChildrenCommentsResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:10:00
func (h *commentHandler) HandleListChildrenComments(ctx context.Context, req *CommentListChildrenRequest) (*protocol.HumaHTTPResponse[*protocol.ListChildrenCommentsResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.ListChildrenCommentsRequest{
		UserID:        userID,
		CommentID:     req.CommentID,
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListChildrenComments(ctx, serviceReq))
}
