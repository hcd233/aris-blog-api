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

// TokenService 令牌服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type TokenService interface {
	RefreshTokenHandler(c *gin.Context)
}

type tokenService struct {
	db                 *gorm.DB
	userDAO            *dao.UserDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewTokenService 创建令牌服务
//
//	@return TokenService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewTokenService() TokenService {
	return &tokenService{
		db:                 database.GetDBInstance(),
		userDAO:            dao.GetUserDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

func (s *tokenService) RefreshTokenHandler(c *gin.Context) {
	body := c.MustGet("body").(*protocol.RefreshTokenBody)

	userID, err := s.refreshTokenSigner.DecodeToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeTokenVerifyError,
			Message: err.Error(),
		})

		return
	}

	_, err = s.userDAO.GetByID(s.db, userID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	accessToken := lo.Must1(s.accessTokenSigner.EncodeToken(userID))
	refreshToken := lo.Must1(s.refreshTokenSigner.EncodeToken(userID))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}
