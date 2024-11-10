package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListUserViewArticlesHandler 列出用户浏览的文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:45:42
func ListUserViewArticlesHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	uri := c.MustGet("uri").(*protocol.UserURI)
	pageParam := c.MustGet("param").(*protocol.PageParam)

	if uri.UserName != c.MustGet("userName").(string) {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to operate other user's view",
		})
		return
	}

	db := database.GetDBInstance()
	_, userViewDAO := dao.GetArticleDAO(), dao.GetUserViewDAO()

	userViews, pageInfo, err := userViewDAO.PaginateByUserID(db, userID, []string{"id", "progress", "last_viewed_at"}, []string{"User", "Article", "Article.Tags", "Article.User"}, pageParam.Page, pageParam.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserViewError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"userViews": lo.Map(*userViews, func(userView model.UserView, idx int) map[string]interface{} {
				return userView.GetViewInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
