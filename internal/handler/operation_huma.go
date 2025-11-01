package handler

import (
	"context"

	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// OperationHandlerForHuma 用户操作处理器（Huma版本）
type OperationHandlerForHuma interface {
	HandleUserLikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleUserLikeComment(ctx context.Context, req *dto.LikeCommentRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleUserLikeTag(ctx context.Context, req *dto.LikeTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleLogUserViewArticle(ctx context.Context, req *dto.LogArticleViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
}

type operationHandlerForHuma struct {
	svc service.OperationService
}

// NewOperationHandlerForHuma 创建用户操作处理器（Huma版本）
func NewOperationHandlerForHuma() OperationHandlerForHuma {
	return &operationHandlerForHuma{
		svc: service.NewOperationService(),
	}
}

func (h *operationHandlerForHuma) HandleUserLikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LikeArticle(ctx, req))
}

func (h *operationHandlerForHuma) HandleUserLikeComment(ctx context.Context, req *dto.LikeCommentRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LikeComment(ctx, req))
}

func (h *operationHandlerForHuma) HandleUserLikeTag(ctx context.Context, req *dto.LikeTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LikeTag(ctx, req))
}

func (h *operationHandlerForHuma) HandleLogUserViewArticle(ctx context.Context, req *dto.LogArticleViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.LogArticleView(ctx, req))
}
