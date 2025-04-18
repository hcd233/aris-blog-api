// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"github.com/hcd233/aris-blog-api/internal/auth"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TokenService 令牌服务
//
//	author centonhuang
//	update 2025-01-04 17:16:27
type TokenService interface {
	RefreshToken(req *protocol.RefreshTokenRequest) (rsp *protocol.RefreshTokenResponse, err error)
}

type tokenService struct {
	db                 *gorm.DB
	userDAO            *dao.UserDAO
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewTokenService 创建令牌服务
//
//	return TokenService
//	author centonhuang
//	update 2025-01-04 17:18:59
func NewTokenService() TokenService {
	return &tokenService{
		db:                 database.GetDBInstance(),
		userDAO:            dao.GetUserDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

func (s *tokenService) RefreshToken(req *protocol.RefreshTokenRequest) (rsp *protocol.RefreshTokenResponse, err error) {
	rsp = &protocol.RefreshTokenResponse{}

	userID, err := s.refreshTokenSigner.DecodeToken(req.RefreshToken)
	if err != nil {
		logger.Logger.Error("[TokenService] failed to decode refresh token", zap.String("refreshToken", req.RefreshToken), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	_, err = s.userDAO.GetByID(s.db, userID, []string{"id"}, []string{})
	if err != nil {
		logger.Logger.Error("[TokenService] failed to get user by id", zap.Uint("userID", userID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	accessToken, err := s.accessTokenSigner.EncodeToken(userID)
	if err != nil {
		logger.Logger.Error("[TokenService] failed to encode access token", zap.Uint("userID", userID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(userID)
	if err != nil {
		logger.Logger.Error("[TokenService] failed to encode refresh token", zap.Uint("userID", userID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	logger.Logger.Info("[TokenService] refresh token success", zap.Uint("userID", userID), zap.String("accessToken", accessToken), zap.String("refreshToken", refreshToken))

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}
