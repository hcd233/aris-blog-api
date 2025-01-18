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
//
//	@Summary 点赞文章
//	@Description 点赞文章
//	@Tags			operation
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		protocol.LikeArticleBody	true	"点赞文章请求体"
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListChildrenCategoriesResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/operation/like/article [post]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-01 05:09:47
func (h *operationHandler) HandleUserLikeArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LikeArticleBody)

	req := &protocol.LikeArticleRequest{
		UserID:    userID,
		ArticleID: body.ArticleID,
		Undo:      body.Undo,
	}

	rsp, err := h.svc.LikeArticle(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUserLikeComment 点赞评论
//
//	@Summary		点赞评论
//	@Description	点赞评论
//	@Tags			operation
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.LikeCommentBody	true	"点赞评论请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.LikeCommentResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/operation/like/comment [post]
func (h *operationHandler) HandleUserLikeComment(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LikeCommentBody)

	req := &protocol.LikeCommentRequest{
		UserID:    userID,
		CommentID: body.CommentID,
		Undo:      body.Undo,
	}

	rsp, err := h.svc.LikeComment(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUserLikeTag 点赞标签
//
//	@Summary		点赞标签
//	@Description	点赞标签
//	@Tags			operation
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.LikeTagBody	true	"点赞标签请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.LikeTagResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/operation/like/tag [post]
func (h *operationHandler) HandleUserLikeTag(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LikeTagBody)

	req := &protocol.LikeTagRequest{
		UserID: userID,
		TagID:  body.TagID,
		Undo:   body.Undo,
	}

	rsp, err := h.svc.LikeTag(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleLogUserViewArticle 记录文章浏览
//
//	@Summary		记录文章浏览
//	@Description	记录文章浏览
//	@Tags			operation
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.LogUserViewArticleBody	true	"记录文章浏览请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.LogArticleViewResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/operation/view/article [post]
func (h *operationHandler) HandleLogUserViewArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.LogUserViewArticleBody)

	req := &protocol.LogArticleViewRequest{
		UserID:    userID,
		ArticleID: body.ArticleID,
		Progress:  body.Progress,
	}

	rsp, err := h.svc.LogArticleView(req)

	util.SendHTTPResponse(c, rsp, err)
}
