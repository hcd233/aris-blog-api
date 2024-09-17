package user

import (
	"net/http"
	"strconv"

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
	query, ok := c.GetQuery("query")
	limit, err1 := strconv.ParseInt(c.DefaultQuery("limit", "5"), 10, 64)
	offset, err2 := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if !ok || err1 != nil || err2 != nil || limit < 0 || offset < 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code: protocol.CodeInvalidQueryError,
		})
		return
	}

	users, err := search.QueryUserIndex(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code: protocol.CodeQueryUserError,
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
