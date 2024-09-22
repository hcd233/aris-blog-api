package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
)

// CreateTagHandler 创建标签
func CreateTagHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	body := c.MustGet("body").(*protocol.CreateTagBody)

	tag := model.Tag{
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		CreateBy:    userID,
	}

	err := tag.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: err.Error(),
		})
		return
	}

	// 同步到搜索引擎
	err = search.AddTagIntoIndex(tag.GetDetailedInfo())
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: tag.GetBasicInfo(),
	})
}
