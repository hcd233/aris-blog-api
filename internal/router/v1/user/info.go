// Package user 用户接口
//
//	@update 2024-09-16 05:29:08
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hcd233/Aris-AI-go/internal/protocol"
	"github.com/hcd233/Aris-AI-go/internal/resource/database/model"
)

// GetInfoHandler 用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func GetInfoHandler(c *gin.Context) {
	userName, ok := c.Params.Get("userName")
	if !ok {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeParamError,
		})
		return
	}

	user, err := model.QueryUserByName(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code: protocol.CodeGetUserError,
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, protocol.Response{
			Code: protocol.CodeUserNotFoundError,
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"created_at": user.CreatedAt,
			"last_login": user.LastLogin.Time,
		},
	})
}
