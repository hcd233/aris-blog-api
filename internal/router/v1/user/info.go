// Package user 用户接口
//
//	@update 2024-09-16 05:29:08
package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
)

// GetUserInfoHandler 用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func GetUserInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)

	dao := dao.GetUserDAO()
	user, err := dao.GetByName(database.DB, uri.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUserNotFoundError,
			Message: fmt.Sprintf("User `%s` not found", uri.UserName),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: user.GetDetailedInfo(),
	})
}

// UpdateInfoHandler 更新用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-18 01:54:05
func UpdateInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.UpdateUserBody)
	userID, userName := c.GetUint("userID"), c.GetString("userName")

	dao := dao.GetUserDAO()

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's info",
		})
		return
	}

	user, err := dao.GetByName(database.DB, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUserNotFoundError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(dao.Update(database.DB, &model.User{ID: userID}, map[string]interface{}{
		"name": body.UserName,
	}))
	user = lo.Must1(dao.GetByID(database.DB, userID, []string{"id", "name", "avatar"}))

	search.UpdateUserInIndex(user.GetBasicInfo())

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Update user info successfully",
	})
}
