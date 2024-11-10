package view

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// LogUserViewArticleHandler 记录文章浏览
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:45:42
func LogUserViewArticleHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LogUserViewArticleBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's view",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, articleDAO, userViewDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetUserViewDAO()

	user, err := userDAO.GetByName(db, body.Author, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, body.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
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
			Message: "You have no permission to like this article",
		})
		return
	}

	userView, err := userViewDAO.GetLatestViewByUserIDAndArticleID(db, userID, article.ID, []string{"id", "created_at", "progress"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserViewError,
			Message: err.Error(),
		})
		return
	}

	if userView.ID == 0 || userView.Progress >= 95 { // 未浏览过或浏览进度大于95%，则创建新的浏览记录
		if body.Progress != 0 {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeLogUserViewError,
				Message: "Invalid progress",
			})
			return
		}

		userView = &model.UserView{
			UserID:    userID,
			ArticleID: article.ID,
			Progress:  body.Progress,
		}

		if err = userViewDAO.Create(db, userView); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeLogUserViewError,
				Message: err.Error(),
			})
			return
		}

	} else { // 更新浏览记录
		if body.Progress-userView.Progress < 5 { // 浏览进度增加小于5%，则忽略此次请求
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeLogUserViewError,
				Message: "Log view too frequently. Ignore this request",
			})
			return
		}

		if err = userViewDAO.Update(db, userView, map[string]interface{}{"progress": body.Progress, "last_viewed_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateUserViewError,
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
