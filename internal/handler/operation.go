package handler

import (
	"context"

	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// OperationHandler 用户操作处理器
type OperationHandler interface {
	HandleUserLikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleUserLikeComment(ctx context.Context, req *dto.LikeCommentRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleUserLikeTag(ctx context.Context, req *dto.LikeTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleLogUserViewArticle(ctx context.Context, req *dto.LogArticleViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
}

type operationHandler struct {
	svc service.OperationService
}

// NewOperationHandler 创建用户操作处理器
func NewOperationHandler() OperationHandler {
	return &operationHandler{
		svc: service.NewOperationService(),
	}
}

func (h *operationHandler) HandleUserLikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LikeArticle(ctx, req))
}

func (h *operationHandler) HandleUserLikeComment(ctx context.Context, req *dto.LikeCommentRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LikeComment(ctx, req))
}

func (h *operationHandler) HandleUserLikeTag(ctx context.Context, req *dto.LikeTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LikeTag(ctx, req))
}

func (h *operationHandler) HandleLogUserViewArticle(ctx context.Context, req *dto.LogArticleViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LogArticleView(ctx, req))
}
