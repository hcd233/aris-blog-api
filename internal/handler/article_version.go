package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
)

// ArticleVersionHandler 文章版本处理器
type ArticleVersionHandler interface {
	HandleCreateArticleVersion(c *gin.Context)
	HandleGetArticleVersionInfo(c *gin.Context)
	HandleGetLatestArticleVersionInfo(c *gin.Context)
	HandleListArticleVersions(c *gin.Context)
}

type articleVersionHandler struct {
	svc service.ArticleVersionService
}

// NewArticleVersionHandler 创建文章版本处理器
func NewArticleVersionHandler() ArticleVersionHandler {
	return &articleVersionHandler{
		svc: service.NewArticleVersionService(),
	}
}

// HandleCreateArticleVersion 创建文章版本
func (h *articleVersionHandler) HandleCreateArticleVersion(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.CreateArticleVersionBody)

	req := &protocol.CreateArticleVersionRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		Content:     body.Content,
	}

	rsp, err := h.svc.CreateArticleVersion(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetArticleVersionInfo 获取文章版本信息
func (h *articleVersionHandler) HandleGetArticleVersionInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleVersionURI)

	req := &protocol.GetArticleVersionInfoRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		Version:     uri.Version,
	}

	rsp, err := h.svc.GetArticleVersionInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetLatestArticleVersionInfo 获取最新文章版本信息
func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	req := &protocol.GetLatestArticleVersionInfoRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
	}

	rsp, err := h.svc.GetLatestArticleVersionInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListArticleVersions 列出文章版本
func (h *articleVersionHandler) HandleListArticleVersions(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListArticleVersionsRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		PageParam:   param,
	}

	rsp, err := h.svc.ListArticleVersions(req)

	util.SendHTTPResponse(c, rsp, err)
}
