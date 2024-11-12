package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	doc_dao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
)

// QueryTagHandler 搜索标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:35:21
func QueryTagHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	docDAO := doc_dao.GetTagDocDAO()

	tags, queryInfo, err := docDAO.QueryDocument(param.Query, param.Filter, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags":      tags,
			"queryInfo": queryInfo,
		},
	})
}

// QueryUserTagHandler 搜索用户标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:35:31
func QueryUserTagHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	params := c.MustGet("param").(*protocol.QueryParam)

	db := database.GetDBInstance()

	userDAO := dao.GetUserDAO()
	docDAO := doc_dao.GetTagDocDAO()

	_, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tags, queryInfo, err := docDAO.QueryDocument(params.Query, append(params.Filter, "creator="+uri.UserName), params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags":      tags,
			"queryInfo": queryInfo,
		},
	})
}
