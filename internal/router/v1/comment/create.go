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

// CreateArticleCommentHandler 创建文章评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-24 04:29:47
func CreateArticleCommentHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.CreateArticleCommentBody)

	db := database.GetDBInstance()

	userDAO, articleDAO, commentDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetCommentDAO()

	author, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, author.ID, []string{"id", "status"})
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
		parent, err = commentDAO.GetByID(db, body.ReplyTo, []string{"id", "article_id"}, []string{})
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

	if err := commentDAO.Create(db, comment); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateCommentError,
			Message: err.Error(),
		})
		return
	}

	comment = lo.Must1(commentDAO.GetByID(db, comment.ID, []string{"id", "created_at", "content", "parent_id"}, []string{}))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: comment.GetBasicInfo(),
	})
}
