package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
)

// QueryArticleHandler 查询文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-18 03:02:01
func QueryArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	params := c.MustGet("param").(*protocol.QueryParam)

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO := dao.GetUserDAO()
	docDAO := docdao.GetArticleDocDAO()

	_, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	articles, err := docDAO.QueryDocument(searchEngine, params.Query, append(params.Filter, "author="+uri.UserName), params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles": articles,
		},
	})
}
