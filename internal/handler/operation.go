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

// OperationService 用户操作服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type OperationService interface {
	UserLikeArticleHandler(c *gin.Context)
	UserLikeCommentHandler(c *gin.Context)
	UserLikeTagHandler(c *gin.Context)
	LogUserViewArticleHandler(c *gin.Context)
}

type operationService struct {
	db          *gorm.DB
	userDAO     *dao.UserDAO
	tagDAO      *dao.TagDAO
	articleDAO  *dao.ArticleDAO
	commentDAO  *dao.CommentDAO
	userLikeDAO *dao.UserLikeDAO
	userViewDAO *dao.UserViewDAO
}

// NewOperationService 创建用户操作服务
//
//	@return OperationService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewOperationService() OperationService {
	return &operationService{
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
func (s *operationService) UserLikeArticleHandler(c *gin.Context) {
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

	user, err := s.userDAO.GetByName(s.db, body.Author, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, body.ArticleSlug, user.ID, []string{"id", "likes", "status", "user_id"}, []string{})
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

	tx := s.db.Begin()
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
		err = s.transactUndoLikeArticle(tx, article, userLike)
	} else {
		err = s.transactLikeArticle(tx, article, userLike)
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
func (s *operationService) UserLikeCommentHandler(c *gin.Context) {
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

	comment, err := s.commentDAO.GetByID(s.db, body.CommentID, []string{"id", "likes"}, []string{})
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

	tx := s.db.Begin()
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
		err = s.transactUndoLikeComment(tx, comment, userLike)
	} else {
		err = s.transactLikeComment(tx, comment, userLike)
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
func (s *operationService) UserLikeTagHandler(c *gin.Context) {
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

	tag, err := s.tagDAO.GetBySlug(s.db, body.TagSlug, []string{"id", "likes"}, []string{})
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

	tx := s.db.Begin()
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
		err = s.transactUndoLikeTag(tx, tag, userLike)
	} else {
		err = s.transactLikeTag(tx, tag, userLike)
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
func (s *operationService) LogUserViewArticleHandler(c *gin.Context) {
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

	user, err := s.userDAO.GetByName(s.db, body.Author, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, body.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
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

	userView, err := s.userViewDAO.GetLatestViewByUserIDAndArticleID(s.db, userID, article.ID, []string{"id", "created_at", "progress"}, []string{})
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

		if err = s.userViewDAO.Create(s.db, userView); err != nil {
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

		if err = s.userViewDAO.Update(s.db, userView, map[string]interface{}{"progress": body.Progress, "last_viewed_at": time.Now()}); err != nil {
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

func (s *operationService) transactLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) (err error) {
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

func (s *operationService) transactUndoLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) (err error) {
	userLikeWithID, err := s.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = s.userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = s.articleDAO.Update(tx, article, map[string]interface{}{"likes": article.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}

	return
}

func (s *operationService) transactLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) (err error) {
	if err = s.userLikeDAO.Create(tx, userLike); err != nil {
		err = fmt.Errorf("transaction create user like failed: %v", err)
		return
	}

	if err = s.commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes + 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}
	return
}

func (s *operationService) transactUndoLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) (err error) {
	userLikeWithID, err := s.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = s.userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = s.commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update comment likes failed: %v", err)
		return
	}

	return
}

func (s *operationService) transactLikeTag(tx *gorm.DB, tag *model.Tag, userLike *model.UserLike) (err error) {
	if err = s.userLikeDAO.Create(tx, userLike); err != nil {
		err = fmt.Errorf("transaction create user like failed: %v", err)
		return
	}

	if err = s.tagDAO.Update(tx, tag, map[string]interface{}{"likes": tag.Likes + 1}); err != nil {
		err = fmt.Errorf("transaction update tag likes failed: %v", err)
		return
	}
	return
}

func (s *operationService) transactUndoLikeTag(tx *gorm.DB, tag *model.Tag, userLike *model.UserLike) (err error) {
	userLikeWithID, err := s.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		err = fmt.Errorf("transaction get user like failed: %v", err)
		return
	}

	userLike.ID = userLikeWithID.ID

	if err = s.userLikeDAO.Delete(tx, userLike); err != nil {
		err = fmt.Errorf("transaction delete user like failed: %v", err)
		return
	}

	if err = s.tagDAO.Update(tx, tag, map[string]interface{}{"likes": tag.Likes - 1}); err != nil {
		err = fmt.Errorf("transaction update tag likes failed: %v", err)
		return
	}

	return
}
