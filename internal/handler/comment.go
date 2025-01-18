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
func (h *commentHandler) HandleCreateArticleComment(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.CreateArticleCommentBody)

	req := &protocol.CreateArticleCommentRequest{
		UserID:    userID,
		ArticleID: body.ArticleID,
		Content:   body.Content,
		ReplyTo:   body.ReplyTo,
	}

	rsp, err := h.svc.CreateArticleComment(req)

	util.SendHTTPResponse(c, rsp, err)
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
func (h *commentHandler) HandleDeleteComment(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.CommentURI)

	req := &protocol.DeleteCommentRequest{
		UserID:    userID,
		CommentID: uri.CommentID,
	}

	rsp, err := h.svc.DeleteComment(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListArticleComments 列出文章评论
func (h *commentHandler) HandleListArticleComments(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListArticleCommentsRequest{
		UserID:    userID,
		ArticleID: uri.ArticleID,
		PageParam: param,
	}

	rsp, err := h.svc.ListArticleComments(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListChildrenComments 列出子评论
//
//	@Summary 列出子评论
//	@Description 列出子评论
//	@Tags comment
//	@Accept json
//	@Produce json
//	@Param path path protocol.CommentURI true "评论ID"
//	@Security ApiKeyAuth
//	@Success 200 {object} protocol.HTTPResponse{data=protocol.ListChildrenCommentsResponse,error=nil} "列出子评论响应"
//	@Failure 400 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 401 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 403 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Failure 500 {object} protocol.HTTPResponse{data=nil,error=string}
//	@Router /v1/comment/{commentID}/subComments [get]
func (h *commentHandler) HandleListChildrenComments(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.CommentURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListChildrenCommentsRequest{
		UserID:    userID,
		CommentID: uri.CommentID,
		PageParam: param,
	}

	rsp, err := h.svc.ListChildrenComments(req)

	util.SendHTTPResponse(c, rsp, err)
}
