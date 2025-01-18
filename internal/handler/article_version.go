package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
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
//
//	@Summary 创建文章版本
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param body body protocol.CreateArticleVersionBody true "创建文章版本请求体"
//	@Success 200 {object} protocol.CreateArticleVersionResponse "创建文章版本响应"
//	@Router /v1/article/version [post]
//	receiver h *articleVersionHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleCreateArticleVersion(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	body := c.MustGet("body").(*protocol.CreateArticleVersionBody)

	req := &protocol.CreateArticleVersionRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		Content:   body.Content,
	}

	rsp, err := h.svc.CreateArticleVersion(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetArticleVersionInfo 获取文章版本信息
//
//	@Summary 获取文章版本信息
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param uri path protocol.ArticleVersionURI true "文章版本路径参数"
//	@Success 200 {object} protocol.GetArticleVersionInfoResponse "获取文章版本信息响应"
//	@Router /v1/article/version/v{version} [get]
//	receiver h *articleVersionHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleGetArticleVersionInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleVersionURI)

	req := &protocol.GetArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		VersionID: uri.Version,
	}

	rsp, err := h.svc.GetArticleVersionInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetLatestArticleVersionInfo 获取最新文章版本信息
//
//	@Summary 获取最新文章版本信息
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param uri path protocol.ArticleURI true "文章路径参数"
//	@Success 200 {object} protocol.GetLatestArticleVersionInfoResponse "获取最新文章版本信息响应"
//	@Router /v1/article/version/latest [get]
//	receiver h *articleVersionHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleURI)

	req := &protocol.GetLatestArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
	}

	rsp, err := h.svc.GetLatestArticleVersionInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListArticleVersions 列出文章版本
//
//	@Summary 列出文章版本
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param uri path protocol.ArticleURI true "文章路径参数"
//	@Param param query protocol.PageParam true "分页参数"
//	@Success 200 {object} protocol.ListArticleVersionsResponse "列出文章版本响应"
//	@Router /v1/article/version/versions [get]
//	receiver h *articleVersionHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleListArticleVersions(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListArticleVersionsRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		PageParam: param,
	}

	rsp, err := h.svc.ListArticleVersions(req)

	util.SendHTTPResponse(c, rsp, err)
}
