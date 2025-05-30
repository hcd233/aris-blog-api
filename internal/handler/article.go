package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleHandler 文章处理器
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ArticleHandler interface {
	HandleCreateArticle(c *gin.Context)
	HandleGetArticleInfo(c *gin.Context)
	HandleUpdateArticle(c *gin.Context)
	HandleUpdateArticleStatus(c *gin.Context)
	HandleDeleteArticle(c *gin.Context)
	HandleListArticles(c *gin.Context)
	HandleGetArticleInfoBySlug(c *gin.Context)
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
//	@Summary 创建文章
//	@Description 创建文章
//	@Tags article
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.CreateArticleBody	true	"创建文章请求"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.CreateArticleResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/article [post]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleCreateArticle(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	body := c.MustGet(constant.CtxKeyBody).(*protocol.CreateArticleBody)

	req := &protocol.CreateArticleRequest{
		UserID:     userID,
		Title:      body.Title,
		Slug:       body.Slug,
		CategoryID: body.CategoryID,
		Tags:       body.Tags,
	}

	rsp, err := h.svc.CreateArticle(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetArticleInfo 获取文章信息
//
//	@Summary 获取文章信息
//	@Description 获取文章信息
//	@Tags article
//	@Accept json
//	@Produce json
//	@Param path path protocol.ArticleURI true "文章ID"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.GetArticleInfoResponse "获取文章信息响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID} [get]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleGetArticleInfo(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	uri := c.MustGet(constant.CtxKeyURI).(*protocol.ArticleURI)

	req := &protocol.GetArticleInfoRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
	}

	rsp, err := h.svc.GetArticleInfo(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetArticleInfoBySlug 获取文章信息
//
//	@Summary 获取文章信息
//	@Description 获取文章信息
//	@Tags article
//	@Accept json
//	@Produce json
//	@Param path path protocol.ArticleSlugURI true "作者名和文章别名"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.GetArticleInfoResponse "获取文章信息响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/slug/{authorName}/{articleSlug} [get]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-19 15:23:26
func (h *articleHandler) HandleGetArticleInfoBySlug(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	uri := c.MustGet(constant.CtxKeyURI).(*protocol.ArticleSlugURI)

	req := &protocol.GetArticleInfoBySlugRequest{
		UserID:      userID,
		AuthorName:  uri.AuthorName,
		ArticleSlug: uri.ArticleSlug,
	}

	rsp, err := h.svc.GetArticleInfoBySlug(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateArticle 更新文章
//
//	@Summary 更新文章
//	@Description 更新文章
//	@Tags article
//	@Accept json
//	@Produce json
//	@Param path path protocol.ArticleURI true "文章ID"
//	@Param body body protocol.UpdateArticleBody true "更新文章请求"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.UpdateArticleResponse "更新文章响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID} [patch]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleUpdateArticle(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	uri := c.MustGet(constant.CtxKeyURI).(*protocol.ArticleURI)
	body := c.MustGet(constant.CtxKeyBody).(*protocol.UpdateArticleBody)

	req := &protocol.UpdateArticleRequest{
		UserID:            userID,
		ArticleID:         uri.ArticleID,
		UpdatedTitle:      body.Title,
		UpdatedSlug:       body.Slug,
		UpdatedCategoryID: body.CategoryID,
	}

	rsp, err := h.svc.UpdateArticle(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateArticleStatus 更新文章状态
//
//	@Summary 更新文章状态
//	@Description 更新文章状态
//	@Tags article
//	@Accept json
//	@Produce json
//	@Param path path protocol.ArticleURI true "文章ID"
//	@Param body body protocol.UpdateArticleStatusBody true "更新文章状态请求"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.UpdateArticleStatusResponse "更新文章状态响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID}/status [put]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleUpdateArticleStatus(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	uri := c.MustGet(constant.CtxKeyURI).(*protocol.ArticleURI)
	body := c.MustGet(constant.CtxKeyBody).(*protocol.UpdateArticleStatusBody)

	req := &protocol.UpdateArticleStatusRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		Status:    body.Status,
	}

	rsp, err := h.svc.UpdateArticleStatus(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteArticle 删除文章
//
//	@Summary 删除文章
//	@Description 删除文章
//	@Tags article
//	@Accept json
//	@Produce json
//	@Param path path protocol.ArticleURI true "文章ID"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.DeleteArticleResponse "删除文章响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID} [delete]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleDeleteArticle(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	uri := c.MustGet(constant.CtxKeyURI).(*protocol.ArticleURI)

	req := &protocol.DeleteArticleRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
	}

	rsp, err := h.svc.DeleteArticle(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListArticles 列出文章
//
//	@Summary 列出文章
//	@Description 列出文章
//	@Tags article
//	@Accept json
//	@Produce json
//	@Param param query protocol.PageParam true "分页参数"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.ListArticlesResponse "列出文章响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/list [get]
//	receiver h *articleHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleHandler) HandleListArticles(c *gin.Context) {
	param := c.MustGet(constant.CtxKeyParam).(*protocol.PageParam)

	req := &protocol.ListArticlesRequest{
		PageParam: param,
	}

	rsp, err := h.svc.ListArticles(c, req)

	util.SendHTTPResponse(c, rsp, err)
}
