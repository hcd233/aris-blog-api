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
