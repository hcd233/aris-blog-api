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

	user, err := model.QueryUserByName(uri.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, protocol.Response{
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

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's info",
		})
		return
	}

	user := lo.Must1(model.UpdateUserInfoByID(userID, map[string]interface{}{
		"name": body.UserName,
	}))
	search.UpdateUserIndex(user.GetBasicInfo())

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Update user info successfully",
	})
}
