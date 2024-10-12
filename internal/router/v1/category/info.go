package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// GetCategoryInfoHandler 获取分类信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 04:58:27
func GetCategoryInfoHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code: protocol.CodeNotPermissionError,
		})
		return
	}

	category, err := model.QueryCategoryByID(uri.CategoryID, []string{"id", "name", "parent_id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: category.GetBasicInfo(),
	})
}

// UpdateCategoryInfoHandler 更新分类信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-02 03:45:55
func UpdateCategoryInfoHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	body := c.MustGet("body").(*protocol.UpdateCategoryBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code: protocol.CodeNotPermissionError,
		})
		return
	}

	updateFields := make(map[string]interface{})

	if body.Name != "" {
		updateFields["name"] = body.Name
	}

	if body.ParentID != 0 {
		updateFields["parent_id"] = body.ParentID
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateCategoryError,
			Message: "No fields to update",
		})
		return
	}

	category, err := model.QueryCategoryByID(uri.CategoryID, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	category, err = model.UpdateCategoryInfoByID(category.ID, updateFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: category.GetBasicInfo(),
	})
}

// DeleteCategoryHandler 删除分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-02 04:55:08
func DeleteCategoryHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code: protocol.CodeNotPermissionError,
		})
		return
	}

	category, err := model.QueryCategoryByID(uri.CategoryID, []string{"id", "parent_id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	if category.ParentID == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteCategoryError,
			Message: "Root category can not be deleted",
		})
		return
	}

	err = model.ReclusiveDeleteCategoryByID(category.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
