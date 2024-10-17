package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

// CreateTagHandler 创建标签
func CreateTagHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	body := c.MustGet("body").(*protocol.CreateTagBody)

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	dao := dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	tag := &model.Tag{
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		CreateBy:    userID,
	}

	if err := dao.Create(db, tag); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: err.Error(),
		})
		return
	}

	// 同步到搜索引擎
	if err := docDAO.AddDocument(searchEngine, document.TransformTagToDocument(tag)); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
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
