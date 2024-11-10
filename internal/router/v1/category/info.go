package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
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
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}
	db := database.GetDBInstance()

	categoryDAO, userDAO := dao.GetCategoryDAO(), dao.GetUserDAO()

	_, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	category, err := categoryDAO.GetByID(db, uri.CategoryID, []string{"id", "name", "parent_id"}, []string{})
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

// ListRootCategoriesHandler 获取根分类信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 03:56:26
func ListRootCategoriesHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)

	db := database.GetDBInstance()

	categoryDAO, userDAO := dao.GetCategoryDAO(), dao.GetUserDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's root category",
		})
		return
	}

	rootCategory, err := categoryDAO.GetRootByUserID(db, user.ID, []string{"id", "name"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: rootCategory.GetBasicInfo(),
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
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's category",
		})
		return
	}

	db := database.GetDBInstance()

	categoryDAO, userDAO := dao.GetCategoryDAO(), dao.GetUserDAO()

	_, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
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

	category, err := categoryDAO.GetByID(db, uri.CategoryID, []string{"id", "name", "parent_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	err = categoryDAO.Update(db, category, updateFields)
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
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's category",
		})
		return
	}

	db := database.GetDBInstance()

	categoryDAO, userDAO := dao.GetCategoryDAO(), dao.GetUserDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	category, err := categoryDAO.GetByID(db, uri.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	if user.ID != category.UserID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's category",
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

	err = categoryDAO.DeleteReclusiveByID(db, category.ID, []string{"id", "name"}, []string{})
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
