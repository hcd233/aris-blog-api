package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// CommentHandler 评论处理器
//
//	@author centonhuang
//	@update 2025-01-04 15:52:48
type CommentHandler interface {
	HandleCreateArticleComment(c *gin.Context)
	HandleGetCommentInfo(c *gin.Context)
	HandleDeleteComment(c *gin.Context)
	HandleListArticleComments(c *gin.Context)
	HandleListChildrenComments(c *gin.Context)
}

type commentHandler struct {
	db         *gorm.DB
	userDAO    *dao.UserDAO
	articleDAO *dao.ArticleDAO
	commentDAO *dao.CommentDAO
}

// NewCommentHandler 创建评论处理器
//
//	@return CommentHandler
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewCommentHandler() CommentHandler {
	return &commentHandler{
		db:         database.GetDBInstance(),
		userDAO:    dao.GetUserDAO(),
		articleDAO: dao.GetArticleDAO(),
		commentDAO: dao.GetCommentDAO(),
	}
}

// CreateArticleCommentHandler 创建文章评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-24 04:29:47
func (h *commentHandler) HandleCreateArticleComment(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.CreateArticleCommentBody)

	author, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, author.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if article.Status == model.ArticleStatusDraft {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "How did you find this article?",
		})
		return
	}

	var parent *model.Comment
	if body.ReplyTo != 0 {
		parent, err = h.commentDAO.GetByID(h.db, body.ReplyTo, []string{"id", "article_id"}, []string{})
		if err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeGetCommentError,
				Message: err.Error(),
			})
			return
		}

		if parent.ArticleID != article.ID {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeURIError,
				Message: "The parent comment is not in this article",
			})
			return
		}
	}

	comment := &model.Comment{
		UserID:    userID,
		ArticleID: article.ID,
		Parent:    parent,
		Content:   body.Content,
	}

	if err := h.commentDAO.Create(h.db, comment); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateCommentError,
			Message: err.Error(),
		})
		return
	}

	comment = lo.Must1(h.commentDAO.GetByID(h.db, comment.ID, []string{"id", "created_at", "content", "parent_id", "user_id"}, []string{"User", "Parent"}))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"comment": comment.GetDetailedInfo(),
		},
	})
}

// GetCommentInfoHandler 获取评论信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-24 05:58:29
func (h *commentHandler) HandleGetCommentInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	comment, err := h.commentDAO.GetByArticleIDAndID(h.db, article.ID, uri.CommentID, []string{"id", "created_at", "content", "user_id", "parent_id", "article_id"}, []string{"User", "Parent", "Article"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"comment": comment.GetDetailedInfo(),
		},
	})
}

// DeleteCommentHandler 删除评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-24 07:05:09
func (h *commentHandler) HandleDeleteComment(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.CommentURI)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	comment, err := h.commentDAO.GetByArticleIDAndID(h.db, article.ID, uri.CommentID, []string{"id", "user_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	if userID != comment.UserID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's comment",
		})
		return
	}

	if err := h.commentDAO.DeleteReclusiveByID(h.db, comment.ID, []string{"id"}, []string{}); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteCommentError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}

// ListArticleCommentsHandler 列出文章评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 05:59:57
func (h *commentHandler) HandleListArticleComments(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if userName != uri.UserName && article.Status != model.ArticleStatusPublish { // 非作者且文章未发布
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view this article's comments",
		})
		return
	}

	comments, pageInfo, err := h.commentDAO.PaginateRootsByArticleID(h.db, article.ID, []string{"id", "content", "created_at", "likes", "user_id"}, []string{"User"}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"comments": lo.Map(*comments, func(comment model.Comment, _ int) map[string]interface{} {
				return comment.GetDetailedInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListChildrenCommentsHandler 列出子评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-24 06:55:45
func (h *commentHandler) HandleListChildrenComments(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	parentComment, err := h.commentDAO.GetByID(h.db, uri.CommentID, []string{"id", "article_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	if parentComment.ArticleID != article.ID {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeURIError,
			Message: "The comment is not in this article",
		})
		return
	}

	comments, pageInfo, err := h.commentDAO.PaginateChildren(h.db, parentComment, []string{"id", "content", "created_at", "likes", "user_id"}, []string{"User"}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"comments": lo.Map(*comments, func(comment model.Comment, _ int) map[string]interface{} {
				return comment.GetDetailedInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
