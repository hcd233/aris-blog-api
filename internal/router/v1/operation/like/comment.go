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

// UserLikeCommentHandler 点赞评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-30 05:52:08
func UserLikeCommentHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LikeCommentBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's like",
		})
		return
	}

	db := database.GetDBInstance()

	commentDAO := dao.GetCommentDAO()

	comment, err := commentDAO.GetByID(db, body.CommentID, []string{"id", "likes"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCommentError,
			Message: err.Error(),
		})
		return
	}

	userLike := &model.UserLike{
		UserID:     userID,
		ObjectID:   comment.ID,
		ObjectType: model.LikeObjectTypeComment,
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
		err = transactUndoLikeComment(tx, comment, userLike)
	} else {
		err = transactLikeComment(tx, comment, userLike)
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

func transactLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) (err error) {
	commentDAO, userLikeDAO := dao.GetCommentDAO(), dao.GetUserLikeDAO()
	if err = userLikeDAO.Create(tx, userLike); err != nil {
		err = fmt.Errorf("transaction create user like failed: %v", err)
		return
	}

	if err = commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes + 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}
	return
}

func transactUndoLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) (err error) {
	commentDAO, userLikeDAO := dao.GetCommentDAO(), dao.GetUserLikeDAO()
	userLikeWithID, err := userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}

	return
}
