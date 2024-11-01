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

// QueryUserArticleHandler 查询用户文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:13:21
func QueryUserArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.QueryParam)

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO := dao.GetUserDAO()
	articleDocDAO := docdao.GetArticleDocDAO()

	_, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	articles, queryInfo, err := articleDocDAO.QueryDocument(searchEngine, param.Query, append(param.Filter, "author="+uri.UserName), param.Page, param.PageSize)
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
			"articles":  articles,
			"queryInfo": queryInfo,
		},
	})
}

// QueryArticleHandler 查询文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:13:17
func QueryArticleHandler(c *gin.Context) {
	params := c.MustGet("param").(*protocol.QueryParam)

	searchEngine := search.GetSearchEngine()

	docDAO := docdao.GetArticleDocDAO()

	articles, queryInfo, err := docDAO.QueryDocument(searchEngine, params.Query, params.Filter, params.Page, params.PageSize)
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
			"articles":  articles,
			"queryInfo": queryInfo,
		},
	})
}
