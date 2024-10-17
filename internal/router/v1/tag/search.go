package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
)

// SearchTagHandler 搜索标签
func SearchTagHandler(c *gin.Context) {
	params := c.MustGet("param").(*protocol.QueryParam)

	searchEngine := search.GetSearchEngine()

	docDAO := docdao.GetTagDocDAO()

	tags, err := docDAO.QueryDocument(searchEngine, params.Query, params.Limit, params.Offset)
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
			"tags": tags,
		},
	})
}
