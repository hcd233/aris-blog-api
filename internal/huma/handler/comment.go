package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type CommentHandlers struct{ svc service.CommentService }

func NewCommentHandlers() *CommentHandlers { return &CommentHandlers{svc: service.NewCommentService()} }

type (
	createCommentInput struct {
		authHeader
		humadto.CreateArticleCommentInput
	}
	commentPathInput struct {
		authHeader
		humadto.CommentPath
	}
	articlePathWithPageInput struct {
		authHeader
		humadto.ArticlePath
		humadto.PaginateParam
	}
	commentPathWithPageInput struct {
		authHeader
		humadto.CommentPath
		humadto.PaginateParam
	}
)

func (h *CommentHandlers) HandleCreateArticleComment(ctx context.Context, input *createCommentInput) (*humadto.Output[protocol.CreateArticleCommentResponse], error) {
	req := &protocol.CreateArticleCommentRequest{UserID: input.UserID, ArticleID: input.Body.ArticleID, Content: input.Body.Content, ReplyTo: input.Body.ReplyTo}
	rsp, err := h.svc.CreateArticleComment(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CreateArticleCommentResponse]{Body: *rsp}, nil
}

func (h *CommentHandlers) HandleDeleteComment(ctx context.Context, input *commentPathInput) (*humadto.Output[protocol.DeleteCommentResponse], error) {
	req := &protocol.DeleteCommentRequest{UserID: input.UserID, CommentID: input.CommentID}
	rsp, err := h.svc.DeleteComment(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.DeleteCommentResponse]{Body: *rsp}, nil
}

func (h *CommentHandlers) HandleListArticleComments(ctx context.Context, input *articlePathWithPageInput) (*humadto.Output[protocol.ListArticleCommentsResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListArticleCommentsRequest{UserID: input.UserID, ArticleID: input.ArticleID, PaginateParam: p}
	rsp, err := h.svc.ListArticleComments(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListArticleCommentsResponse]{Body: *rsp}, nil
}

func (h *CommentHandlers) HandleListChildrenComments(ctx context.Context, input *commentPathWithPageInput) (*humadto.Output[protocol.ListChildrenCommentsResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListChildrenCommentsRequest{UserID: input.UserID, CommentID: input.CommentID, PaginateParam: p}
	rsp, err := h.svc.ListChildrenComments(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListChildrenCommentsResponse]{Body: *rsp}, nil
}
