package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/search"
)

// SearchTagHandler 搜索标签
func SearchTagHandler(c *gin.Context) {
	params := c.MustGet("param").(*protocol.QueryParam)

	query, limit, offset := params.Query, params.Limit, params.Offset
	tags, err := search.QueryTagFromIndex(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
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
