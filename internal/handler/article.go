package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
)

// ArticleHandler 文章处理器
//
//	@author centonhuang
//	@update 2025-01-05 15:23:26
type ArticleHandler interface {
	HandleCreateArticle(c *gin.Context)
	HandleGetArticleInfo(c *gin.Context)
	HandleUpdateArticle(c *gin.Context)
	HandleUpdateArticleStatus(c *gin.Context)
	HandleDeleteArticle(c *gin.Context)
	HandleListArticles(c *gin.Context)
	HandleListUserArticles(c *gin.Context)
	HandleQueryArticle(c *gin.Context)
	HandleQueryUserArticle(c *gin.Context)
}

type articleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 创建文章处理器
//
//	@return ArticleHandler
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		svc: service.NewArticleService(),
	}
}

// HandleCreateArticle 创建文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleCreateArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateArticleBody)

	req := &protocol.CreateArticleRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		UserID:      userID,
		Title:       body.Title,
		Slug:        body.Slug,
		CategoryID:  body.CategoryID,
		Tags:        body.Tags,
	}

	rsp, err := h.svc.CreateArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetArticleInfo 获取文章信息
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleGetArticleInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	req := &protocol.GetArticleInfoRequest{
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
	}

	rsp, err := h.svc.GetArticleInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateArticle 更新文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleUpdateArticle(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.UpdateArticleBody)

	req := &protocol.UpdateArticleRequest{
		CurUserName:       userName,
		UserName:          uri.UserName,
		ArticleSlug:       uri.ArticleSlug,
		UpdatedTitle:      body.Title,
		UpdatedSlug:       body.Slug,
		UpdatedCategoryID: body.CategoryID,
	}

	rsp, err := h.svc.UpdateArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateArticleStatus 更新文章状态
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleUpdateArticleStatus(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.UpdateArticleStatusBody)

	req := &protocol.UpdateArticleStatusRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		Status:      body.Status,
	}

	rsp, err := h.svc.UpdateArticleStatus(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteArticle 删除文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleDeleteArticle(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	req := &protocol.DeleteArticleRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
	}

	rsp, err := h.svc.DeleteArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListArticles 列出文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleListArticles(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListArticlesRequest{
		PageParam: param,
	}

	rsp, err := h.svc.ListArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserArticles 列出用户文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleListUserArticles(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserArticlesRequest{
		UserName:  uri.UserName,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleQueryArticle 查询文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleQueryArticle(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	req := &protocol.QueryArticleRequest{
		QueryParam: param,
	}

	rsp, err := h.svc.QueryArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleQueryUserArticle 查询用户文章
//
//	@receiver h *articleHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 15:23:26
func (h *articleHandler) HandleQueryUserArticle(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.QueryParam)

	req := &protocol.QueryUserArticleRequest{
		UserName:   uri.UserName,
		QueryParam: param,
	}

	rsp, err := h.svc.QueryUserArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}
