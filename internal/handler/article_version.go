package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleVersionHandler 文章版本处理器
type ArticleVersionHandler interface {
	HandleCreateArticleVersion(c *fiber.Ctx) error
	HandleGetArticleVersionInfo(c *fiber.Ctx) error
	HandleGetLatestArticleVersionInfo(c *fiber.Ctx) error
	HandleListArticleVersions(c *fiber.Ctx) error
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
//	@Param uri path protocol.ArticleURI true "文章路径参数"
//	@Param body body protocol.CreateArticleVersionBody true "创建文章版本请求体"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.CreateArticleVersionResponse "创建文章版本响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID}/version [post]
//	receiver h *articleVersionHandler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleCreateArticleVersion(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.ArticleURI)
	body := c.Locals(constant.CtxKeyBody).(*protocol.CreateArticleVersionBody)

	req := &protocol.CreateArticleVersionRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		Content:   body.Content,
	}

	rsp, err := h.svc.CreateArticleVersion(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleGetArticleVersionInfo 获取文章版本信息
//
//	@Summary 获取文章版本信息
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param uri path protocol.ArticleVersionURI true "文章版本路径参数"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.GetArticleVersionInfoResponse "获取文章版本信息响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID}/version/v{version} [get]
//	receiver h *articleVersionHandler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleGetArticleVersionInfo(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.ArticleVersionURI)

	req := &protocol.GetArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		VersionID: uri.Version,
	}

	rsp, err := h.svc.GetArticleVersionInfo(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleGetLatestArticleVersionInfo 获取最新文章版本信息
//
//	@Summary 获取最新文章版本信息
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param uri path protocol.ArticleURI true "文章路径参数"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.GetLatestArticleVersionInfoResponse "获取最新文章版本信息响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID}/version/latest [get]
//	receiver h *articleVersionHandler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.ArticleURI)

	req := &protocol.GetLatestArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
	}

	rsp, err := h.svc.GetLatestArticleVersionInfo(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleListArticleVersions 列出文章版本
//
//	@Summary 列出文章版本
//	@Tags articleVersion
//	@Accept json
//	@Produce json
//	@Param uri path protocol.ArticleURI true "文章路径参数"
//	@Param param query protocol.PageParam true "分页参数"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.ListArticleVersionsResponse "列出文章版本响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/article/{articleID}/version/list [get]
//	receiver h *articleVersionHandler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (h *articleVersionHandler) HandleListArticleVersions(c *fiber.Ctx) error {
		userID := c.Locals(constant.CtxKeyUserID).(uint).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.ArticleURI)
	param := c.Locals(constant.CtxKeyParam).(*protocol.PageParam)

	req := &protocol.ListArticleVersionsRequest{
		UserID:     userID,
		ArticleID:  uri.ArticleID,
		PageParam:  param,
	}

	rsp, err := h.svc.ListArticleVersions(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}
