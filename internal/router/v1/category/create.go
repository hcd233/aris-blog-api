package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// CreateCategoryHandler 创建分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-28 07:03:28
func CreateCategoryHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateCategoryBody)
	category := &model.Category{
		Name:     body.Name,
		ParentID: body.ParentID,
		UserID:   userID,
	}

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code: protocol.CodeNotPermissionError,
		})
		return
	}

	if err := category.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeCreateCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: category.GetBasicInfo(),
	})
}
