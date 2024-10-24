package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
)

// GetCommentInfoHandler 获取评论信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-24 05:58:29
func GetCommentInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CommentURI)

	db := database.GetDBInstance()

	userDAO, articleDAO, commentDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetCommentDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	comment, err := commentDAO.GetByArticleIDAndID(db, article.ID, uri.CommentID, []string{"id", "created_at", "content", "parent_id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: comment.GetBasicInfo(),
	})
}
