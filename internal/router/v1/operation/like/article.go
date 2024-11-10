// Package like 用户点赞接口
//
//	@update 2024-10-29 07:09:59
package like

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserLikeArticleHandler 点赞文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-30 05:52:24
func UserLikeArticleHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LikeArticleBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's like",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, articleDAO := dao.GetUserDAO(), dao.GetArticleDAO()

	user, err := userDAO.GetByName(db, body.Author, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, body.ArticleSlug, user.ID, []string{"id", "likes", "status", "user_id"}, []string{})
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

	userLike := &model.UserLike{
		UserID:     userID,
		ObjectID:   article.ID,
		ObjectType: model.LikeObjectTypeArticle,
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if body.Undo {
		err = transactUndoLikeArticle(tx, article, userLike)
	} else {
		err = transactLikeArticle(tx, article, userLike)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeLikeCommentError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}

func transactLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) (err error) {
	articleDAO, userLikeDAO := dao.GetArticleDAO(), dao.GetUserLikeDAO()
	if err = userLikeDAO.Create(tx, userLike); err != nil {
		err = fmt.Errorf("transaction create user like failed: %v", err)
		return
	}

	if err = articleDAO.Update(tx, article, map[string]interface{}{"likes": article.Likes + 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}
	return
}

func transactUndoLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) (err error) {
	articleDAO, userLikeDAO := dao.GetArticleDAO(), dao.GetUserLikeDAO()
	userLikeWithID, err := userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = articleDAO.Update(tx, article, map[string]interface{}{"likes": article.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}

	return
}
