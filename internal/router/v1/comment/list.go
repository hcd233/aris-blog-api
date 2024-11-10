// Package comment 评论接口
//
//	@update 2024-10-23 05:56:38
package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListArticleCommentsHandler 列出文章评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 05:59:57
func ListArticleCommentsHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	userDAO, articleDAO, commentDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetCommentDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
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

	comments, pageInfo, err := commentDAO.PaginateRootsByArticleID(db, article.ID, []string{"id", "content", "created_at", "likes"}, []string{}, param.Page, param.PageSize)
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
				return comment.GetBasicInfo()
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
func ListChildrenCommentsHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	userDAO, articleDAO, commentDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetCommentDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	parentComment, err := commentDAO.GetByID(db, uri.CommentID, []string{"id", "article_id"}, []string{})
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

	comments, pageInfo, err := commentDAO.PaginateChildren(db, parentComment, []string{"id", "content", "created_at", "likes"}, []string{}, param.Page, param.PageSize)
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
				return comment.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
