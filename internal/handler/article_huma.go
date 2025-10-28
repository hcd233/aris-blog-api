package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// CreateArticleHuma 创建文章（Huma 版本）
func CreateArticleHuma(ctx context.Context, input *protocol.CreateArticleInput) (*protocol.ArticleOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.CreateArticleRequest{
		UserID:     userID,
		Title:      input.Body.Title,
		Slug:       input.Body.Slug,
		CategoryID: input.Body.CategoryID,
		Tags:       input.Body.Tags,
	}

	rsp, err := service.NewArticleService().CreateArticle(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.ArticleOutput{
		Body: *rsp.Article,
	}, nil
}

// GetArticleInfoHuma 获取文章信息（Huma 版本）
func GetArticleInfoHuma(ctx context.Context, input *protocol.ArticleInput) (*protocol.ArticleOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.GetArticleInfoRequest{
		UserID:    userID,
		ArticleID: input.ArticleID,
	}

	rsp, err := service.NewArticleService().GetArticleInfo(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.ArticleOutput{
		Body: *rsp.Article,
	}, nil
}

// GetArticleInfoBySlugHuma 通过 slug 获取文章信息（Huma 版本）
func GetArticleInfoBySlugHuma(ctx context.Context, input *protocol.ArticleSlugInput) (*protocol.ArticleOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.GetArticleInfoBySlugRequest{
		UserID:      userID,
		AuthorName:  input.AuthorName,
		ArticleSlug: input.ArticleSlug,
	}

	rsp, err := service.NewArticleService().GetArticleInfoBySlug(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.ArticleOutput{
		Body: *rsp.Article,
	}, nil
}

// UpdateArticleHuma 更新文章（Huma 版本）
func UpdateArticleHuma(ctx context.Context, input *protocol.UpdateArticleInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.UpdateArticleRequest{
		UserID:            userID,
		ArticleID:         input.ArticleID,
		UpdatedTitle:      input.Body.Title,
		UpdatedSlug:       input.Body.Slug,
		UpdatedCategoryID: input.Body.CategoryID,
	}

	_, err := service.NewArticleService().UpdateArticle(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// UpdateArticleStatusHuma 更新文章状态（Huma 版本）
func UpdateArticleStatusHuma(ctx context.Context, input *protocol.UpdateArticleStatusInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.UpdateArticleStatusRequest{
		UserID:    userID,
		ArticleID: input.ArticleID,
		Status:    protocol.GetArticleStatusFromString(input.Body.Status),
	}

	_, err := service.NewArticleService().UpdateArticleStatus(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// DeleteArticleHuma 删除文章（Huma 版本）
func DeleteArticleHuma(ctx context.Context, input *protocol.ArticleInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.DeleteArticleRequest{
		UserID:    userID,
		ArticleID: input.ArticleID,
	}

	_, err := service.NewArticleService().DeleteArticle(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// ListArticlesHuma 列出文章（Huma 版本）
func ListArticlesHuma(ctx context.Context, input *protocol.PaginatedSearchParams) (*protocol.ArticleListOutput, error) {
	req := &protocol.ListArticlesRequest{
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     input.Page,
				PageSize: input.PageSize,
			},
			QueryParam: &protocol.QueryParam{
				Query: input.Query,
			},
		},
	}

	rsp, err := service.NewArticleService().ListArticles(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	articles := make([]protocol.Article, len(rsp.Articles))
	for i, a := range rsp.Articles {
		articles[i] = *a
	}
	return &protocol.ArticleListOutput{
		Body: struct {
			Articles []protocol.Article `json:"articles"`
			PageInfo protocol.PageInfo  `json:"pageInfo"`
		}{
			Articles: articles,
			PageInfo: *rsp.PageInfo,
		},
	}, nil
}
