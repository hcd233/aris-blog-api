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
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateCategoryBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's category",
		})
		return
	}

	db := database.GetDBInstance()

	categoryDAO, userDAO := dao.GetCategoryDAO(), dao.GetUserDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	var parentCategory *model.Category
	if body.ParentID == 0 {
		parentCategory, err = categoryDAO.GetRootByUserID(db, user.ID, []string{"id"}, []string{})
	} else {
		parentCategory, err = categoryDAO.GetByID(db, body.ParentID, []string{"id"}, []string{})
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}
	body.ParentID = parentCategory.ID

	category := &model.Category{
		Name:     body.Name,
		ParentID: body.ParentID,
		UserID:   user.ID,
	}

	if err := categoryDAO.Create(db, category); err != nil {
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
