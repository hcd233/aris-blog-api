package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CommentHandler 评论处理器
type CommentHandler interface {
	HandleCreateArticleComment(c *gin.Context)
	HandleGetCommentInfo(c *gin.Context)
	HandleDeleteComment(c *gin.Context)
	HandleListArticleComments(c *gin.Context)
	HandleListChildrenComments(c *gin.Context)
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
func (h *commentHandler) HandleCreateArticleComment(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	body := c.MustGet("body").(*protocol.CreateArticleCommentBody)

	req := &protocol.CreateArticleCommentRequest{
		CurUserID:   userID,
		Author:      uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		Content:     body.Content,
		ReplyTo:     body.ReplyTo,
	}

	rsp, err := h.svc.CreateArticleComment(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetCommentInfo 获取评论信息
func (h *commentHandler) HandleGetCommentInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)

	req := &protocol.GetCommentInfoRequest{
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		CommentID:   uri.CommentID,
	}

	rsp, err := h.svc.GetCommentInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteComment 删除评论
func (h *commentHandler) HandleDeleteComment(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CommentURI)

	req := &protocol.DeleteCommentRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		CommentID:   uri.CommentID,
	}

	rsp, err := h.svc.DeleteComment(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListArticleComments 列出文章评论
func (h *commentHandler) HandleListArticleComments(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListArticleCommentsRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		PageParam:   param,
	}

	rsp, err := h.svc.ListArticleComments(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListChildrenComments 列出子评论
func (h *commentHandler) HandleListChildrenComments(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListChildrenCommentsRequest{
		UserName:    uri.UserName,
		ArticleSlug: uri.ArticleSlug,
		CommentID:   uri.CommentID,
		PageParam:   param,
	}

	rsp, err := h.svc.ListChildrenComments(req)

	util.SendHTTPResponse(c, rsp, err)
}
