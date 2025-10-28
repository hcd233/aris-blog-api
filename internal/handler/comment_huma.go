package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// CreateCommentHuma 创建评论（Huma 版本）
func CreateCommentHuma(ctx context.Context, input *protocol.CreateCommentInput) (*protocol.CommentOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.CreateArticleCommentRequest{
		UserID:    userID,
		ArticleID: input.Body.ArticleID,
		Content:   input.Body.Content,
		ReplyTo:   input.Body.ReplyTo,
	}

	rsp, err := service.NewCommentService().CreateArticleComment(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.CommentOutput{
		Body: *rsp.Comment,
	}, nil
}

// DeleteCommentHuma 删除评论（Huma 版本）
func DeleteCommentHuma(ctx context.Context, input *protocol.CommentInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.DeleteCommentRequest{
		UserID:    userID,
		CommentID: input.CommentID,
	}

	_, err := service.NewCommentService().DeleteComment(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// ListArticleCommentsHuma 列出文章评论（Huma 版本）
func ListArticleCommentsHuma(ctx context.Context, input *protocol.ArticleCommentInput) (*protocol.CommentListOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.ListArticleCommentsRequest{
		UserID:    userID,
		ArticleID: input.ArticleID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     1,
				PageSize: 50,
			},
		},
	}

	rsp, err := service.NewCommentService().ListArticleComments(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	comments := make([]protocol.Comment, len(rsp.Comments))
	for i, c := range rsp.Comments {
		comments[i] = *c
	}
	return &protocol.CommentListOutput{
		Body: struct {
			Comments []protocol.Comment `json:"comments"`
			PageInfo protocol.PageInfo  `json:"pageInfo"`
		}{
			Comments: comments,
			PageInfo: *rsp.PageInfo,
		},
	}, nil
}

// ListChildrenCommentsHuma 列出子评论（Huma 版本）
func ListChildrenCommentsHuma(ctx context.Context, input *protocol.CommentInput) (*protocol.CommentListOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.ListChildrenCommentsRequest{
		UserID:    userID,
		CommentID: input.CommentID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     1,
				PageSize: 50,
			},
		},
	}

	rsp, err := service.NewCommentService().ListChildrenComments(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	comments := make([]protocol.Comment, len(rsp.Comments))
	for i, c := range rsp.Comments {
		comments[i] = *c
	}
	return &protocol.CommentListOutput{
		Body: struct {
			Comments []protocol.Comment `json:"comments"`
			PageInfo protocol.PageInfo  `json:"pageInfo"`
		}{
			Comments: comments,
			PageInfo: *rsp.PageInfo,
		},
	}, nil
}
