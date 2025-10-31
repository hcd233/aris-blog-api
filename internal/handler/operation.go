package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// OperationHandler 用户操作处理器
//
//	author centonhuang
//	update 2025-10-30
type OperationHandler interface {
	HandleUserLikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (*protocol.HumaHTTPResponse[*dto.LikeArticleResponse], error)
	HandleUserLikeComment(ctx context.Context, req *dto.LikeCommentRequest) (*protocol.HumaHTTPResponse[*dto.LikeCommentResponse], error)
	HandleUserLikeTag(ctx context.Context, req *dto.LikeTagRequest) (*protocol.HumaHTTPResponse[*dto.LikeTagResponse], error)
	HandleLogUserViewArticle(ctx context.Context, req *dto.LogArticleViewRequest) (*protocol.HumaHTTPResponse[*dto.LogArticleViewResponse], error)
}

type operationHandler struct {
	svc service.OperationService
}

// NewOperationHandler 创建用户操作处理器
//
//	return OperationHandler
//	author centonhuang
//	update 2025-10-30
func NewOperationHandler() OperationHandler {
	return &operationHandler{
		svc: service.NewOperationService(),
	}
}

func (h *operationHandler) HandleUserLikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (*protocol.HumaHTTPResponse[*dto.LikeArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.LikeArticleRequest{
		UserID:    userID,
		ArticleID: req.Body.ArticleID,
		Undo:      req.Body.Undo,
	}

	_, err := h.svc.LikeArticle(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.LikeArticleResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.LikeArticleResponse{}, nil)
}

func (h *operationHandler) HandleUserLikeComment(ctx context.Context, req *dto.LikeCommentRequest) (*protocol.HumaHTTPResponse[*dto.LikeCommentResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.LikeCommentRequest{
		UserID:    userID,
		CommentID: req.Body.CommentID,
		Undo:      req.Body.Undo,
	}

	_, err := h.svc.LikeComment(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.LikeCommentResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.LikeCommentResponse{}, nil)
}

func (h *operationHandler) HandleUserLikeTag(ctx context.Context, req *dto.LikeTagRequest) (*protocol.HumaHTTPResponse[*dto.LikeTagResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.LikeTagRequest{
		UserID: userID,
		TagID:  req.Body.TagID,
		Undo:   req.Body.Undo,
	}

	_, err := h.svc.LikeTag(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.LikeTagResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.LikeTagResponse{}, nil)
}

func (h *operationHandler) HandleLogUserViewArticle(ctx context.Context, req *dto.LogArticleViewRequest) (*protocol.HumaHTTPResponse[*dto.LogArticleViewResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.LogArticleViewRequest{
		UserID:    userID,
		ArticleID: req.Body.ArticleID,
		Progress:  req.Body.Progress,
	}

	_, err := h.svc.LogArticleView(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.LogArticleViewResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.LogArticleViewResponse{}, nil)
}
