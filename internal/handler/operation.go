package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// OperationHandler 用户操作处理器
//
//	@author centonhuang
//	@update 2025-01-04 15:52:48
type OperationHandler interface {
	HandleUserLikeArticle(c *gin.Context)
	HandleUserLikeComment(c *gin.Context)
	HandleUserLikeTag(c *gin.Context)
	HandleLogUserViewArticle(c *gin.Context)
}

type operationHandler struct {
	db          *gorm.DB
	userDAO     *dao.UserDAO
	tagDAO      *dao.TagDAO
	articleDAO  *dao.ArticleDAO
	commentDAO  *dao.CommentDAO
	userLikeDAO *dao.UserLikeDAO
	userViewDAO *dao.UserViewDAO
}

// NewOperationHandler 创建用户操作处理器
//
//	@return OperationHandler
//	@author centonhuang
//	@update 2025-01-04 15:52:48
func NewOperationHandler() OperationHandler {
	return &operationHandler{
		db:          database.GetDBInstance(),
		userDAO:     dao.GetUserDAO(),
		tagDAO:      dao.GetTagDAO(),
		articleDAO:  dao.GetArticleDAO(),
		commentDAO:  dao.GetCommentDAO(),
		userLikeDAO: dao.GetUserLikeDAO(),
		userViewDAO: dao.GetUserViewDAO(),
	}
}

// UserLikeArticleHandler 点赞文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-30 05:52:24
func (h *operationHandler) HandleUserLikeArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LikeArticleBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's like",
		})
		return
	}

	user, err := h.userDAO.GetByName(h.db, body.Author, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, body.ArticleSlug, user.ID, []string{"id", "likes", "status", "user_id"}, []string{})
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

	tx := h.db.Begin()
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
		err = h.transactUndoLikeArticle(tx, article, userLike)
	} else {
		err = h.transactLikeArticle(tx, article, userLike)
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

// UserLikeCommentHandler 点赞评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-30 05:52:08
func (h *operationHandler) HandleUserLikeComment(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LikeCommentBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's like",
		})
		return
	}

	comment, err := h.commentDAO.GetByID(h.db, body.CommentID, []string{"id", "likes"}, []string{})
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

	tx := h.db.Begin()
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
		err = h.transactUndoLikeComment(tx, comment, userLike)
	} else {
		err = h.transactLikeComment(tx, comment, userLike)
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

// UserLikeTagHandler 点赞标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-29 07:10:45
func (h *operationHandler) HandleUserLikeTag(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LikeTagBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's like",
		})
		return
	}

	tag, err := h.tagDAO.GetBySlug(h.db, body.TagSlug, []string{"id", "likes"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	userLike := &model.UserLike{
		UserID:     userID,
		ObjectID:   tag.ID,
		ObjectType: model.LikeObjectTypeTag,
	}

	tx := h.db.Begin()
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
		err = h.transactUndoLikeTag(tx, tag, userLike)
	} else {
		err = h.transactLikeTag(tx, tag, userLike)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeLikeTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}

// LogUserViewArticleHandler 记录文章浏览
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:45:42
func (h *operationHandler) HandleLogUserViewArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.LogUserViewArticleBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's view",
		})
		return
	}

	user, err := h.userDAO.GetByName(h.db, body.Author, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, body.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
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

	userView, err := h.userViewDAO.GetLatestViewByUserIDAndArticleID(h.db, userID, article.ID, []string{"id", "created_at", "progress"}, []string{})
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

		if err = h.userViewDAO.Create(h.db, userView); err != nil {
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

		if err = h.userViewDAO.Update(h.db, userView, map[string]interface{}{"progress": body.Progress, "last_viewed_at": time.Now()}); err != nil {
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

func (h *operationHandler) transactLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) (err error) {
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

func (h *operationHandler) transactUndoLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) (err error) {
	userLikeWithID, err := h.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = h.userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = h.articleDAO.Update(tx, article, map[string]interface{}{"likes": article.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}

	return
}

func (h *operationHandler) transactLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) (err error) {
	if err = h.userLikeDAO.Create(tx, userLike); err != nil {
		err = fmt.Errorf("transaction create user like failed: %v", err)
		return
	}

	if err = h.commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes + 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}
	return
}

func (h *operationHandler) transactUndoLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) (err error) {
	userLikeWithID, err := h.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = h.userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = h.commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}

	return
}

func (h *operationHandler) transactLikeTag(tx *gorm.DB, tag *model.Tag, userLike *model.UserLike) (err error) {
	if err = h.userLikeDAO.Create(tx, userLike); err != nil {
		err = fmt.Errorf("transaction create user like failed: %v", err)
		return
	}

	if err = h.tagDAO.Update(tx, tag, map[string]interface{}{"likes": tag.Likes + 1}); err != nil {
		err = fmt.Errorf("transaction update tag likes failed: %v", err)
		return
	}
	return
}

func (h *operationHandler) transactUndoLikeTag(tx *gorm.DB, tag *model.Tag, userLike *model.UserLike) (err error) {
	userLikeWithID, err := h.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = h.userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = h.tagDAO.Update(tx, tag, map[string]interface{}{"likes": tag.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update tag likes failed: %v", err)
		return
	}

	return
}
