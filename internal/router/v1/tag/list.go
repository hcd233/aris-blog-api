// Package tag 标签接口
//
//	@update 2024-09-22 02:40:12
package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListTagsHandler 标签列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-22 02:41:01
func ListTagsHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	dao := dao.GetTagDAO()

	tags, pageInfo, err := dao.Paginate(db, []string{"id", "slug", "title"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags": lo.Map(*tags, func(tag model.Tag, index int) map[string]interface{} {
				return tag.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserTagsHandler 列出用户标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-22 02:41:01
func ListUserTagsHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	userDAO, tagDAO := dao.GetUserDAO(), dao.GetTagDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tags, pageInfo, err := tagDAO.PaginateByUserID(db, user.ID, []string{"id", "slug", "name"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags": lo.Map(*tags, func(tag model.Tag, index int) map[string]interface{} {
				return tag.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
