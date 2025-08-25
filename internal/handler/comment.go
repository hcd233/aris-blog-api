package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CommentHandler 评论处理器
type CommentHandler interface {
	HandleCreateArticleComment(c *fiber.Ctx) error
	HandleDeleteComment(c *fiber.Ctx) error
	HandleListArticleComments(c *fiber.Ctx) error
	HandleListChildrenComments(c *fiber.Ctx) error
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

// HandleCreateArticleComment 创建文章评论
//
//	@Summary 创建文章评论
//	@Description 创建文章评论
//	@Tags comment
//	@Accept json
//	@Produce json
//	@Param body body protocol.CreateArticleCommentBody true "创建文章评论请求"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.HTTPResponse{data=protocol.CreateArticleCommentResponse,error=nil} "创建文章评论响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/comment [post]
//	receiver h *commentHandler
func (h *commentHandler) HandleCreateArticleComment(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	body := c.Locals(constant.CtxKeyBody).(*protocol.CreateArticleCommentBody)

	req := &protocol.CreateArticleCommentRequest{
		UserID:    userID,
		ArticleID: body.ArticleID,
		Content:   body.Content,
		ReplyTo:   body.ReplyTo,
	}

	rsp, err := h.svc.CreateArticleComment(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleDeleteComment 删除评论
//
//	@Summary 删除评论
//	@Description 删除评论
//	@Tags comment
//	@Accept json
//	@Produce json
//	@Param path path protocol.CommentURI true "评论ID"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.HTTPResponse{data=protocol.DeleteCommentResponse,error=nil} "删除评论响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/comment/{commentID} [delete]
func (h *commentHandler) HandleDeleteComment(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.CommentURI)

	req := &protocol.DeleteCommentRequest{
		UserID:    userID,
		CommentID: uri.CommentID,
	}

	rsp, err := h.svc.DeleteComment(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleListArticleComments 列出文章一级评论
//
//	@Summary 列出文章一级评论
//	@Description 列出文章一级评论
//	@Tags comment
//	@Accept json
//	@Produce json
//	@Param path path protocol.ArticleURI true "文章ID"
//	@Param param query protocol.PageParam true "分页参数"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.HTTPResponse{data=protocol.ListArticleCommentsResponse,error=nil} "列出文章评论响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/comment/article/{articleID}/list [get]
func (h *commentHandler) HandleListArticleComments(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.ArticleURI)
	param := c.Locals(constant.CtxKeyParam).(*protocol.PaginateParam)

	req := &protocol.ListArticleCommentsRequest{
		UserID:         userID,
		ArticleID:      uri.ArticleID,
		PaginateParam:  param,
	}

	rsp, err := h.svc.ListArticleComments(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleListChildrenComments 列出子评论
//
//	@Summary 列出子评论
//	@Description 列出子评论
//	@Tags comment
//	@Accept json
//	@Produce json
//	@Param path path protocol.CommentURI true "评论ID"
//	@Param param query protocol.PageParam true "分页参数"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.HTTPResponse{data=protocol.ListChildrenCommentsResponse,error=nil} "列出子评论响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/comment/{commentID}/subComments [get]
func (h *commentHandler) HandleListChildrenComments(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.CommentURI)
	param := c.Locals(constant.CtxKeyParam).(*protocol.PaginateParam)

	req := &protocol.ListChildrenCommentsRequest{
		UserID:         userID,
		CommentID:      uri.CommentID,
		PaginateParam:  param,
	}

	rsp, err := h.svc.ListChildrenComments(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}
