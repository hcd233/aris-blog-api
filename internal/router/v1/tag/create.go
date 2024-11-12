package tag

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

// CreateTagHandler 创建标签
func CreateTagHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	body := c.MustGet("body").(*protocol.CreateTagBody)

	db := database.GetDBInstance()

	tagDAO := dao.GetTagDAO()
	docDAO := docdao.GetTagDocDAO()

	tag := &model.Tag{
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		User:        &model.User{ID: userID, Name: userName},
	}

	var wg sync.WaitGroup
	var createTagErr, createDocErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		createTagErr = tagDAO.Create(db, tag)
	}()

	go func() {
		defer wg.Done()
		createDocErr = docDAO.AddDocument(document.TransformTagToDocument(tag))
	}()

	wg.Wait()

	if createTagErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: createTagErr.Error(),
		})
		return
	}

	if createDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: createDocErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: tag.GetBasicInfo(),
	})
}
