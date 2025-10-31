package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CommentHandler 评论处理器
type CommentHandler interface {
	HandleCreateArticleComment(ctx context.Context, req *dto.CommentCreateRequest) (*protocol.HumaHTTPResponse[*dto.CommentCreateResponse], error)
	HandleDeleteComment(ctx context.Context, req *dto.CommentDeleteRequest) (*protocol.HumaHTTPResponse[*dto.CommentDeleteResponse], error)
	HandleListArticleComments(ctx context.Context, req *dto.CommentListArticleRequest) (*protocol.HumaHTTPResponse[*dto.CommentListArticleResponse], error)
	HandleListChildrenComments(ctx context.Context, req *dto.CommentListChildrenRequest) (*protocol.HumaHTTPResponse[*dto.CommentListChildrenResponse], error)
}

type commentHandler struct {
	svc service.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler() CommentHandler {
	return &commentHandler{
		svc: service.NewCommentService(),
	}
}

func (h *commentHandler) HandleCreateArticleComment(ctx context.Context, req *dto.CommentCreateRequest) (*protocol.HumaHTTPResponse[*dto.CommentCreateResponse], error) {
	if req == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.CreateArticleComment(ctx, req))
}

func (h *commentHandler) HandleDeleteComment(ctx context.Context, req *dto.CommentDeleteRequest) (*protocol.HumaHTTPResponse[*dto.CommentDeleteResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.DeleteComment(ctx, req))
}

func (h *commentHandler) HandleListArticleComments(ctx context.Context, req *dto.CommentListArticleRequest) (*protocol.HumaHTTPResponse[*dto.CommentListArticleResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.ListArticleComments(ctx, req))
}

func (h *commentHandler) HandleListChildrenComments(ctx context.Context, req *dto.CommentListChildrenRequest) (*protocol.HumaHTTPResponse[*dto.CommentListChildrenResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.ListChildrenComments(ctx, req))
}
