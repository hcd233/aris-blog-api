package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
	"github.com/samber/lo"
)

// GetTagInfoHandler 获取标签信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 04:58:01
func GetTagInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TagURI)

	db := database.GetDBInstance()

	userDAO, tagDAO := dao.GetUserDAO(), dao.GetTagDAO()

	_, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tag, err := tagDAO.GetAllBySlug(db, uri.TagSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: tag.GetDetailedInfoWithUser(),
	})
}

// UpdateTagHandler 更新标签
func UpdateTagHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's tag",
		})
		return
	}

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO, tagDAO := dao.GetUserDAO(), dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
	}

	tag, err := tagDAO.GetBySlugAndUserID(db, uri.TagSlug, user.ID, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	updateFields := make(map[string]interface{})
	if body.Name != "" {
		updateFields["name"] = body.Name
	}
	if body.Slug != "" {
		updateFields["slug"] = body.Slug
	}
	if body.Description != "" {
		updateFields["description"] = body.Description
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: "No fields to update",
		})
		return
	}
	err = tagDAO.Update(db, tag, updateFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	tag = lo.Must1(tagDAO.GetBySlugAndUserID(db, uri.TagSlug, user.ID, []string{"id"}))

	// 同步到搜索引擎
	err = docDAO.UpdateDocument(searchEngine, document.TransformTagToDocument(tag))
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: tag.GetDetailedInfo(),
	})
}

// DeleteTagHandler 删除标签
func DeleteTagHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.TagURI)
	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's tag",
		})
		return
	}

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()
	userDAO, tagDAO := dao.GetUserDAO(), dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
	}

	tag, err := tagDAO.GetBySlugAndUserID(db, uri.TagSlug, user.ID, []string{"id", "name", "slug"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	err = tagDAO.Delete(db, tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteTagError,
			Message: err.Error(),
		})
		return
	}
	docDAO.DeleteDocument(searchEngine, tag.ID)

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
