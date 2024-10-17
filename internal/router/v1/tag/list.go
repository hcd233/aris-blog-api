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

	dao := dao.GetTagDAO()

	tags, err := dao.Paginate(database.DB, []string{"id", "slug"}, param.Limit, param.Offset)
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
			"tags": lo.Map(*tags, func(article model.Tag, index int) map[string]interface{} {
				return article.GetBasicInfo()
			}),
		},
	})
}
