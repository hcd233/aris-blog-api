package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/auth"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// TokenHandler 令牌处理器
//
//	@author centonhuang
//	@update 2025-01-04 15:56:10
type TokenHandler interface {
	HandleRefreshToken(c *gin.Context)
}

type tokenHandler struct {
	db                 *gorm.DB
	userDAO            *dao.UserDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewTokenHandler 创建令牌处理器
//
//	@return TokenHandler
//	@author centonhuang
//	@update 2025-01-04 15:56:04
func NewTokenHandler() TokenHandler {
	return &tokenHandler{
		db:                 database.GetDBInstance(),
		userDAO:            dao.GetUserDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

// HandleRefreshToken 刷新令牌
//
//	@receiver s *tokenHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:56:10
func (h *tokenHandler) HandleRefreshToken(c *gin.Context) {
	body := c.MustGet("body").(*protocol.RefreshTokenBody)

	userID, err := h.refreshTokenSigner.DecodeToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeTokenVerifyError,
			Message: err.Error(),
		})

		return
	}

	_, err = h.userDAO.GetByID(h.db, userID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	accessToken := lo.Must1(h.accessTokenSigner.EncodeToken(userID))
	refreshToken := lo.Must1(h.refreshTokenSigner.EncodeToken(userID))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}
