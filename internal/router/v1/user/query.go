package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/protocol"
	"github.com/hcd233/Aris-AI-go/internal/resource/search"
)

// QueryUserHandler 查询用户
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func QueryUserHandler(c *gin.Context) {
	var params protocol.QueryUserParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeParamError,
			Message: err.Error(),
		})
		return
	}

	query, limit, offset := params.Query, params.Limit, params.Offset
	users, err := search.QueryUserIndex(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"users": users,
		},
	})
}
