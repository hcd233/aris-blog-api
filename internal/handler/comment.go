package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CommentHandler 评论处理器
//
//	author centonhuang
//	update 2025-10-30
type CommentHandler interface {
	HandleCreateArticleComment(ctx context.Context, req *dto.CreateArticleCommentRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleCommentResponse], error)
	HandleDeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) (*protocol.HumaHTTPResponse[*dto.DeleteCommentResponse], error)
	HandleListArticleComments(ctx context.Context, req *dto.ListArticleCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleCommentsResponse], error)
	HandleListChildrenComments(ctx context.Context, req *dto.ListChildrenCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListChildrenCommentsResponse], error)
}

type commentHandler struct {
	svc service.CommentService
}

// NewCommentHandler 创建评论处理器
//
//	return CommentHandler
//	author centonhuang
//	update 2025-10-30
func NewCommentHandler() CommentHandler {
	return &commentHandler{
		svc: service.NewCommentService(),
	}
}

func (h *commentHandler) HandleCreateArticleComment(ctx context.Context, req *dto.CreateArticleCommentRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleCommentResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.CreateArticleCommentRequest{
		UserID:    userID,
		ArticleID: req.Body.ArticleID,
		Content:   req.Body.Content,
		ReplyTo:   req.Body.ReplyTo,
	}

	svcRsp, err := h.svc.CreateArticleComment(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.CreateArticleCommentResponse](nil, err)
	}

	rsp := &dto.CreateArticleCommentResponse{
		Comment: &dto.Comment{
			CommentID: svcRsp.Comment.CommentID,
			Content:   svcRsp.Comment.Content,
			UserID:    svcRsp.Comment.UserID,
			ReplyTo:   svcRsp.Comment.ReplyTo,
			CreatedAt: svcRsp.Comment.CreatedAt,
			Likes:     svcRsp.Comment.Likes,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *commentHandler) HandleDeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) (*protocol.HumaHTTPResponse[*dto.DeleteCommentResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.DeleteCommentRequest{
		UserID:    userID,
		CommentID: req.CommentID,
	}

	_, err := h.svc.DeleteComment(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.DeleteCommentResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.DeleteCommentResponse{}, nil)
}

func (h *commentHandler) HandleListArticleComments(ctx context.Context, req *dto.ListArticleCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleCommentsResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListArticleCommentsRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListArticleComments(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListArticleCommentsResponse](nil, err)
	}

	comments := make([]*dto.Comment, len(svcRsp.Comments))
	for i, comment := range svcRsp.Comments {
		comments[i] = &dto.Comment{
			CommentID: comment.CommentID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			ReplyTo:   comment.ReplyTo,
			CreatedAt: comment.CreatedAt,
			Likes:     comment.Likes,
		}
	}

	rsp := &dto.ListArticleCommentsResponse{
		Comments: comments,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *commentHandler) HandleListChildrenComments(ctx context.Context, req *dto.ListChildrenCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListChildrenCommentsResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListChildrenCommentsRequest{
		UserID:    userID,
		CommentID: req.CommentID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListChildrenComments(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListChildrenCommentsResponse](nil, err)
	}

	comments := make([]*dto.Comment, len(svcRsp.Comments))
	for i, comment := range svcRsp.Comments {
		comments[i] = &dto.Comment{
			CommentID: comment.CommentID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			ReplyTo:   comment.ReplyTo,
			CreatedAt: comment.CreatedAt,
			Likes:     comment.Likes,
		}
	}

	rsp := &dto.ListChildrenCommentsResponse{
		Comments: comments,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}
