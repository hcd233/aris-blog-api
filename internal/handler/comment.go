package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CommentHandler 评论处理器
type CommentHandler interface {
	HandleCreateArticleComment(ctx context.Context, req *dto.CreateCommentRequest) (*protocol.HTTPResponse[*dto.CreateCommentResponse], error)
	HandleDeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error)
	HandleListArticleComments(ctx context.Context, req *dto.ListArticleCommentRequest) (*protocol.HTTPResponse[*dto.ListArticleCommentResponse], error)
	HandleListChildrenComments(ctx context.Context, req *dto.ListChildrenCommentRequest) (*protocol.HTTPResponse[*dto.ListChildrenCommentResponse], error)
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

func (h *commentHandler) HandleCreateArticleComment(ctx context.Context, req *dto.CreateCommentRequest) (*protocol.HTTPResponse[*dto.CreateCommentResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreateArticleComment(ctx, req))
}

func (h *commentHandler) HandleDeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteComment(ctx, req))
}

func (h *commentHandler) HandleListArticleComments(ctx context.Context, req *dto.ListArticleCommentRequest) (*protocol.HTTPResponse[*dto.ListArticleCommentResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListArticleComments(ctx, req))
}

func (h *commentHandler) HandleListChildrenComments(ctx context.Context, req *dto.ListChildrenCommentRequest) (*protocol.HTTPResponse[*dto.ListChildrenCommentResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListChildrenComments(ctx, req))
}
