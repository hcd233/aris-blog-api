package tag

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
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

	tagDAO := dao.GetTagDAO()

	tag, err := tagDAO.GetBySlug(db, uri.TagSlug, []string{"id", "name", "slug", "description", "create_by"}, []string{"User"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: tag.GetDetailedInfo(),
	})
}

// UpdateTagHandler 更新标签
func UpdateTagHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	db := database.GetDBInstance()

	tagDAO := dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	tag, err := tagDAO.GetBySlug(db, uri.TagSlug, []string{"id", "create_by"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	if tag.CreatedBy != userID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's tag",
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

	if err := tagDAO.Update(db, tag, updateFields); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	tag = lo.Must1(tagDAO.GetBySlug(db, uri.TagSlug, []string{"id", "name", "slug", "description", "created_by"}, []string{"User"}))

	// 同步到搜索引擎
	err = docDAO.UpdateDocument(document.TransformTagToDocument(tag))
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
	userID := c.MustGet("userID").(uint)
	uri := c.MustGet("uri").(*protocol.TagURI)

	db := database.GetDBInstance()

	tagDAO := dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	tag, err := tagDAO.GetBySlug(db, uri.TagSlug, []string{"id", "name", "slug", "create_by"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	if tag.CreatedBy != userID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's tag",
		})
		return
	}

	var wg sync.WaitGroup
	var deleteTagErr, deleteDocErr error
	wg.Add(2)

	go func() {
		defer wg.Done()
		deleteTagErr = tagDAO.Delete(db, tag)
	}()

	go func() {
		defer wg.Done()
		deleteDocErr = docDAO.DeleteDocument(tag.ID)
	}()

	wg.Wait()

	if deleteTagErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteTagError,
			Message: deleteTagErr.Error(),
		})
		return
	}

	if deleteDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteTagError,
			Message: deleteDocErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
