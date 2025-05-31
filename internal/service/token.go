// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/auth"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"go.uber.org/zap"
)

// TokenService 令牌服务
//
//	author centonhuang
//	update 2025-01-04 17:16:27
type TokenService interface {
	RefreshToken(ctx context.Context, req *protocol.RefreshTokenRequest) (rsp *protocol.RefreshTokenResponse, err error)
}

type tokenService struct {
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
		userDAO:            dao.GetUserDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

func (s *tokenService) RefreshToken(ctx context.Context, req *protocol.RefreshTokenRequest) (rsp *protocol.RefreshTokenResponse, err error) {
	rsp = &protocol.RefreshTokenResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	userID, err := s.refreshTokenSigner.DecodeToken(req.RefreshToken)
	if err != nil {
		logger.Error("[TokenService] failed to decode refresh token", zap.String("refreshToken", req.RefreshToken), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	_, err = s.userDAO.GetByID(db, userID, []string{"id"}, []string{})
	if err != nil {
		logger.Error("[TokenService] failed to get user by id", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	accessToken, err := s.accessTokenSigner.EncodeToken(userID)
	if err != nil {
		logger.Error("[TokenService] failed to encode access token", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(userID)
	if err != nil {
		logger.Error("[TokenService] failed to encode refresh token", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[TokenService] refresh token success", zap.String("accessToken", accessToken), zap.String("refreshToken", refreshToken))

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}
