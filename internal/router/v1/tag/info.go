package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	"github.com/samber/lo"
)

// GetTagInfoHandler 获取标签信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 04:58:01
func GetTagInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TagURI)

	dao := dao.GetTagDAO()

	tag, err := dao.GetBySlug(database.DB, uri.TagSlug, []string{"id", "name", "slug", "description", "create_by"})
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
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	dao := dao.GetTagDAO()

	tag, err := dao.GetBySlug(database.DB, uri.TagSlug, []string{"id"})
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
	err = dao.Update(database.DB, tag, updateFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	tag = lo.Must1(dao.GetBySlug(database.DB, uri.TagSlug, []string{"id"}))

	// 同步到搜索引擎
	err = search.UpdateTagInIndex(tag.GetDetailedInfo())
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
	uri := c.MustGet("uri").(*protocol.TagURI)

	dao := dao.GetTagDAO()

	tag, err := dao.GetBySlug(database.DB, uri.TagSlug, []string{"id", "name", "slug"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	err = dao.Delete(database.DB, tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteTagError,
			Message: err.Error(),
		})
		return
	}
	search.DeleteTagFromIndex(tag.ID)

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
