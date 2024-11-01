package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
)

// QueryUserHandler 查询用户
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func QueryUserHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	searchEngine := search.GetSearchEngine()

	docDAO := docdao.GetUserDocDAO()

	users, queryInfo, err := docDAO.QueryDocument(searchEngine, param.Query, param.Filter, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"users":     users,
			"queryInfo": queryInfo,
		},
	})
}
