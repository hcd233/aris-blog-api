package token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/auth"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/samber/lo"
)

func RefreshTokenHandler(c *gin.Context) {
	body := c.MustGet("body").(*protocol.RefreshTokenBody)

	db := database.GetDBInstance()

	userDAO := dao.GetUserDAO()

	jwtAccessTokenSvc := auth.GetJwtAccessTokenSvc()
	jwtRefreshTokenSvc := auth.GetJwtRefreshTokenSvc()

	userID, err := jwtRefreshTokenSvc.DecodeToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeTokenVerifyError,
			Message: err.Error(),
		})

		return
	}

	_, err = userDAO.GetByID(db, userID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	accessToken := lo.Must1(jwtAccessTokenSvc.EncodeToken(userID))
	refreshToken := lo.Must1(jwtRefreshTokenSvc.EncodeToken(userID))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}
