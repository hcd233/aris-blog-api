package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// OperationHandler 用户操作处理器
type OperationHandler interface {
	HandleUserLikeArticle(c *gin.Context)
	HandleUserLikeComment(c *gin.Context)
	HandleUserLikeTag(c *gin.Context)
	HandleLogUserViewArticle(c *gin.Context)
}

type operationHandler struct {
	svc service.OperationService
}

// NewOperationHandler 创建用户操作处理器
func NewOperationHandler() OperationHandler {
	return &operationHandler{
		svc: service.NewOperationService(),
	}
}

// HandleUserLikeArticle 点赞文章
func (h *operationHandler) HandleUserLikeArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LikeArticleBody)

	req := &protocol.LikeArticleRequest{
		CurUserID:   userID,
		Author:      body.Author,
		ArticleSlug: body.ArticleSlug,
		Undo:        body.Undo,
	}

	rsp, err := h.svc.LikeArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUserLikeComment 点赞评论
func (h *operationHandler) HandleUserLikeComment(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LikeCommentBody)

	req := &protocol.LikeCommentRequest{
		CurUserID: userID,
		CommentID: body.CommentID,
		Undo:      body.Undo,
	}

	rsp, err := h.svc.LikeComment(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUserLikeTag 点赞标签
func (h *operationHandler) HandleUserLikeTag(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LikeTagBody)

	req := &protocol.LikeTagRequest{
		CurUserID: userID,
		TagSlug:   body.TagSlug,
		Undo:      body.Undo,
	}

	rsp, err := h.svc.LikeTag(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleLogUserViewArticle 记录文章浏览
func (h *operationHandler) HandleLogUserViewArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LogUserViewArticleBody)

	req := &protocol.LogArticleViewRequest{
		CurUserID:   userID,
		Author:      body.Author,
		ArticleSlug: body.ArticleSlug,
		Progress:    body.Progress,
	}

	rsp, err := h.svc.LogArticleView(req)

	util.SendHTTPResponse(c, rsp, err)
}
