package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleHandler 文章处理器
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleResponse], error)
	HandleGetArticleInfo(ctx context.Context, req *dto.GetArticleInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoResponse], error)
	HandleGetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleInfoBySlugRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoBySlugResponse], error)
	HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleResponse], error)
	HandleUpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleStatusResponse], error)
	HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (*protocol.HumaHTTPResponse[*dto.DeleteArticleResponse], error)
	HandleListArticles(ctx context.Context, req *dto.ListArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListArticlesResponse], error)
}

type articleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 创建文章处理器
//
//	return ArticleHandler
//	author centonhuang
//	update 2025-01-05 15:23:26
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		svc: service.NewArticleService(),
	}
}

// HandleCreateArticle 创建文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.CreateArticleRequest
//	return *protocol.HumaHTTPResponse[*dto.CreateArticleResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.CreateArticleRequest{
		UserID:     userID,
		Title:      req.Body.Title,
		Slug:       req.Body.Slug,
		CategoryID: req.Body.CategoryID,
		Tags:       req.Body.Tags,
	}

	serviceRsp, err := h.svc.CreateArticle(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.CreateArticleResponse{
		Article: convertArticle(serviceRsp.Article),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleGetArticleInfo 获取文章信息
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.GetArticleInfoRequest
//	return *protocol.HumaHTTPResponse[*dto.GetArticleInfoResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleGetArticleInfo(ctx context.Context, req *dto.GetArticleInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.GetArticleInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	serviceRsp, err := h.svc.GetArticleInfo(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.GetArticleInfoResponse{
		Article: convertArticle(serviceRsp.Article),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleGetArticleInfoBySlug 通过slug获取文章信息
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.GetArticleInfoBySlugRequest
//	return *protocol.HumaHTTPResponse[*dto.GetArticleInfoBySlugResponse]
//	return error
//	author centonhuang
//	update 2025-01-19 15:23:26
func (h *articleHandler) HandleGetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleInfoBySlugRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoBySlugResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.GetArticleInfoBySlugRequest{
		UserID:      userID,
		AuthorName:  req.AuthorName,
		ArticleSlug: req.ArticleSlug,
	}

	serviceRsp, err := h.svc.GetArticleInfoBySlug(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.GetArticleInfoBySlugResponse{
		Article: convertArticle(serviceRsp.Article),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleUpdateArticle 更新文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.UpdateArticleRequest
//	return *protocol.HumaHTTPResponse[*dto.UpdateArticleResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.UpdateArticleRequest{
		UserID:            userID,
		ArticleID:         req.ArticleID,
		UpdatedTitle:      req.Body.Title,
		UpdatedSlug:       req.Body.Slug,
		UpdatedCategoryID: req.Body.CategoryID,
	}

	_, err := h.svc.UpdateArticle(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.UpdateArticleResponse{}
	return util.WrapHTTPResponse(rsp, nil)
}

// HandleUpdateArticleStatus 更新文章状态
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.UpdateArticleStatusRequest
//	return *protocol.HumaHTTPResponse[*dto.UpdateArticleStatusResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleUpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleStatusResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.UpdateArticleStatusRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Status:    req.Body.Status,
	}

	_, err := h.svc.UpdateArticleStatus(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.UpdateArticleStatusResponse{}
	return util.WrapHTTPResponse(rsp, nil)
}

// HandleDeleteArticle 删除文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.DeleteArticleRequest
//	return *protocol.HumaHTTPResponse[*dto.DeleteArticleResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (*protocol.HumaHTTPResponse[*dto.DeleteArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.DeleteArticleRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	_, err := h.svc.DeleteArticle(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.DeleteArticleResponse{}
	return util.WrapHTTPResponse(rsp, nil)
}

// HandleListArticles 列出文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *dto.ListArticlesRequest
//	return *protocol.HumaHTTPResponse[*dto.ListArticlesResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleListArticles(ctx context.Context, req *dto.ListArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListArticlesResponse], error) {
	serviceReq := &protocol.ListArticlesRequest{
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     req.Page,
				PageSize: req.PageSize,
			},
		},
	}

	serviceRsp, err := h.svc.ListArticles(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	articles := make([]*dto.Article, len(serviceRsp.Articles))
	for i, article := range serviceRsp.Articles {
		articles[i] = convertArticle(article)
	}

	rsp := &dto.ListArticlesResponse{
		Articles: articles,
		PageInfo: &dto.PageInfo{
			Page:     serviceRsp.PageInfo.Page,
			PageSize: serviceRsp.PageInfo.PageSize,
			Total:    serviceRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// convertArticle 转换文章模型
func convertArticle(article *protocol.Article) *dto.Article {
	if article == nil {
		return nil
	}

	tags := make([]*dto.Tag, len(article.Tags))
	for i, tag := range article.Tags {
		tags[i] = &dto.Tag{
			TagID:       tag.TagID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			UserID:      tag.UserID,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
			Likes:       tag.Likes,
		}
	}

	var category *dto.Category
	if article.Category != nil {
		category = &dto.Category{
			CategoryID: article.Category.CategoryID,
			Name:       article.Category.Name,
			ParentID:   article.Category.ParentID,
			CreatedAt:  article.Category.CreatedAt,
			UpdatedAt:  article.Category.UpdatedAt,
		}
	}

	var user *dto.User
	if article.User != nil {
		user = &dto.User{
			UserID:    article.User.UserID,
			Name:      article.User.Name,
			Email:     article.User.Email,
			Avatar:    article.User.Avatar,
			CreatedAt: article.User.CreatedAt,
			LastLogin: article.User.LastLogin,
		}
	}

	return &dto.Article{
		ArticleID:   article.ArticleID,
		Title:       article.Title,
		Slug:        article.Slug,
		Status:      article.Status,
		User:        user,
		Category:    category,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		PublishedAt: article.PublishedAt,
		Likes:       article.Likes,
		Views:       article.Views,
		Tags:        tags,
		Comments:    article.Comments,
	}
}
