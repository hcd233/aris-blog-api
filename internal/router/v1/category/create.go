package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
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

	dao := dao.GetCategoryDAO()

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's category",
		})
		return
	}

	var err error
	var parentCategory *model.Category
	if body.ParentID == 0 {
		parentCategory, err = dao.GetRootByUserID(database.DB, userID, []string{"id"})
	} else {
		parentCategory, err = dao.GetByID(database.DB, body.ParentID, []string{"id"})
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}
	body.ParentID = parentCategory.ID

	category := &model.Category{
		Name:     body.Name,
		ParentID: body.ParentID,
		UserID:   userID,
	}

	if err := dao.Create(database.DB, category); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
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
