package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

// CreateTagHandler 创建标签
func CreateTagHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateTagBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's tag",
		})
		return
	}

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO, tagDAO := dao.GetUserDAO(), dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id", "name"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
	}

	tag := &model.Tag{
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		User:        user,
	}

	if err := tagDAO.Create(db, tag); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: err.Error(),
		})
		return
	}

	// 同步到搜索引擎
	if err := docDAO.AddDocument(searchEngine, document.TransformTagToDocument(tag)); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: tag.GetBasicInfo(),
	})
}
