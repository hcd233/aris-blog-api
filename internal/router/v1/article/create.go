package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// CreateArticleHandler 创建文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 09:58:14
func CreateArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateArticleBody)
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)

	db := database.GetDBInstance()

	userDAO, tagDAO, articleDAO := dao.GetUserDAO(), dao.GetTagDAO(), dao.GetArticleDAO()

	if uri.UserName != userName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's article",
		})
		return
	}

	if body.Slug == "" {
		body.Slug = body.Title
	}

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tags := []model.Tag{}
	for _, tag := range body.Tags {
		tag, err := tagDAO.GetBySlugAndUserID(db, tag, user.ID, []string{"id"})
		if err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeGetTagError,
				Message: err.Error(),
			})
			return
		}
		tags = append(tags, *tag)
	}

	article := &model.Article{
		UserID:     userID,
		Status:     model.ArticleStatusDraft,
		Title:      body.Title,
		Slug:       body.Slug,
		Tags:       tags,
		CategoryID: body.CategoryID,
		Comments:   []model.Comment{},
		Versions:   []model.ArticleVersion{},
	}

	if err := articleDAO.Create(db, article); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: article.GetBasicInfo(),
	})
}
