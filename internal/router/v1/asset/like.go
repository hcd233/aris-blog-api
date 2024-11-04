// Package asset 用户资产接口
//
//	@update 2024-11-01 07:26:04
package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// ListUserLikeArticlesHandler 列出用户喜欢的文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:45:42
func ListUserLikeArticlesHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view other user's like articles",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, userLikeDAO, articleDAO := dao.GetUserDAO(), dao.GetUserLikeDAO(), dao.GetArticleDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userLikes, pageInfo, err := userLikeDAO.PaginateByUserIDAndObjectType(db, user.ID, model.LikeObjectTypeArticle, param.Page, param.PageSize, []string{"object_id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserLikeError,
			Message: err.Error(),
		})
		return
	}

	articleIDs := lo.Map(*userLikes, func(like model.UserLike, idx int) uint {
		return like.ObjectID
	})

	articles, err := articleDAO.BatchGetAllByIDs(db, articleIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if len(articleIDs) != len(*articles) {
		_, deletedIDs := lo.Difference(articleIDs, lo.Map(*articles, func(article model.Article, idx int) uint {
			return article.ID
		}))

		logger.Logger.Warn("[List User Like Articles]", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, idx int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeArticle,
			}
		})

		err = userLikeDAO.BatchDelete(db, &deleteLikes)
		if err != nil {
			logger.Logger.Error("[List User Like Articles]", zap.Error(err))
			err = nil
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
				return article.GetLikeInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserLikeCommentsHandler 列出用户喜欢的评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:47:41
func ListUserLikeCommentsHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view other user's like comments",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, userLikeDAO, commentDAO := dao.GetUserDAO(), dao.GetUserLikeDAO(), dao.GetCommentDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userLikes, pageInfo, err := userLikeDAO.PaginateByUserIDAndObjectType(db, user.ID, model.LikeObjectTypeComment, param.Page, param.PageSize, []string{"object_id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserLikeError,
			Message: err.Error(),
		})
		return
	}

	commentIDs := lo.Map(*userLikes, func(like model.UserLike, idx int) uint {
		return like.ObjectID
	})

	comments, err := commentDAO.BatchGetAllByIDs(db, commentIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if len(commentIDs) != len(*comments) {
		_, deletedIDs := lo.Difference(commentIDs, lo.Map(*comments, func(comment model.Comment, idx int) uint {
			return comment.ID
		}))

		logger.Logger.Warn("[List User Like Comments]", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, idx int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeComment,
			}
		})

		err = userLikeDAO.BatchDelete(db, &deleteLikes)
		if err != nil {
			logger.Logger.Error("[List User Like Comments]", zap.Error(err))
			err = nil
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"comments": lo.Map(*comments, func(comment model.Comment, index int) map[string]interface{} {
				return comment.GetLikeInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserLikeTagsHandler 列出用户喜欢的标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:47:43
func ListUserLikeTagsHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view other user's like comments",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, userLikeDAO, tagDAO := dao.GetUserDAO(), dao.GetUserLikeDAO(), dao.GetTagDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userLikes, pageInfo, err := userLikeDAO.PaginateByUserIDAndObjectType(db, user.ID, model.LikeObjectTypeTag, param.Page, param.PageSize, []string{"object_id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserLikeError,
			Message: err.Error(),
		})
		return
	}

	tagIDs := lo.Map(*userLikes, func(like model.UserLike, idx int) uint {
		return like.ObjectID
	})

	tags, err := tagDAO.BatchGetAllByIDs(db, tagIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if len(tagIDs) != len(*tags) {
		_, deletedIDs := lo.Difference(tagIDs, lo.Map(*tags, func(tag model.Tag, idx int) uint {
			return tag.ID
		}))

		logger.Logger.Warn("[List User Like Tags]", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, idx int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeTag,
			}
		})

		err = userLikeDAO.BatchDelete(db, &deleteLikes)
		if err != nil {
			logger.Logger.Error("[List User Like Tags]", zap.Error(err))
			err = nil
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags": lo.Map(*tags, func(tag model.Tag, index int) map[string]interface{} {
				return tag.GetLikeInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
