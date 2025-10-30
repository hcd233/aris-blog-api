// Package service 业务逻辑
//
//	update 2025-01-04 21:13:05
package service

import (
	"context"
	"errors"

	"github.com/hcd233/aris-blog-api/internal/auth"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TokenService 令牌服务
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type TokenService interface {
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (rsp *dto.RefreshTokenResponse, err error)
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
//	update 2025-01-05 21:00:00
func NewTokenService() TokenService {
	return &tokenService{
		userDAO:            dao.GetUserDAO(),
		accessTokenSigner:  auth.GetJwtAccessTokenSigner(),
		refreshTokenSigner: auth.GetJwtRefreshTokenSigner(),
	}
}

// RefreshToken 刷新令牌
//
//	receiver s *tokenService
//	param ctx context.Context
//	param req *dto.RefreshTokenRequest
//	return rsp *dto.RefreshTokenResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 21:00:00
func (s *tokenService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (rsp *dto.RefreshTokenResponse, err error) {
	rsp = &dto.RefreshTokenResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID, err := s.refreshTokenSigner.DecodeToken(req.Body.RefreshToken)
	if err != nil {
		logger.Error("[TokenService] failed to decode refresh token", zap.String("refreshToken", req.Body.RefreshToken), zap.Error(err))
		return nil, protocol.ErrUnauthorized
	}

	_, err = s.userDAO.GetByID(db, userID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TokenService] user not found", zap.Uint("userID", userID))
			return nil, protocol.ErrDataNotExists
		}
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

	logger.Info("[TokenService] refresh token success", zap.Uint("userID", userID))

	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken

	return rsp, nil
}
