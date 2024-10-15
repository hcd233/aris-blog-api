package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListArticleHandler 用户文章列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 08:59:40
func ListArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := model.QueryUserByName(uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}
	articles, err := model.QueryArticlesByUserID(user.ID, param.Limit, param.Offset, []string{"id", "title", "slug"})
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
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
				return article.GetBasicInfo()
			}),
		},
	})
}
