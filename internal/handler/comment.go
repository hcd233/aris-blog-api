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

// CommentService 评论服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type CommentService interface {
	CreateArticleCommentHandler(c *gin.Context)
	GetCommentInfoHandler(c *gin.Context)
	DeleteCommentHandler(c *gin.Context)
	ListArticleCommentsHandler(c *gin.Context)
	ListChildrenCommentsHandler(c *gin.Context)
}

type commentService struct {
	db         *gorm.DB
	userDAO    *dao.UserDAO
	articleDAO *dao.ArticleDAO
	commentDAO *dao.CommentDAO
}

// NewCommentService 创建评论服务
//
//	@return CommentService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewCommentService() CommentService {
	return &commentService{
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
func (s *commentService) CreateArticleCommentHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.CreateArticleCommentBody)

	author, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, author.ID, []string{"id", "status"}, []string{})
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
		parent, err = s.commentDAO.GetByID(s.db, body.ReplyTo, []string{"id", "article_id"}, []string{})
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

	if err := s.commentDAO.Create(s.db, comment); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateCommentError,
			Message: err.Error(),
		})
		return
	}

	comment = lo.Must1(s.commentDAO.GetByID(s.db, comment.ID, []string{"id", "created_at", "content", "parent_id", "user_id"}, []string{"User", "Parent"}))

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
func (s *commentService) GetCommentInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	comment, err := s.commentDAO.GetByArticleIDAndID(s.db, article.ID, uri.CommentID, []string{"id", "created_at", "content", "user_id", "parent_id", "article_id"}, []string{"User", "Parent", "Article"})
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
func (s *commentService) DeleteCommentHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.CommentURI)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	comment, err := s.commentDAO.GetByArticleIDAndID(s.db, article.ID, uri.CommentID, []string{"id", "user_id"}, []string{})
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

	if err := s.commentDAO.DeleteReclusiveByID(s.db, comment.ID, []string{"id"}, []string{}); err != nil {
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
func (s *commentService) ListArticleCommentsHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
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

	comments, pageInfo, err := s.commentDAO.PaginateRootsByArticleID(s.db, article.ID, []string{"id", "content", "created_at", "likes", "user_id"}, []string{"User"}, param.Page, param.PageSize)
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
			"comments": lo.Map(*comments, func(comment model.Comment, idx int) map[string]interface{} {
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
func (s *commentService) ListChildrenCommentsHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	parentComment, err := s.commentDAO.GetByID(s.db, uri.CommentID, []string{"id", "article_id"}, []string{})
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

	comments, pageInfo, err := s.commentDAO.PaginateChildren(s.db, parentComment, []string{"id", "content", "created_at", "likes", "user_id"}, []string{"User"}, param.Page, param.PageSize)
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
			"comments": lo.Map(*comments, func(comment model.Comment, idx int) map[string]interface{} {
				return comment.GetDetailedInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
